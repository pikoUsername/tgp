package tgp

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
	"unsafe"

	"github.com/pikoUsername/tgp/objects"
)

// This file stores ALL method configs

// functions with New prefix, use with Context.Reply method

// Configurable is interface for using by method
type Configurable interface {
	values() (url.Values, error)
	method() string
}

// FileableConf config using for Files storing
type FileableConf interface {
	Configurable
	params() (map[string]string, error)
	getFiles() []*objects.InputFile
}

// BaseChat taken from go-telegram-bot-api
type BaseChat struct {
	ChatID              int64
	ChannelUsername     string
	ReplyToMessageID    int64
	ReplyMarkup         interface{}
	DisableNotification bool
}

// helper method
func (c *BaseChat) params() (map[string]string, error) {
	params := make(map[string]string)

	if c.ChannelUsername != "" {
		params["chat_id"] = strconv.FormatInt(c.ChatID, 10)
	} else {
		params["chat_id"] = c.ChannelUsername
	}
	if c.ReplyToMessageID != 0 {
		params["reply_to_message_id"] = strconv.FormatInt(c.ChatID, 10)
	}

	if c.ReplyMarkup != nil {
		params["reply_markup"] = FormatMarkup(c.ReplyMarkup)
	}
	params["disable_notification"] = strconv.FormatBool(c.DisableNotification)

	return params, nil
}

// values returns url.Values representation of BaseChat
func (c *BaseChat) values() (url.Values, error) {
	v := url.Values{}
	if c.ChannelUsername != "" {
		v.Add("chat_id", c.ChannelUsername)
	} else {
		v.Add("chat_id", strconv.FormatInt(c.ChatID, 10))
	}

	if c.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(c.ReplyToMessageID, 10))
	}

	if c.ReplyMarkup != nil {
		v.Add("reply_markup", FormatMarkup(c.ReplyMarkup))
	}

	v.Add("disable_notification", strconv.FormatBool(c.DisableNotification))

	return v, nil
}

// BaseFile taken from go-telegram-bot-api
type BaseFile struct {
	BaseChat
	File        *objects.InputFile
	FileID      string
	UseExisting bool
	MimeType    string
	FileSize    int
}

// params ...
func (bf *BaseFile) params() (v map[string]string, err error) {
	v = make(map[string]string)

	if bf.FileID != "" {
		v["file_id"] = bf.FileID
	}
	v["use_existing"] = strconv.FormatBool(bf.UseExisting)
	if bf.MimeType != "" {
		v["mime_type"] = bf.MimeType
	}
	if bf.FileSize != 0 {
		v["file_size"] = strconv.Itoa(bf.FileSize)
	}

	cv, _ := bf.values()
	urlValuesToMapString(cv, v)

	return v, nil
}

type UserPermissionsConfig struct {
	CanManageChat       bool `json:"can_manage_chat"`
	CanPostMessage      bool `json:"can_post_message"`
	CanEditMessages     bool `json:"can_edit_messages"`
	CanDeleteMessages   bool `json:"can_delete_messages"`
	CanManageVideoChats bool `json:"can_manage_video_chats"`
	CanRestrictMembers  bool `json:"can_restrict_members"`
	CanPromoteMembers   bool `json:"can_promote_members"`
	CanChangeInfo       bool `json:"can_change_info"`
	CanInviteUsers      bool `json:"can_invite_users"`
	CanPinMessages      bool `json:"can_pin_messages"`
}

func (upc *UserPermissionsConfig) values() (url.Values, error) {
	v := url.Values{}

	// TODO:
	// f := func(n string, b bool, iv url.Values) {
	// 	iv.Add(n, strconv.FormatBool(b))
	// }

	return v, nil
}

// For CopyMessage method config
// https://core.telegram.org/bots/api#copymessage
type CopyMessageConfig struct {
	DisableNotifications  bool
	AllowSendingWithReply bool
	Caption               string
	ChatID                int64 // required
	FromChatID            int64 // required
	MessageID             int64 // required
	ReplyToMessageId      int64
	CaptionEntities       []*objects.MessageEntity
	ProtectContent        bool

	// type: Union[objects.InlineKeyboardMarkup, ReplyKeyboardMarkup]
	ReplyMarkup interface{}
}

func (cmc *CopyMessageConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(cmc.ChatID, 10))
	v.Add("from_chat_id", strconv.FormatInt(cmc.ChatID, 10))
	v.Add("protect_content", strconv.FormatBool(cmc.ProtectContent))
	v.Add("message_id", strconv.FormatInt(cmc.MessageID, 10))
	if cmc.Caption != "" {
		v.Add("caption", cmc.Caption)
	}
	if cmc.CaptionEntities != nil {
		v.Add("caption_entities", ObjectToJson(cmc.CaptionEntities))
	}
	v.Add("disable_notifications", strconv.FormatBool(cmc.DisableNotifications))
	if cmc.ReplyToMessageId != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(cmc.ReplyToMessageId, 10))
	}
	v.Add("allow_sending_with_reply", strconv.FormatBool(cmc.AllowSendingWithReply))
	if cmc.ReplyMarkup != nil {
		v.Add("reply_keyboards", FormatMarkup(cmc.ReplyMarkup))
	}
	return v, nil
}

func (cmc *CopyMessageConfig) method() string {
	return "copyMessage"
}

// SendMessageConfig respresnests method,
// and fields of sendMessage method of telegram
// https://core.telegram.org/bots/api#sendmessage
type SendMessageConfig struct {
	// Required Field
	ChatID int64

	// It s too, Telegram excepts
	Text                  string // required
	ParseMode             string
	Entities              []*objects.MessageEntity
	DisableWebPagePreview bool

	DisableNotifiaction bool
	ReplyKeyboard       *objects.InlineKeyboardMarkup
	ProtectContent      bool
}

