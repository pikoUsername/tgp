package dispatcher_test

import (
	"testing"

	"github.com/pikoUsername/tgp/bot"
	"github.com/pikoUsername/tgp/dispatcher"
)

var (
	TestToken = "1780004238:AAGFsgj2pxzXWoUqn25YohCEb1ENKIQOr1Q"
)

func GetDispatcher(t *testing.T) (error, *dispatcher.Dispatcher) {
	b, err := bot.NewBot(TestToken, true, "HTML")
	if err != nil {
		t.Error(err)
	}
	return nil, &dispatcher.Dispatcher{Bot: b}
}

func TestNewDispatcher(t *testing.T) {
	err, dp := GetDispatcher(t)
	if err != nil {
		t.Error(err)
	}
	if dp == nil {
		t.Error("Oh no, Dispatcher didnt created, fix it")
	}
}