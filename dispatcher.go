package tgp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/pikoUsername/tgp/fsm"
	"github.com/pikoUsername/tgp/fsm/storage"
	"github.com/pikoUsername/tgp/objects"
)

// Dispatcher's purpose is run bot, and comfortable pipeline
// Bot struct uses as API wrapper
// Dispatcher uses as Bot starter
// Another level of abstraction
type Dispatcher struct {
	Bot *Bot
	// Handlers
	MessageHandler       *HandlerObj
	CallbackQueryHandler *HandlerObj
	ChannelPostHandler   *HandlerObj
	PollHandler          *HandlerObj
	ChatMemberHandler    *HandlerObj
	PollAnswerHandler    *HandlerObj
	MyChatMemberHandler  *HandlerObj
	Storage              storage.Storage

	currentUpdate *objects.Update

	// If you want to add onshutdown function
	// just append to this object, :P
	OnWebhookShutdown []OnStartAndShutdownFunc
	OnPollingShutdown []OnStartAndShutdownFunc
	OnWebhookStartup  []OnStartAndShutdownFunc
	OnPollingStartup  []OnStartAndShutdownFunc

	Welcome    bool
	Synchronus bool
	polling    bool
	webhook    bool
}

var (
	ErrorTypeAssertion = tgpErr.New("impossible to do type assertion to this callback")
	ErrorConflictModes = tgpErr.New("enabled two conflicting modes at the same time, polling and webhook")
)

type OnStartAndShutdownFunc func(dp *Dispatcher)

// NewDispathcer get a new Dispatcher with default values
// settings -
// 		syncronus: true
// 		welcome: true
func NewDispatcher(bot *Bot, storage storage.Storage) *Dispatcher {
	dp := &Dispatcher{
		Bot:        bot,
		Synchronus: true,
		Storage:    storage,
		Welcome:    true,
	}

	dp.MessageHandler = NewHandlerObj(dp)
	dp.CallbackQueryHandler = NewHandlerObj(dp)
	dp.ChannelPostHandler = NewHandlerObj(dp)
	dp.ChatMemberHandler = NewHandlerObj(dp)
	dp.PollHandler = NewHandlerObj(dp)
	dp.PollAnswerHandler = NewHandlerObj(dp)
	dp.ChannelPostHandler = NewHandlerObj(dp)

	return dp
}

// OnConfig using as argument for OnStartup, OnShutdown methods
// You can add multiple functions to startup, or shutdown mthds
// Example:
// c := &OnConfig{}
// // could be added a multiple functions in one call
// c.Add(func(...) {})
// dp.OnStartup(c)
type OnConfig struct {
	Polling bool
	Webhook bool
	cb      []OnStartAndShutdownFunc
}

func (oc *OnConfig) Add(cb OnStartAndShutdownFunc) {
	oc.cb = append(oc.cb, cb)
}

func NewOnConf(cb OnStartAndShutdownFunc) *OnConfig {
	return &OnConfig{
		cb:      []OnStartAndShutdownFunc{cb},
		Webhook: true,
		Polling: true,
	}
}

func callListFuncs(funcs []OnStartAndShutdownFunc, dp *Dispatcher) {
	for _, cb := range funcs {
		if dp.Synchronus {
			cb(dp)
		} else {
			go cb(dp)
		}
	}
}

// Config for start polling method
// idk where to put this config, configs or dispatcher?
type StartPollingConfig struct {
	*GetUpdatesConfig
	SkipUpdates  bool
	SafeExit     bool
	ResetWebhook bool
	ErrorSleep   uint
	Relax        time.Duration
	Timeout      time.Duration
}

func NewStartPollingConf(skip_updates bool) *StartPollingConfig {
	return &StartPollingConfig{
		GetUpdatesConfig: &GetUpdatesConfig{
			Timeout: 20,
			Limit:   0,
		},
		Relax:        1 * time.Second,
		ResetWebhook: false,
		ErrorSleep:   5,
		SkipUpdates:  skip_updates,
		SafeExit:     true,
		Timeout:      5 * time.Second,
	}
}