// values ...
func (smc *SendMessageConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(smc.ChatID, 10))

	v.Add("text", smc.Text)

	if smc.ParseMode != "" {
		v.Add("parse_mode", smc.ParseMode)
	}

	if smc.ReplyKeyboard != nil {
		v.Add("reply_markup", FormatMarkup(smc.ReplyKeyboard))
	}
	v.Add("disable_web_page_preview", strconv.FormatBool(smc.DisableWebPagePreview))
	v.Add("protect_content", strconv.FormatBool(smc.ProtectContent))
	if smc.Entities != nil {
		// Must be work!
		v.Add("entities", ObjectToJson(smc.Entities))
	}

	return v, nil
}

func (smc *SendMessageConfig) method() string {
	return "sendMessage"
}

func NewSendMessage(text string, chat_id int64) *SendMessageConfig {
	return &SendMessageConfig{
		Text:   text,
		ChatID: chat_id,
	}
}

func NewReplyMessage(text string) *SendMessageConfig {
	return &SendMessageConfig{
		Text: text,
	}
}

// SetWebhookConfig uses for Using as arguemnt
// You may not fill all fields in struct
// https://core.telegram.org/bots/api#setwebhook
type SetWebhookConfig struct {
	URL                string // required
	Offset             int
	MaxConnections     int
	AllowedUpdates     []string
	DropPendingUpdates bool
	IP                 string // if you need u can use it ;)
	Certificate        *objects.InputFile
}

func (wc *SetWebhookConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("url", wc.URL)
	if wc.IP != "" {
		v.Add("ip_address", wc.IP) // required field
	}
	if wc.MaxConnections != 0 {
		v.Add("max_connections", strconv.Itoa(wc.MaxConnections))
	}
	if len(wc.AllowedUpdates) != 0 {
		bs, err := json.Marshal(wc.AllowedUpdates)
		if err != nil {
			return nil, err
		}
		v.Add("allowed_updates", BytesToString(bs))
	}
	v.Add("drop_pending_updates", strconv.FormatBool(wc.DropPendingUpdates))

	return v, nil
}

func (wc *SetWebhookConfig) method() string {
	return "setWebhook"
}

// URL must be full URL to your domain
func NewSetWebhook(url string) *SetWebhookConfig {
	return &SetWebhookConfig{
		URL: url,
	}
}

// SendPhotoConfig represnts telegram api method fields
// https://core.telegram.org/bots/api#sendphoto
type SendPhotoConfig struct {
	*BaseFile
	Caption        string
	ProtectContent bool
}

func (spc *SendPhotoConfig) values() (url.Values, error) {
	v, _ := spc.BaseFile.values()
	if spc.Caption != "" {
		v.Add("caption", spc.Caption)
	}
	v.Add("protect_content", strconv.FormatBool(spc.ProtectContent))
	return v, nil
}

func (spc *SendPhotoConfig) method() string {
	return "sendPhoto"
}

func (spc *SendPhotoConfig) params() (map[string]string, error) {
	v, _ := spc.BaseFile.params()
	if spc.Caption != "" {
		v["caption"] = spc.Caption
	}
	return v, nil
}

func (spc *SendPhotoConfig) getFiles() []*objects.InputFile {
	return []*objects.InputFile{spc.File}
}

func NewSendPhoto(photo *objects.InputFile) *SendPhotoConfig {
	return &SendPhotoConfig{
		BaseFile: &BaseFile{
			BaseChat: BaseChat{},
			File:     photo,
		},
	}
}

// represents a sendAudio fields
type SendAudioConfig struct {
	BaseFile
	Caption         string
	ParseMode       string
	Duration        uint
	Performer       string
	Title           string
	ProtectContent  bool
	Thumb           *objects.InputFile
	CaptionEntities []*objects.MessageEntity
}

func (sac *SendAudioConfig) values() (url.Values, error) {
	v, _ := sac.BaseFile.values()

	v.Add("chat_id", strconv.FormatInt(sac.ChatID, 10))

	if sac.Caption != "" {
		v.Add("caption", sac.Caption)
		if sac.ParseMode != "" {
			v.Add("parse_mode", sac.ParseMode)
		}
		if sac.CaptionEntities != nil {
			v.Add("caption_entities", ObjectToJson(sac.CaptionEntities))
		}
	}
	if sac.Duration != 0 {
		v.Add("duration", strconv.FormatUint(uint64(sac.Duration), 10))
	}
	if sac.Performer != "" {
		v.Add("performer", sac.Performer)
	}
	if sac.Title != "" {
		v.Add("title", sac.Title)
	}
	v.Add("protect_content", strconv.FormatBool(sac.ProtectContent))
	return v, nil
}

func (sac *SendAudioConfig) params() (map[string]string, error) {
	v, _ := sac.values()
	m := make(map[string]string)

	urlValuesToMapString(v, m)

	return m, nil
}

func (sac *SendAudioConfig) method() string {
	return "sendAudio"
}

func (sac *SendAudioConfig) getFiles() []*objects.InputFile {
	return []*objects.InputFile{sac.File, sac.Thumb}
}

func NewSendAudio(audio *objects.InputFile) *SendAudioConfig {
	return &SendAudioConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{},
			File:        audio,
			UseExisting: false,
		},
	}
}

// SendDocumentConfig represents sendDoucument method fields
type SendDocumentConfig struct {
	ChatID                      int64              // required
	Document                    *objects.InputFile // required
	Thumb                       *objects.InputFile
	Caption                     string
	ParseMode                   string
	CaptionEntities             []*objects.MessageEntity
	DisableContentTypeDetection bool
	DisableNotifiaction         bool
	ProtectContent              bool
	ReplyToMessageID            int64
	AllowSendingWithoutReply    bool
	ReplyMarkup                 interface{}
}

