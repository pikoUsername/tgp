package tgp

import "fmt"

// Interface for outer usage
type ITelegramServer interface {
	ApiURL(string, string) string
	FileURL(string, string) string
}

// TelegramApiServer(just copy paste from aiogram)
// make easier use custom telegram api server
type TelegramAPIServer struct {
	// Base telegram, sendMessage and etc.
	Base string `json:"base"`

	// Url for file transfer, CDN and etc.
	File string `json:"file"`
}

// NewTelegramApiServer ...
func NewTelegramApiServer(Base string) *TelegramAPIServer {
	template := "/bot%s/%s"
	// /bot%s/%s is /bot<TOKEN>/<METHOD>
	return &TelegramAPIServer{
		Base: fmt.Sprint(Base, template),
		File: fmt.Sprint(Base, "/file", template),
	}
}

// ApiUrl creates from base telegram url
func (tas *TelegramAPIServer) ApiURL(token string, method string) string {
	return fmt.Sprintf(tas.Base, token, method)
}

// FileUrl Creates at base of tas.File string
// a url for send a request
func (tas *TelegramAPIServer) FileURL(Token string, File string) string {
	return fmt.Sprintf(tas.File, Token, File)
}

// Default telegram api server url
var (
	DefaultTelegramServer = NewTelegramApiServer("https://api.telegram.org")
)