type StartWebhookConfig struct {
	*SetWebhookConfig
	Handler            http.Handler
	KeyFile            interface{}
	BotURL             string
	Address            string
	DropPendingUpdates bool
	SafeExit           bool
}

func NewStartWebhookConf(url string, address string) *StartWebhookConfig {
	return &StartWebhookConfig{
		BotURL:  url,
		Address: address,
	}
}

// ResetWebhook uses for reset webhook for telegram
func (dp *Dispatcher) ResetWebhook(check bool) error {
	if check {
		wi, err := dp.Bot.GetWebhookInfo()
		if err != nil {
			return err
		}
		if wi.URL == "" {
			return errors.New("url is nothing")
		}
	}
	_, err := dp.Bot.DeleteWebhook(&DeleteWebhookConfig{})
	return err
}

// ProcessOneUpdate you guess, processes ONLY one comming update
// Support only one Message update
func (dp *Dispatcher) ProcessOneUpdate(update *objects.Update) error {
	var err error

	// very bad code, please dont see this bullshit
	// ============================================
	if update.Message != nil {
		dp.MessageHandler.TriggerMiddleware(dp.Bot, update, PREMIDDLEWARE)
		for _, h := range dp.MessageHandler.handlers {
			cb, ok := h.Callback.(func(*Bot, *objects.Message))
			if !ok {
				return tgpErr.New("message handler type assertion error, need type func(*Bot, *Message), current type is - " + fmt.Sprintln(reflect.TypeOf(h.Callback)))
			}

			err = dp.MessageHandler.TriggerMiddleware(dp.Bot, update, PROCESSMIDDLEWARE)
			if err != nil {
				dp.Bot.logger.Println(err)
				continue
			}

			h.Call(update, func() { cb(dp.Bot, update.Message) })
		}
		dp.MessageHandler.TriggerMiddleware(dp.Bot, update, POSTMIDDLEWARE)

	} else if update.CallbackQuery != nil {
		dp.CallbackQueryHandler.TriggerMiddleware(dp.Bot, update, PREMIDDLEWARE)
		for _, h := range dp.CallbackQueryHandler.handlers {
			cb, ok := h.Callback.(func(*Bot, *objects.CallbackQuery))
			if !ok {
				return tgpErr.New("callbackquery handler type assertion error, need type func(*Bot, *CallbackQuery), current type is - " + fmt.Sprintln(reflect.TypeOf(h.Callback)))
			}
			err = dp.CallbackQueryHandler.TriggerMiddleware(dp.Bot, update, PROCESSMIDDLEWARE)
			if err != nil {
				dp.Bot.logger.Println(err)
				continue
			}

			h.Call(update, func() { cb(dp.Bot, update.CallbackQuery) })
		}
		dp.CallbackQueryHandler.TriggerMiddleware(dp.Bot, update, POSTMIDDLEWARE)

	} else if update.ChannelPost != nil {
		dp.ChannelPostHandler.TriggerMiddleware(dp.Bot, update, PREMIDDLEWARE)
		for _, h := range dp.ChannelPostHandler.handlers {
			cb, ok := h.Callback.(func(*Bot, *objects.Message))
			if !ok {
				return tgpErr.New("channelPost handler type assertion error, need type func(*Bot, *ChannelPost), current type is - " + fmt.Sprintln(reflect.TypeOf(h.Callback)))
			}
			err = dp.ChannelPostHandler.TriggerMiddleware(dp.Bot, update, PROCESSMIDDLEWARE)
			if err != nil {
				dp.Bot.logger.Println(err)
				continue
			}

			h.Call(update, func() { cb(dp.Bot, update.ChannelPost) })
		}
		dp.ChannelPostHandler.TriggerMiddleware(dp.Bot, update, POSTMIDDLEWARE)

	} else if update.Poll != nil {
		dp.PollHandler.TriggerMiddleware(dp.Bot, update, PREMIDDLEWARE)
		for _, h := range dp.PollHandler.handlers {
			cb, ok := h.Callback.(func(*Bot, *objects.Poll))
			if !ok {
				return tgpErr.New("poll handler type assertion error, need type func(*Bot, *Poll), current type is - " + fmt.Sprintln(reflect.TypeOf(h.Callback)))
			}
			err = dp.PollHandler.TriggerMiddleware(dp.Bot, update, PROCESSMIDDLEWARE)
			if err != nil {
				dp.Bot.logger.Println(err)
				continue
			}

			h.Call(update, func() { cb(dp.Bot, update.Poll) })
		}
		dp.PollHandler.TriggerMiddleware(dp.Bot, update, POSTMIDDLEWARE)

	} else if update.PollAnswer != nil {
		dp.PollAnswerHandler.TriggerMiddleware(dp.Bot, update, PREMIDDLEWARE)
		for _, h := range dp.PollAnswerHandler.handlers {
			cb, ok := h.Callback.(func(*Bot, *objects.PollAnswer))
			if !ok {
				return tgpErr.New("pollAnswer handler type assertion error, need type func(*Bot, *PollAnswer), current type is - " + fmt.Sprintln(reflect.TypeOf(h.Callback)))
			}
			err = dp.PollAnswerHandler.TriggerMiddleware(dp.Bot, update, PROCESSMIDDLEWARE)
			if err != nil {
				dp.Bot.logger.Println(err)
				continue
			}

			h.Call(update, func() { cb(dp.Bot, update.PollAnswer) })
		}
		dp.PollAnswerHandler.TriggerMiddleware(dp.Bot, update, POSTMIDDLEWARE)

	} else if update.ChatMember != nil {
		dp.ChatMemberHandler.TriggerMiddleware(dp.Bot, update, PREMIDDLEWARE)
		for _, h := range dp.ChatMemberHandler.handlers {
			cb, ok := h.Callback.(func(*Bot, *objects.ChatMember))
			if !ok {
				return tgpErr.New("ChatMember handler type assertion error, need type func(*Bot, *ChatMember), current type is - " + fmt.Sprintln(reflect.TypeOf(h.Callback)))
			}
			err = dp.ChatMemberHandler.TriggerMiddleware(dp.Bot, update, PROCESSMIDDLEWARE)
			if err != nil {
				dp.Bot.logger.Println(err)
				continue
			}

			h.Call(update, func() { cb(dp.Bot, update.ChatMember) })
		}
		dp.ChatMemberHandler.TriggerMiddleware(dp.Bot, update, POSTMIDDLEWARE)

	} else if update.MyChatMember != nil {

		dp.MyChatMemberHandler.TriggerMiddleware(dp.Bot, update, PREMIDDLEWARE)
		for _, h := range dp.MyChatMemberHandler.handlers {
			cb, ok := h.Callback.(func(*Bot, *objects.ChatMemberUpdated))
			if !ok {
				return tgpErr.New("MyChatMember handler type assertion error, need type func(*Bot, *ChatMemberUpdated), current type is - " + fmt.Sprintln(reflect.TypeOf(h.Callback)))
			}
			err = dp.MyChatMemberHandler.TriggerMiddleware(dp.Bot, update, PROCESSMIDDLEWARE)
			if err != nil {
				dp.Bot.logger.Println(err)
				continue
			}

			h.Call(update, func() { cb(dp.Bot, update.MyChatMember) })
		}
		dp.MyChatMemberHandler.TriggerMiddleware(dp.Bot, update, POSTMIDDLEWARE)

	} else {
		text := "detected not supported type of updates, seems like telegram bot api updated before this package updated"
		return tgpErr.New(text)
	}

	// end of adventure
	return nil
}