func (sdc *SendDocumentConfig) values() (v url.Values, err error) {
	v = url.Values{}

	v.Add("chat_id", strconv.FormatInt(sdc.ChatID, 10))
	if sdc.Caption != "" {
		v.Add("caption", sdc.Caption)
		if sdc.ParseMode != "" {
			v.Add("parse_mode", sdc.ParseMode)
		}
		if sdc.CaptionEntities != nil {
			v.Add("caption_entities", ObjectToJson(sdc.CaptionEntities))
		}
	}
	v.Add("disable_notification", strconv.FormatBool(sdc.DisableNotifiaction))
	v.Add("protect_content", strconv.FormatBool(sdc.ProtectContent))
	if sdc.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(sdc.ReplyToMessageID, 10))
	}
	v.Add("allow_sending_without_reply", strconv.FormatBool(sdc.AllowSendingWithoutReply))
	if sdc.ReplyMarkup != nil {
		v.Add("reply_markup", FormatMarkup(sdc.ReplyMarkup))
	}
	return v, nil
}

func (sdc *SendDocumentConfig) getFiles() []*objects.InputFile {
	return []*objects.InputFile{sdc.Document, sdc.Thumb}
}

func (sdc *SendDocumentConfig) params() (map[string]string, error) {
	params := make(map[string]string)

	v, _ := sdc.values()
	urlValuesToMapString(v, params)

	return params, nil
}

func (sdc *SendDocumentConfig) method() string {
	return "sendDocument"
}

func NewDocumentConfig(cid int64, r *objects.InputFile) *SendDocumentConfig {
	return &SendDocumentConfig{
		ChatID:   cid,
		Document: r,
	}
}

// SendVideoConfig Represents sendVideo fields
// https://core.telegram.org/bots/api#sendvideo
type SendVideoConfig struct {
	*BaseFile
	Duration       uint32
	Width          uint16
	Height         uint16
	ProtectContent bool
	Thumb          *objects.InputFile
}

func (svc *SendVideoConfig) values() (url.Values, error) {
	v, _ := svc.BaseFile.values()
	if svc.Duration != 0 {
		v.Add("duration", strconv.FormatUint((uint64)(svc.Duration), 10))
	}

	if svc.Width != 0 {
		v.Add("width", strconv.FormatUint((uint64)(svc.Width), 10))
	}
	if svc.Height != 0 {
		v.Add("height", strconv.FormatUint((uint64)(svc.Height), 10))
	}
	v.Add("protect_content", strconv.FormatBool(svc.ProtectContent))

	return v, nil
}

func (svc *SendVideoConfig) params() (map[string]string, error) {
	v, _ := svc.BaseFile.params()
	uv, _ := svc.values()

	urlValuesToMapString(uv, v)

	return v, nil
}

func (svc *SendVideoConfig) getFiles() []*objects.InputFile {
	return []*objects.InputFile{svc.File, svc.Thumb}
}

func (svc *SendVideoConfig) method() string {
	return "sendVideo"
}

// Represents Method SendAnimation Fields
// https://core.telegram.org/bots/api#sendanimation
type SendAnimationConfig struct {
	ChatID    int64
	Animation *objects.InputFile

	Duration uint32
	Width    uint32
	Height   uint32

	Thumb     *objects.InputFile
	Caption   string
	ParseMode string

	ProtectContent bool
}

func (sac *SendAnimationConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(sac.ChatID, 10))
	v.Add("duration", strconv.FormatUint(uint64(sac.Duration), 10))
	v.Add("width", strconv.FormatUint(uint64(sac.Width), 10))
	v.Add("height", strconv.FormatUint(uint64(sac.Height), 10))
	if sac.Caption != "" {
		v.Add("caption", sac.Caption)
	}
	if sac.ParseMode != "" {
		v.Add("parse_mode", sac.ParseMode)
	}
	v.Add("protect_content", strconv.FormatBool(sac.ProtectContent))

	return v, nil
}

func (sac *SendAnimationConfig) method() string {
	return "sendAnimation"
}

func (sac *SendAnimationConfig) params() (map[string]string, error) {
	m := map[string]string{}
	v, _ := sac.values()
	urlValuesToMapString(v, m)
	return m, nil
}

func (sac *SendAnimationConfig) getFiles() []*objects.InputFile {
	return []*objects.InputFile{sac.Animation, sac.Thumb}
}

func NewSendAnimtion(chat_id int64, animation *objects.InputFile) *SendAnimationConfig {
	return &SendAnimationConfig{
		ChatID:    chat_id,
		Animation: animation,
	}
}

type SendVoiceConfig struct {
	*BaseFile
	ChatID               int64
	Caption              string
	ParseMode            string
	CaptionEntities      []*objects.MessageEntity
	Duration             int
	DisableNotifications bool
	ReplyToMessageID     int64
	ProtectContent       bool

	// for first time you can use InlineKeyboardMarkup
	ReplyMarkup *objects.InlineKeyboardMarkup
}

func (svc *SendVoiceConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(svc.ChatID, 10))
	v.Add("caption", svc.Caption)
	if svc.Caption != "" {
		v.Add("parse_mode", svc.Caption)
	}
	v.Add("disable_notifications", strconv.FormatBool(svc.DisableNotifications))
	if svc.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(svc.ReplyToMessageID, 10))
	}

	v.Add("caption_entities", ObjectToJson(svc.CaptionEntities))
	v.Add("protect_content", strconv.FormatBool(svc.ProtectContent))
	if svc.ReplyMarkup != nil {
		v.Add("reply_markup", FormatMarkup(svc.ReplyMarkup))
	}

	return v, nil
}

func (s *SendVoiceConfig) getFiles() []*objects.InputFile {
	return []*objects.InputFile{s.File}
}

func (svc *SendVoiceConfig) method() string {
	return "sendVoice"
}

type SendVideoNoteConfig struct {
	*BaseFile
	Duration                 time.Duration
	Length                   int64
	Thumb                    *objects.InputFile
	AllowSendingWithoutReply bool
	ProtectContent           bool
}

