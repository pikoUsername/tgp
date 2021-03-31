package bot

import (
	"encoding/json"
	"net/url"

	"github.com/pikoUsername/tgp/tgp/configs"
	"github.com/pikoUsername/tgp/tgp/objects"
	"github.com/pikoUsername/tgp/tgp/utils"
)

// Bot ...
type Bot struct {
	Token string
	ParseMode string

	// Using prefix Bot, for avoid names conflict
	// and golang dont love name conflicts
	BotParseMode string
	Me           *objects.User `json:"-"`
}

// NewBot get a new Bot
// This Fucntion checks a token
// for spaces and etc.
func NewBot(token string, checkToken bool, parseMode string) (*Bot, error) {
	if checkToken {
		err := utils.CheckToken(token)
		if err != nil {
			return nil, err
		}
	}
	return &Bot{
		Token: token,
		ParseMode: parseMode,
	}, nil
}

// GetMe reporesents telegram method
// https://core.telegram.org/bots/api#getme
func (bot *Bot) GetMe() (*objects.User, error) {
	if bot.Me != nil  {
		return bot.Me, nil
	}
	resp, err := MakeRequest("getMe", bot.Token, &url.Values{})
	if err != nil {
		return &objects.User{}, err
	}
	var user objects.User
	err = json.Unmarshal(resp.Result, &user)
	if err != nil {
		return &user, err
	}

	bot.Me = &user
	return &user, nil
}

// Logout your bot from telegram
// https://core.telegram.org/bots/api#logout
func (bot *Bot) Logout() (bool, error) {
	_, err := MakeRequest("logout", bot.Token, &url.Values{})
	if err != nil {
		return false, err
	}
	return true, nil
} // Indeed

// Send uses as sender for almost all stuff
func (bot *Bot) SendMessageable(c configs.Configurable) (*objects.Message, error) {
	v, err := c.Values()
	if err != nil {
		return &objects.Message{}, err
	}
	resp, err := MakeRequest(c.Method(), bot.Token, v)
	if err != nil {
		return &objects.Message{}, err
	}
	var msg objects.Message
	err = json.Unmarshal(resp.Result, &msg)
	if err != nil {
		return &msg, err
	}
	return &msg, nil
}

// Send ...
func (bot *Bot) Send(config configs.Configurable) (*objects.Message, error) {
	switch config.(type) {
	case configs.FileableConf:
		return nil, nil
	default:
		return bot.SendMessageable(config)
	}
}

// CopyMessage copies message
// https://core.telegram.org/bots/api#copymessage
func (bot *Bot) CopyMessage(config *configs.CopyMessageConfig) (*objects.MessageID, error) {
	// Stub here, TODO: make for every config a values function/method
	v, err := config.Values()
	resp, err := MakeRequest(config.Method(), bot.Token, v)
	if err != nil {
		return &objects.MessageID{}, err
	}
	var msg objects.MessageID

	err = json.Unmarshal(resp.Result, &msg)
	if err != nil { return &msg, err }
	return &msg, nil
}

// ----------------HardCoded for now Dont watch this place pls-----------------

// SendPhoto ...
func (bot *Bot) SendPhoto(config *configs.SendPhotoConfig) (*objects.Message, error) {
	return &objects.Message{}, nil
}

// SendAudio ...
func (bot *Bot) SendAudio(config *configs.SendAudioConfig) (*objects.Message, error) {
	return &objects.Message{}, nil
}

// SendDocument ...
func (bot *Bot) SendDocument(config *configs.SendDocumentConfig) (*objects.Message, error) {
	return &objects.Message{}, nil
}

// SendVideo ...
func (bot *Bot) SendVideo(config *configs.SendVideoConfig) (*objects.Message, error) {
	return &objects.Message{}, nil
}

// SendAnimation ...
func (bot *Bot) SendAnimation(config *configs.SendAnimation) (*objects.Message, error) {
	return &objects.Message{}, nil
}

// SendVoice ...
func (bot *Bot) SendVoice(config *configs.SendVoiceConfig) (*objects.Message, error) {
	return &objects.Message{}, nil
}

// SendVideoName ...
func (bot *Bot) SendVideoName(config *configs.SendVideoNameConfig) (*objects.Message, error) {
	return &objects.Message{}, nil
}

// SendMediaGroup ...
func (bot *Bot) SendMediaGroup(config *configs.SendMediaGroupConfig) (*objects.Message, error) {
	return &objects.Message{}, nil
}

// SendLocation ...
func (bot *Bot) SendLocation(config *configs.SendLocationConfig) (*objects.Message, error) {
	return &objects.Message{}, nil
}

// editMessageLiveLocation ...
func (bot *Bot) EditMessageLiveLocation(config *configs.LiveLocationConfig) (*objects.Message, error) {
	return &objects.Message{}, nil
}

// ------------------EndHardCode, Phew----------------------

// GetUpdates uses for long polling
// https://core.telegram.org/bots/api#getupdates
func (bot *Bot) GetUpdates(c *configs.GetUpdatesConfig) (*objects.Update, error) {
	v, err := c.Values()
	if err != nil {
		return &objects.Update{}, err
	}
	resp, err := MakeRequest(c.Method(), bot.Token, v)
	if err != nil {
		return &objects.Update{}, &objects.TelegramApiError{
			Code:               resp.ErrorCode,
			Description:        resp.Description,
			ResponseParameters: objects.ResponseParameters{},
		}
	}
	var upd objects.Update
	err = json.Unmarshal(resp.Result, &upd)
	if err != nil {
		return &upd, err
	}
	return &upd, nil
}

// SendMessage sends message using ChatID
// see: https://core.telegram.org/bots/api#sendmessage
func (bot *Bot) SendMessage(config *configs.SendMessageConfig) (*objects.Message, error) {
	v, err := config.Values()
	if err != nil {
		return &objects.Message{}, err
	}
	resp, err := MakeRequest("sendMessage", bot.Token, v)
	if err != nil {
		return &objects.Message{}, &objects.TelegramApiError{
			Code:               resp.ErrorCode,
			Description:        resp.Description,
			ResponseParameters: objects.ResponseParameters{},
		}
	}
	var msg objects.Message
	err = json.Unmarshal(resp.Result, &msg)
	if err != nil {
		return &msg, nil
	}
	return &msg, nil
}

// SetWebhook make subscribe to telegram events
// or sends to telegram a request for make
// Subscribe to specific IP, and when user
// sends a message to your bot, Telegram know
// Your bot IP and sends to your bot a Update
// https://core.telegram.org/bots/api#setwebhook
func (bot *Bot) SetWebhook(config *configs.SetWebhookConfig) error {
	v, err := config.Values()
	if err != nil {
		return err
	}
	_, err = MakeRequest(config.Method(), bot.Token, v)
	if err != nil {
		return err
	}
	return nil
}