// SkipUpdates skip comming updates, sending to telegram servers
func (dp *Dispatcher) SkipUpdates() {
	dp.Bot.GetUpdates(&GetUpdatesConfig{
		Offset:  -1,
		Timeout: 1,
	})
}

// SetState set a state which passed for a current user in current chat
// works only in handler, or in middleware, nor outside
func (dp *Dispatcher) SetState(state *fsm.State) error {
	u := dp.currentUpdate
	if u != nil {
		cid, uid := getUidAndCidFromUpd(u)
		return dp.Storage.SetState(cid, uid, state.GetFullState())
	}
	return nil
}

// ResetState reset state for current user, and current chat
func (dp *Dispatcher) ResetState() error {
	if dp.currentUpdate != nil {
		cid, uid := getUidAndCidFromUpd(dp.currentUpdate)
		return dp.Storage.SetState(cid, uid, fsm.DefaultState.GetFullState())
	}
	return nil
}

// ========================================
//   Startup and Shutdown related methods
// ========================================

// Shutdown calls when you enter ^C(which means SIGINT)
// And SafeExit catch it, before you exit
func (dp *Dispatcher) shutdownPolling() {
	callListFuncs(dp.OnPollingShutdown, dp)
}

// startUpPolling function, iterate over a callbacks from OnStartupCallbacks
// Calls in StartPolling function
func (dp *Dispatcher) startupPolling() {
	callListFuncs(dp.OnPollingStartup, dp)
	dp.welcome()
}