func (svnc *SendVideoNoteConfig) values() (url.Values, error) {
	v, _ := svnc.BaseFile.values()
	if svnc.Duration != 0 {
		v.Add("duration", strconv.FormatInt(*(*int64)(unsafe.Pointer(&svnc.Duration)), 10))
	}
	if svnc.Length != 0 {
		v.Add("length", strconv.FormatInt(svnc.Length, 10))
	}
	v.Add("allow_sending_without_reply", strconv.FormatBool(svnc.AllowSendingWithoutReply))
	v.Add("protect_content", strconv.FormatBool(svnc.ProtectContent))

	return v, nil
}

func (svnc *SendVideoNoteConfig) params() (v map[string]string, err error) {
	v, _ = svnc.BaseFile.params()
	uv, _ := svnc.values()
	urlValuesToMapString(uv, v)

	return
}

func (svnc *SendVideoNoteConfig) getFiles() []*objects.InputFile {
	return []*objects.InputFile{svnc.File, svnc.Thumb}
}
func (svnc *SendVideoNoteConfig) method() string {
	return "sendVideoName"
}

func NewSendVideoNote(video_note *objects.InputFile) *SendVideoNoteConfig {
	return &SendVideoNoteConfig{
		BaseFile: &BaseFile{
			BaseChat: BaseChat{},
			File:     video_note,
		},
	}
}

type SendMediaGroupConfig struct {
	// required fields
	ChatID int64
	Media  []interface{} // type: Union[[]InputMediaAudio, []InputMediaDocument, []InputMediaPhoto, []InputMediaVideo]

	// Optional fields
	DisableNotification      bool
	ProtectContent           bool
	ReplyToMessageID         int64
	AllowSendingWithoutReply bool
}

func (smgc *SendMediaGroupConfig) values() (url.Values, error) {
	v := url.Values{}

	v.Add("protect_sending", strconv.FormatBool(smgc.ProtectContent))
	v.Add("chat_id", strconv.FormatInt(smgc.ChatID, 10))
	// TOOD: media types
	// v.Add("media", smgc.Media)
	v.Add("disable_notification", strconv.FormatBool(smgc.DisableNotification))

	if smgc.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(smgc.ReplyToMessageID, 10))
	}

	v.Add("allow_sending_without_reply", strconv.FormatBool(smgc.AllowSendingWithoutReply))

	return v, nil
}

func (smgc *SendMediaGroupConfig) method() string {
	return "sendMediaGroup"
}

func NewSendMediaGroupConfig(media []interface{}) *SendMediaGroupConfig {
	return &SendMediaGroupConfig{
		Media: media,
	}
}

type SendLocationConfig struct {
	ChatID                   int64   // req
	Latitude                 float32 // req
	Longitude                float32 // req
	HorizontalAccuracy       float32
	LivePeriod               uint
	Heading                  int
	ProximityAlertRadius     int
	DisableNotification      bool
	ReplyToMessageID         int
	AllowSendingWithoutReply bool
	ProtectContent           bool
}

func (slc *SendLocationConfig) values() (url.Values, error) {
	v := url.Values{}

	v.Add("chat_id", strconv.FormatInt(slc.ChatID, 10))

	// Same lines, broken DRY
	v.Add("latitude", strconv.FormatFloat(float64(slc.Latitude), 'E', -1, 64))
	v.Add("longitude", strconv.FormatFloat(float64(slc.Longitude), 'E', -1, 64))
	v.Add("horizontal_accuracy", strconv.FormatFloat(float64(slc.HorizontalAccuracy), 'E', -1, 64))

	if slc.LivePeriod != 0 {
		v.Add("live_period", strconv.FormatUint(uint64(slc.LivePeriod), 10))
	}

	if slc.Heading != 0 {
		v.Add("heading", strconv.FormatInt(int64(slc.Heading), 10))
	}

	v.Add("proximity_alert_radius", strconv.FormatInt(int64(slc.ProximityAlertRadius), 10))
	v.Add("disable_notification", strconv.FormatBool(slc.DisableNotification))
	v.Add("protect_content", strconv.FormatBool(slc.ProtectContent))

	if slc.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(slc.ReplyToMessageID))
	}

	v.Add("allow_sending_without_reply", strconv.FormatBool(slc.AllowSendingWithoutReply))

	return v, nil
}

func NewSendLocationConf(latitude float32, longitude float32) *SendLocationConfig {
	return &SendLocationConfig{
		Latitude:  latitude,
		Longitude: longitude,
	}
}

func (slc *SendLocationConfig) method() string {
	return "sendLocation"
}

// LiveLocationConfig represents Telegram method fields of editmessageliveLocation
// https://core.telegram.org/bots/api#editmessagelivelocation
type EditMessageLLConf struct { // too long name anyway
	Longitude            float64 // required
	Latitude             float64 // required
	InlineMessageID      int64
	ChatID               int64
	MessageID            int64
	HorizontalAccuracy   float64
	Heading              int64
	ProximityAlertRadius int64
	ReplyMarkup          *objects.InlineKeyboardMarkup
}

func (llc *EditMessageLLConf) values() (url.Values, error) {
	v := url.Values{}

	v.Add("longitude", strconv.FormatFloat(llc.Longitude, 'E', -1, 64))
	v.Add("latitude", strconv.FormatFloat(llc.Latitude, 'E', -1, 64))
	if llc.InlineMessageID != 0 {
		v.Add("inline_message_id", strconv.FormatInt(llc.InlineMessageID, 10))
	}
	if llc.ChatID != 0 {
		v.Add("chat_id", strconv.FormatInt(llc.ChatID, 10))
	}
	if llc.MessageID != 0 {
		v.Add("message_id", strconv.FormatInt(llc.MessageID, 10))
	}
	if llc.HorizontalAccuracy != 0.0 {
		v.Add("horizontal_accuracy", strconv.FormatFloat(llc.HorizontalAccuracy, 'E', -1, 64))
	}
	if llc.Heading != 0 {
		v.Add("heading", strconv.FormatInt(llc.Heading, 10))
	}
	if llc.ProximityAlertRadius != 0 {
		v.Add("proximity_alert_radius", strconv.FormatInt(llc.ProximityAlertRadius, 10))
	}
	if llc.ReplyMarkup != nil {
		v.Add("reply_markup", FormatMarkup(llc.ReplyMarkup))
	}

	return v, nil
}

