package dispatcher

import (
	"errors"

	"github.com/pikoUsername/tgp/bot"
	"github.com/pikoUsername/tgp/configs"
	"github.com/pikoUsername/tgp/objects"
)

// Dispatcher need for Polling, and webhook
// For Bot run,
// Bot struct uses as API wrapper
// Dispatcher uses as Bot starter
// Middlewares uses function
// Another level of abstraction
type Dispatcher struct {
	Bot *bot.Bot

	// Handlers
	MessageHandler       HandlerObj
	CallbackQueryHandler HandlerObj
	UpdatesHandler       HandlerObj

	polling bool
	webhook bool
}

// NewDispathcer get a new Dispatcher
// And with autoconfiguration, need to run once
func NewDispatcher(bot *bot.Bot) (*Dispatcher, error) {
	dp := &Dispatcher{
		Bot: bot,
	}

	dp.MessageHandler = &DefaultHandlerObj{}
	dp.CallbackQueryHandler = &DefaultHandlerObj{}
	dp.UpdatesHandler = &DefaultHandlerObj{}

	return dp, nil
}

func (dp *Dispatcher) ResetWebhook(check bool) error {
	if check {
		wi, err := dp.Bot.GetWebhookInfo()
		if err != nil {
			return err
		}
		if wi.URL == "" {
			return nil
		}
	}
	return dp.Bot.DeleteWebhook(&configs.DeleteWebhookConfig{})
}

// RegisterMessageHandler excepts you pass to parametrs a your function
func (dp *Dispatcher) RegisterMessageHandler(callback HandlerType) {
	dp.MessageHandler.Register(callback)
}

// ProcessUpdates havenot got any efficient
// if you use webhook and long polling
func (dp *Dispatcher) ProcessPollingUpdates(updates []objects.Update) error {
	return nil // TODO
}

// ProcessUpdates using for process updates from any way
func (dp *Dispatcher) ProcessUpdates(updates []objects.Update) error {
	for _, upd := range updates {
		err := dp.ProcessOneUpdate(upd)
		if err != nil {
			return err
		}
	}

	return nil
}

// ProcessOneUpdate you guess, processes ONLY one comming update
// Support only one Message update
func (dp *Dispatcher) ProcessOneUpdate(update objects.Update) error {
	if update.Message != nil {
		event := update.Message
		return dp.MessageHandler.Trigger(event, *dp.Bot)
	} else if update.CallbackQuery != nil {
		event := update.CallbackQuery
		return dp.MessageHandler.Trigger(event, *dp.Bot)
	} else {
		text := "Detected Not supported type of updates Seems like Telegram bot api updated brfore this package updated"
		return errors.New(text)
	}
}

// StartPolling check out to comming updates
// If yes, Telegram Get to your bot a Update
// Using GetUpdates method in Bot structure
// GetUpdates config using for getUpdates method
func (dp *Dispatcher) StartPolling(c *configs.GetUpdatesConfig) error {
	for {
		// TODO: timeout
		updates, err := dp.Bot.GetUpdates(c)
		if err != nil {
			return err
		}
		if updates != nil {
			// I cant understand how it s works, and where need to use goroutines
			err := dp.ProcessPollingUpdates(updates)
			if err != nil {
				return err
			}
		}
	}
}

// StartWebhook ...
func (dp *Dispatcher) StartWebhook() error {
	return nil
}