// shutdownWebhook method, iterate over a callbacks from OnWebhookShutdown
func (dp *Dispatcher) shutdownWebhook() {
	callListFuncs(dp.OnWebhookShutdown, dp)
}

// startupPolling method, iterate over a callbacks from OnWebhookStartup
func (dp *Dispatcher) startupWebhook() {
	callListFuncs(dp.OnWebhookStartup, dp)
	dp.welcome()
}

// Onstartup method append to OnStartupCallbaks a callbacks
// Using pointers bc cant unregister function using copy of object
// And golang doesnot support generics, and type equals
func (dp *Dispatcher) OnStartup(c *OnConfig) {
	if !c.Webhook && !c.Polling {
		dp.Bot.logger.Println("this expression have not got any effect")
	}

	if c.Webhook {
		dp.OnWebhookStartup = append(dp.OnWebhookStartup, c.cb...)
	}
	if c.Polling {
		dp.OnPollingStartup = append(dp.OnPollingStartup, c.cb...)
	}
}

// OnShutdown method using for register OnShutdown callbacks
// Same code like OnStartup
func (dp *Dispatcher) OnShutdown(c *OnConfig) {
	if !c.Webhook && !c.Polling {
		dp.Bot.logger.Println("!polling and !webhook expression have not got any effect")
	}

	if c.Webhook {
		dp.OnWebhookShutdown = append(dp.OnWebhookShutdown, c.cb...)
	}
	if c.Polling {
		dp.OnPollingShutdown = append(dp.OnPollingShutdown, c.cb...)
	}
}

func (dp *Dispatcher) Start() {
	if dp.polling {
		dp.startupPolling()
	}
	if dp.webhook {
		dp.startupWebhook()
	}
}

func (dp *Dispatcher) Shutdown() {
	if dp.polling {
		dp.shutdownPolling()
	}
	if dp.webhook {
		dp.shutdownWebhook()
	}
}

// SafeExit method uses for notify about exit from program
// Thanks: https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
func (dp *Dispatcher) SafeExit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		dp.shutDown()
		os.Exit(0)
	}()
}

// ShutDownDP calls ResetWebhook for reset webhook in telegram servers, if yes
func (dp *Dispatcher) shutDown() {
	dp.Bot.logger.Println("Stop polling!")
	dp.ResetWebhook(true)
	dp.Storage.Close()
	dp.Shutdown()
}

func (dp *Dispatcher) welcome() {
	if dp.Welcome {
		dp.Bot.GetMe()
		dp.Bot.logger.Println("Bot: ", dp.Bot.Me)
	}
}