func (llc *EditMessageLLConf) method() string {
	return "editMessageLiveLocation"
}

func NewEditMessageLL(longitude float64, latit float64) *EditMessageLLConf {
	return &EditMessageLLConf{
		Longitude: longitude,
		Latitude:  latit,
	}
}

type StopMessageLiveLocation struct {
	ChatID          int64
	MessageID       int64
	InlineMessageID int64
	ReplyMarkup     objects.InlineKeyboardMarkup
}

// GetUpdate method fields
// https://core.telegram.org/bots/api#getting-updates
type GetUpdatesConfig struct {
	Offset         int64
	Limit          uint
	Timeout        uint
	AllowedUpdates []string
}

func (guc *GetUpdatesConfig) values() (url.Values, error) {
	v := url.Values{}
	if guc.Offset != 0 {
		v.Add("offset", strconv.FormatInt(guc.Offset, 10))
	}
	if guc.Limit != 0 {
		v.Add("limit", strconv.FormatUint(uint64(guc.Limit), 10))
	}
	if guc.Timeout != 0 {
		v.Add("timeout", strconv.FormatUint(uint64(guc.Timeout), 10))
	}
	if len(guc.AllowedUpdates) > 0 {
		bs, err := json.Marshal(guc.AllowedUpdates)
		if err != nil {
			return v, err
		}
		v.Add("allowed_updates", BytesToString(bs))
	}

	return v, nil
}

func (guc *GetUpdatesConfig) method() string {
	return "getUpdates"
}

// Uses for default values for Sending updates
func NewGetUpdateConfig(Offset int64) *GetUpdatesConfig {
	return &GetUpdatesConfig{
		Offset:  Offset,
		Limit:   20,
		Timeout: 5,
	}
}

type GetMyCommandsConfig struct {
	Scope        objects.BotCommandScope // optional
	LanguageCode string                  // optional
}

func (gmcc *GetMyCommandsConfig) values() (url.Values, error) {
	v := url.Values{}
	if gmcc.Scope != nil {
		v.Add("scope", ObjectToJson(gmcc.Scope))
	}
	if gmcc.LanguageCode != "" {
		v.Add("language_code", gmcc.LanguageCode)
	}
	return v, nil
}

func (gmcc *GetMyCommandsConfig) method() string {
	return "getMyCommands"
}

// DeleteMyCommandsConfig ...
type DeleteMyCommandsConfig struct {
	Scope        objects.BotCommandScope // optional
	LanguageCode string                  // optional
}

func (dmcc *DeleteMyCommandsConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("scope", ObjectToJson(dmcc.Scope))
	if dmcc.LanguageCode != "" {
		v.Add("language_code", dmcc.LanguageCode)
	}
	return v, nil
}

func (dmcc *DeleteMyCommandsConfig) method() string {
	return "deleteMyCommands"
}

func NewDeleteMyCommandsConf() *DeleteMyCommandsConfig {
	return &DeleteMyCommandsConfig{}
}

// SetMyCommandsConfig ...
type SetMyCommandsConfig struct {
	Commands     []*objects.BotCommand
	Scope        objects.BotCommandScope
	LanguageCode string
}

func (smcc *SetMyCommandsConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("commands", ObjectToJson(smcc.Commands))
	if smcc.LanguageCode != "" {
		v.Add("language_code", smcc.LanguageCode)
	}
	if smcc.Scope != nil {
		v.Add("scope", ObjectToJson(smcc.Scope))
	}
	return v, nil
}

func (smcc *SetMyCommandsConfig) method() string {
	return "setMyCommands"
}

func NewSetMyCommands(commands ...*objects.BotCommand) *SetMyCommandsConfig {
	return &SetMyCommandsConfig{
		Commands: commands,
	}
}

// DeleteWebhookConfig ...
type DeleteWebhookConfig struct {
	DropPendingUpdates bool
}

func (dwc *DeleteWebhookConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("drop_pending_updates", strconv.FormatBool(dwc.DropPendingUpdates))
	return v, nil
}

func (dwc *DeleteWebhookConfig) method() string {
	return "deleteWebhook"
}

func NewDeleteWebHook(drop_pending_updates bool) *DeleteWebhookConfig {
	return &DeleteWebhookConfig{
		DropPendingUpdates: drop_pending_updates,
	}
}

// SendDiceConfig https://core.telegram.org/bots/api#senddice
type SendDiceConfig struct {
	ChatID                   int64
	Emoji                    string
	DisableNotifications     bool
	ReplyToMessageId         int64
	AllowSendingWithoutReply bool
	ProtectContent           bool
	// ReplyMarkup will be type of objects.KeynoardMarkup not inline, and reply and etc.
	ReplyMarkup interface{}
}

func (sdc *SendDiceConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(sdc.ChatID, 10))
	if sdc.Emoji != "" {
		v.Add("emoji", sdc.Emoji)
	}
	v.Add("disable_notification", strconv.FormatBool(sdc.DisableNotifications))
	if sdc.ReplyToMessageId != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(sdc.ReplyToMessageId, 10))
	}
	v.Add("allow_sending_without_reply", strconv.FormatBool(sdc.AllowSendingWithoutReply))
	if sdc.ReplyMarkup != nil {
		v.Add("reply_markup", FormatMarkup(sdc.ReplyMarkup))
	}
	v.Add("protect_content", strconv.FormatBool(sdc.ProtectContent))
	return v, nil
}

func (sdc *SendDiceConfig) method() string {
	return "sendDice"
}

func NewSendDice(emoji string) *SendDiceConfig {
	return &SendDiceConfig{
		Emoji: emoji,
	}
}

// SendPollConfig Use this method to send a native poll
// https://core.telegram.org/bots/api#sendpoll
type SendPollConfig struct {
	ChatID   int64
	Question string   // VarChar(300) limit 300 chars
	Options  []string // starts with 2->10 limit, 1-100 char limit

	// Vezet, Vezet
	IsAnonymous bool
	Type        string

	AllowsMultipleAnswers bool
	CorrectOptionId       int64
	Explanation           string
	ExpalnationParseMode  string
	ExplnationEntites     []*objects.MessageEntity

	// Using int time, here can be used time.Time
	OpenPeriod     int64
	CloseDate      int64
	IsClosed       bool
	ProtectContent bool

	// Please, always turn off this
	DisableNotifications     bool
	ReplyToMessageID         int64
	AllowSendingWithoutReply bool
	// ReplyMarkup              *objects.KeyboardMarkup
}

func (spc *SendPollConfig) values() (url.Values, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(spc.ChatID, 10))
	v.Add("question", spc.Question)
	// lucky, lucky
	v.Add("is_anonymous", strconv.FormatBool(spc.IsAnonymous))
	if spc.Type != "" {
		v.Add("type", spc.Type)
	}
	v.Add("allows_multiple_answers", strconv.FormatBool(spc.AllowsMultipleAnswers))
	v.Add("correct_option_id", strconv.FormatInt(spc.CorrectOptionId, 10))
	if spc.Explanation != "" {
		v.Add("explanation", spc.Explanation)
	}
	if spc.ExpalnationParseMode != "" {
		v.Add("explanation_parse_mode", spc.ExpalnationParseMode)
	}
	if spc.ExplnationEntites != nil {
		v.Add("explanation_entities", ObjectToJson(spc.ExplnationEntites))
	}
	v.Add("open_period", strconv.FormatInt(spc.OpenPeriod, 10))
	v.Add("close_date", strconv.FormatInt(spc.CloseDate, 10))
	v.Add("is_closed", strconv.FormatBool(spc.IsClosed))
	v.Add("disable_notifications", strconv.FormatBool(spc.DisableNotifications))
	if spc.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(spc.ReplyToMessageID, 10))
	}
	v.Add("protect_content", strconv.FormatBool(spc.ProtectContent))
	return v, nil
}

func (spc *SendPollConfig) method() string {
	return "sendPoll"
}

func NewSendPoll(question string, options []string) *SendPollConfig {
	return &SendPollConfig{
		Question: question,
		Options:  options,
	}
}

// GetUserProfilePhotosConf represents getUserProfilePhotos method fields
// https://core.telegram.org/bots/api#getUserProfilePhotos
type GetUserProfilePhotosConf struct {
	UserId int64
	Offset int
	Limit  int
}

func (guppc *GetUserProfilePhotosConf) values() (url.Values, error) {
	v := url.Values{}

	v.Add("user_id", strconv.FormatInt(guppc.UserId, 10))
	v.Add("offset", strconv.Itoa(guppc.Offset))
	v.Add("limit", strconv.Itoa(guppc.Limit))

	return v, nil
}

func (guppc *GetUserProfilePhotosConf) method() string {
	return "getUserProfilePhotos"
}

type SendChatActionConf struct {
	ChatID int64
	Action string // see utils for actions type
}

func (scac *SendChatActionConf) values() (url.Values, error) {
	v := url.Values{}

	v.Add("chat_id", strconv.FormatInt(scac.ChatID, 10))
	v.Add("action", scac.Action)

	return v, nil
}

func (scac *SendChatActionConf) method() string {
	return "sendChatAction"
}

type SendContactConfig struct {
	ChatID                   interface{} // req
	PhoneNumber              string      // req
	FirstName                string      // req
	LastName                 string
	Vcard                    string
	DisableNotifiaction      bool
	ReplyToMessageID         int64
	AllowSendingWithoutReply bool
	ReplyKeyboard            interface{}
	ProtectContent           bool
}

func (scc *SendContactConfig) values() (url.Values, error) {
	v := url.Values{}
	switch t := scc.ChatID.(type) {
	case int64:
		v.Add("chat_id", strconv.FormatInt(t, 10))
	case string:
		v.Add("chat_id", t)
	}
	v.Add("phone_number", scc.PhoneNumber)
	v.Add("first_name", scc.FirstName)
	if scc.LastName != "" {
		v.Add("last_name", scc.LastName)
	}
	if scc.Vcard != "" {
		v.Add("vcard", scc.Vcard)
	}
	v.Add("disable_notification", strconv.FormatBool(!scc.DisableNotifiaction))
	if scc.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(scc.ReplyToMessageID, 10))
	}
	if scc.ReplyKeyboard != nil {
		v.Add("reply_keyboard", FormatMarkup(scc.ReplyKeyboard))
	}
	v.Add("protect_content", strconv.FormatBool(scc.ProtectContent))
	return v, nil
}

func (scc *SendContactConfig) method() string {
	return "sendContact"
}

// SendVenueConfig ...
type SendVenueConfig struct {
	ChatID                   interface{} // req
	Latitude                 float64     // req
	Longitude                float64     // req
	Title                    string      // req
	Address                  string      // req
	FoursQuareId             string
	FoursQuareType           string
	GooglePlaceId            string
	GooglePlaceType          string
	DisableNotification      bool
	ProtectContent           bool
	ReplyToMessageId         int64
	AllowSendingWithoutReply bool
	ReplyMarkup              interface{}
}