// =========================================
//    Polling and webhook related methods
// =========================================

// GetUpdatesChan makes getUpdates request to telegram servers
// sends update to updates channel
// Time.Sleep here for stop goroutine for a c.Relax time
//
// yeah it bad, and works only on crutches, but works
func (dp *Dispatcher) MakeUpdatesChan(c *StartPollingConfig, ch chan *objects.Update) {
	go func() {
		for {
			if c.Relax != 0 {
				time.Sleep(c.Relax)
			}

			updates, err := dp.Bot.GetUpdates(c.GetUpdatesConfig)
			if err != nil {
				dp.Bot.logger.Println(err.Error())
				dp.Bot.logger.Println("Error with getting updates")
				time.Sleep(time.Duration(c.ErrorSleep))

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= c.Offset {
					c.Offset = update.UpdateID + 1
					ch <- update
				}
			}
		}
	}()
}

// ProcessUpdates iterate <-chan *objects.Update
//
// Note: use after a MakeUpdatesChan call
func (dp *Dispatcher) ProcessUpdates(ch <-chan *objects.Update) error {
	for upd := range ch {
		if upd == nil {
			continue
		}
		dp.currentUpdate = upd
		err := dp.ProcessOneUpdate(upd)
		if err != nil {
			return err
		}
	}

	return nil
}

// StartPolling check out to comming updates
// If yes, Telegram Get to your bot a Update
// Using GetUpdates method in Bot structure
// GetUpdates config using for getUpdates method
func (dp *Dispatcher) StartPolling(c *StartPollingConfig) error {
	if dp.webhook {
		panic(ErrorConflictModes)
	}

	dp.polling = true
	dp.Start()
	if c.SafeExit {
		dp.SafeExit()
	}
	if c.ResetWebhook {
		dp.ResetWebhook(true)
	}

	if c.SkipUpdates {
		dp.SkipUpdates()
	}

	ch := make(chan *objects.Update)

	dp.MakeUpdatesChan(c, ch)
	dp.ProcessUpdates(ch)

	return nil
}

// MakeWebhookChan adds a http Handler with c.BotURL path
func (dp *Dispatcher) MakeWebhookChan(c *StartWebhookConfig, ch chan *objects.Update) {
	http.HandleFunc(c.BotURL, func(wr http.ResponseWriter, req *http.Request) {
		update, err := requestToUpdate(req)
		if err != nil {
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			wr.WriteHeader(http.StatusBadRequest)
			wr.Header().Set("Content-Type", "application/json")
			wr.Write(errMsg)
			return
		}

		ch <- update
	})
}

// StartWebhook method registers BotUrl uri a function which handles every comming update
// Using In Pair of SetWebhook method
// Startup method executes after SetWebhook method call
//
// NOTE: you should to add a webhook close callback function, using OnShutdown
func (dp *Dispatcher) StartWebhook(c *StartWebhookConfig) error {
	if dp.polling {
		panic(ErrorConflictModes)
	}
	_, err := dp.Bot.SetWebhook(c.SetWebhookConfig)
	if err != nil {
		return err
	}
	dp.webhook = true
	dp.Start()
	if c.SafeExit {
		dp.SafeExit()
	}
	http.HandleFunc(c.BotURL, func(wr http.ResponseWriter, req *http.Request) {
		update, err := requestToUpdate(req)
		if err != nil {
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			wr.WriteHeader(http.StatusBadRequest)
			wr.Header().Set("Content-Type", "application/json")
			wr.Write(errMsg)
			return
		}

		dp.currentUpdate = update
		err = dp.ProcessOneUpdate(update)
		if err != nil {
			fmt.Println(err)
		}
	})
	certPath, err := guessFileName(c.Certificate)
	if err != nil {
		return err
	}
	keyfile, err := guessFileName(c.KeyFile)
	if err != nil {
		return err
	}
	return http.ListenAndServeTLS(c.Address, certPath, keyfile, c.Handler)
}