func (svc *SendVenueConfig) values() (url.Values, error) {
	v := url.Values{}
	switch t := svc.ChatID.(type) {
	case int64:
		v.Add("chat_id", strconv.FormatInt(t, 10))
	case string:
		v.Add("chat_id", t)
	}
	v.Add("latitude", strconv.FormatFloat(svc.Latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(svc.Longitude, 'f', -1, 64))
	v.Add("title", svc.Title)
	v.Add("address", svc.Address)
	v.Add("allow_sending_without_reply", strconv.FormatBool(svc.AllowSendingWithoutReply))
	v.Add("protect_content", strconv.FormatBool(svc.ProtectContent))
	if svc.ReplyToMessageId != 0 {
		v.Add("reply_to_message_id", strconv.FormatInt(svc.ReplyToMessageId, 10))
	}
	v.Add("disable_notification", strconv.FormatBool(!svc.DisableNotification))
	if svc.GooglePlaceId != "" {
		v.Add("google_place_id", svc.GooglePlaceId)
	}
	if svc.GooglePlaceType != "" {
		v.Add("google_place_type", svc.GooglePlaceType)
	}
	if svc.FoursQuareId != "" {
		v.Add("four_square_id", svc.FoursQuareId)
	}
	if svc.FoursQuareType != "" {
		v.Add("four_square_type", svc.FoursQuareType)
	}
	if svc.ReplyMarkup != nil {
		v.Add("reply_markup", FormatMarkup(svc.ReplyMarkup))
	}

	return v, nil
}

func (svc *SendVenueConfig) method() string {
	return "sendVenue"
}

// BanChatMemberConfig ...
type BanChatMemberConfig struct {
	ChatID         int64
	UserID         int64
	UntilDate      time.Duration
	RevokeMessages bool
}

func (bcm *BanChatMemberConfig) values() (url.Values, error) {
	v := url.Values{}

	v.Add("chat_id", strconv.FormatInt(bcm.ChatID, 10))
	v.Add("user_id", strconv.FormatInt(bcm.UserID, 10))

	if bcm.UntilDate != 0 {
		v.Add("until_date", strconv.FormatInt((int64)(bcm.UntilDate), 10))
	}
	v.Add("revoke_messages", strconv.FormatBool(bcm.RevokeMessages))

	return v, nil
}

func (bcm *BanChatMemberConfig) method() string {
	return "banChatMember"
}

func NewBanChatMember(chat_id int64, user_id int64) *BanChatMemberConfig {
	return &BanChatMemberConfig{
		ChatID: chat_id,
		UserID: user_id,
	}
}

type RestrictChatMemberConfig struct {
	ChatID      int64
	UserID      int64
	Permissions *objects.ChatMemberPermissions
	UntilDate   time.Duration
}

func (rc *RestrictChatMemberConfig) method() string {
	return "restrictChatMember"
}

func (rc *RestrictChatMemberConfig) values() (url.Values, error) {
	v := url.Values{}

	v.Add("chat_id", strconv.FormatInt(rc.ChatID, 10))
	v.Add("user_id", strconv.FormatInt(rc.UserID, 10))
	v.Add("permissions", ObjectToJson(rc.Permissions))
	if rc.UntilDate != 0 {
		v.Add("until_date", fmt.Sprintln(rc.UntilDate))
	}

	return v, nil
}

func NewRestrictMember(chat_id, user_id int64, perms *objects.ChatMemberPermissions) *RestrictChatMemberConfig {
	return &RestrictChatMemberConfig{
		ChatID:      chat_id,
		UserID:      user_id,
		Permissions: perms,
	}
}

type EditChatInviteLinkConf struct {
	ChatID             int64
	InviteLink         string
	Name               string
	ExpireDate         int64
	MemberLimit        int
	CreatesJoinRequest bool
}

func (eilc *EditChatInviteLinkConf) values() (v url.Values, _ error) {
	v.Add("chat_id", strconv.FormatInt(eilc.ChatID, 10))
	v.Add("invite_link", eilc.InviteLink)
	if eilc.Name != "" {
		v.Add("name", eilc.Name)
	}
	if eilc.ExpireDate != 0 {
		v.Add("expire_date", strconv.FormatInt(eilc.ExpireDate, 10))
	}
	if eilc.MemberLimit != 0 {
		v.Add("member_limit", strconv.Itoa(eilc.MemberLimit))
	}
	v.Add("creates_join_request", strconv.FormatBool(eilc.CreatesJoinRequest))
	return
}

func (eilc *EditChatInviteLinkConf) method() string {
	return "editChatInviteLink"
}

func NeweditChatInviteLink(ChatID int64, invite_link string) *EditChatInviteLinkConf {
	return &EditChatInviteLinkConf{
		ChatID:     ChatID,
		InviteLink: invite_link,
	}
}

type PromoteChatMemberConfig struct {
	UserPermissionsConfig

	IsAnonymous bool  `json:"is_anonymous"`
	ChatID      int64 `json:"chat_id"` // required
	UserID      int64 `json:"user_id"` // required
}

func (pcmc PromoteChatMemberConfig) values() (url.Values, error) {
	v, _ := pcmc.UserPermissionsConfig.values()

	v.Add("chat_id", strconv.FormatInt(pcmc.ChatID, 10))
	v.Add("user_id", strconv.FormatInt(pcmc.ChatID, 10))
	v.Add("is_anonymous", strconv.FormatBool(pcmc.IsAnonymous))

	return v, nil
}

func (pcmc PromoteChatMemberConfig) method() string {
	return "promoteChatMember"
}

func NewPromoteChat(ChatId int64, UserID int64) *PromoteChatMemberConfig {
	return &PromoteChatMemberConfig{
		ChatID: ChatId,
		UserID: UserID,
	}
}

// TODO: SendInvoice

type SendGameConfig struct {
	BaseChat                        // chat_id is required
	GameShortName            string // required
	DisableNotfication       bool
	ProtectContent           bool
	AllowSendingWithoutReply bool
}

func (sg *SendGameConfig) values() (url.Values, error) {
	v, _ := sg.BaseChat.values()
	v.Add("game_short_name", sg.GameShortName)
	v.Add("disable_notification", strconv.FormatBool(sg.DisableNotfication))
	v.Add("protect_content", strconv.FormatBool(sg.ProtectContent))
	v.Add("allow_sending_without_reply", strconv.FormatBool(sg.AllowSendingWithoutReply))
	return v, nil
}

func (sg *SendGameConfig) method() string {
	return "sendGame"
}

type SendStickerConfig struct {
	BaseChat                                    // chat_id required
	Sticker                  *objects.InputFile // required
	ProtectContent           bool
	AllowSendingWithoutReply bool
}

func (stc *SendStickerConfig) values() (url.Values, error) {
	v, _ := stc.BaseChat.values()
	v.Add("protect_content", strconv.FormatBool(stc.ProtectContent))
	v.Add("allow_sending_without_reply", strconv.FormatBool(stc.AllowSendingWithoutReply))
	return v, nil
}

func (stc *SendStickerConfig) method() string {
	return "sendSticker"
}

func (stc *SendStickerConfig) params() (map[string]string, error) {
	v, _ := stc.BaseChat.params()
	uv, _ := stc.values()
	urlValuesToMapString(uv, v)
	return v, nil
}

func (stc *SendStickerConfig) getFiles() []*objects.InputFile {
	return []*objects.InputFile{stc.Sticker}
}

func NewSendStickerConf(base BaseChat, Sticker *objects.InputFile) *SendStickerConfig {
	return &SendStickerConfig{
		BaseChat: base,
		Sticker:  Sticker,
	}
}

type CreateNewStickerSetConf struct {
	UserID        int64  // required
	Name          string // required
	Title         string // required
	PngSticker    *objects.InputFile
	TgsSticker    *objects.InputFile
	WebmSticker   *objects.InputFile
	Emojis        string // required
	ContainsMasks bool
	MaskPosition  *objects.MaskPosition
}

func (cns *CreateNewStickerSetConf) getFiles() []*objects.InputFile {
	return []*objects.InputFile{cns.PngSticker, cns.TgsSticker, cns.WebmSticker}
}

func (cns *CreateNewStickerSetConf) values() (url.Values, error) {
	v := url.Values{}

	v.Add("user_id", strconv.FormatInt(cns.UserID, 10))
	v.Add("name", cns.Name)
	v.Add("title", cns.Title)
	v.Add("emojis", cns.Emojis)
	v.Add("conatins_masks", strconv.FormatBool(cns.ContainsMasks))

	if cns.MaskPosition != nil {
		bs, err := json.Marshal(cns.MaskPosition)
		if err != nil {
			return nil, err
		}
		v.Add("mask_poistion", BytesToString(bs))
	}

	return v, nil
}

func (cns *CreateNewStickerSetConf) params() (map[string]string, error) {
	w := make(map[string]string)
	v, err := cns.values()
	if err != nil {
		return nil, err
	}
	urlValuesToMapString(v, w)
	return w, nil
}

func (cns *CreateNewStickerSetConf) method() string {
	return "createNewStickerSet"
}

func NewCreateStickerSet(user_id int64, name, title, emojis string) *CreateNewStickerSetConf {
	return &CreateNewStickerSetConf{
		UserID: user_id,
		Name:   name,
		Title:  title,
		Emojis: emojis,
	}
}

type AddStickerToSetConf struct {
	UserID       int64  // required
	Name         string // required
	PngSticker   *objects.InputFile
	TgsSticker   *objects.InputFile
	WebmSticker  *objects.InputFile
	Emojis       string // required
	MaskPosition *objects.MaskPosition
}

func (ast *AddStickerToSetConf) values() (url.Values, error) {
	v := url.Values{}

	v.Add("user_id", strconv.FormatInt(ast.UserID, 10))
	v.Add("name", ast.Name)
	v.Add("emojis", ast.Emojis)

	if ast.MaskPosition != nil {
		bs, err := json.Marshal(ast.MaskPosition)
		if err != nil {
			return nil, err
		}
		v.Add("mask_poistion", BytesToString(bs))
	}

	return v, nil
}

func (ast *AddStickerToSetConf) params() (map[string]string, error) {
	w := make(map[string]string)
	v, err := ast.values()
	if err != nil {
		return nil, err
	}
	urlValuesToMapString(v, w)
	return w, nil
}

func (ast *AddStickerToSetConf) getFiles() []*objects.InputFile {
	return []*objects.InputFile{ast.PngSticker, ast.TgsSticker, ast.WebmSticker}
}

func (ast *AddStickerToSetConf) method() string {
	return "addStickerToSet"
}

func NewAddStickerToSet(user_id int64, name, emojis string) *AddStickerToSetConf {
	return &AddStickerToSetConf{
		UserID: user_id,
		Name:   name,
		Emojis: emojis,
	}
}

type SetStickerSetThumbConf struct {
	Name   string
	UserId int64
	Thumb  *objects.InputFile
}

func (sst *SetStickerSetThumbConf) values() (url.Values, error) {
	v := url.Values{}

	v.Add("name", sst.Name)
	v.Add("user_id", strconv.FormatInt(sst.UserId, 10))

	return v, nil
}

func (sst *SetStickerSetThumbConf) method() string {
	return "setStickerSetThumb"
}

func (sst *SetStickerSetThumbConf) params() (map[string]string, error) {
	p := make(map[string]string)
	p["name"] = sst.Name
	p["user_id"] = strconv.FormatInt(sst.UserId, 10)
	return p, nil
}

func (sst *SetStickerSetThumbConf) getFiles() []*objects.InputFile {
	return []*objects.InputFile{sst.Thumb}
}

func NewSetStickerSetThumb(name string, userId int64) *SetStickerSetThumbConf {
	return &SetStickerSetThumbConf{
		Name:   name,
		UserId: userId,
	}
}

type AnswerWebAppQueryConf struct {
	WebAppQueryId string
	Result        objects.InlineQueryResult
}

func (awpq *AnswerWebAppQueryConf) method() string {
	return "answerWebAppQuery"
}

func (awpq *AnswerWebAppQueryConf) values() (url.Values, error) {
	v := url.Values{}

	v.Add("web_app_query_id", awpq.WebAppQueryId)

	bs, err := json.Marshal(awpq.Result)
	if err != nil {
		return nil, err
	}
	v.Add("result", BytesToString(bs))

	return v, nil
}

func NewAnswerWebQuery(web_app_query_id string, result objects.InlineQueryResult) *AnswerWebAppQueryConf {
	return &AnswerWebAppQueryConf{
		WebAppQueryId: web_app_query_id,
		Result:        result,
	}
}
