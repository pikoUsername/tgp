package dispatcher

import (
	"github.com/pikoUsername/tgp/objects"
)

// HandlerFunc ...
type HandlerFunc func(upd *objects.Update)

// Another level of abstraction
type HandlerType struct {
	Callback *HandlerFunc
	Filters  []Filter
}

// CheckForFilters iterate all filters and call Check method for check
func (ht *HandlerType) CheckForFilters(u *objects.Update) bool {
	for _, f := range ht.Filters {
		b := f.Check(u)
		if !b {
			return false
		}
	}
	return true
}

// Call uses for checking using filters
func (ht *HandlerType) Call(u *objects.Update) {
	fr := ht.CheckForFilters(u)
	if !fr {
		return
	}

	if u != nil {
		cb := *ht.Callback
		cb(u)
	}
}

// Interface for creating custom HandlerObj
type HandlerObj interface {
	Register(handler HandlerFunc, filters ...Filter)
	Unregister(handler *HandlerFunc)
	Notify(update *objects.Update)
	RegisterMiddleware(middlewares ...MiddlewareFunc)
}

// HandlerObj uses for save Callback
type DefaultHandlerObj struct {
	handlers   []HandlerType
	Middleware MiddlewareManager
}

// NEwDHandlerObj creates new DefaultHandlerObj
func NewDHandlerObj(dp *Dispatcher) *DefaultHandlerObj {
	return &DefaultHandlerObj{
		Middleware: NewMiddlewareManager(dp),
	}
}

// Register, append to Callbacks, e.g handler functions
func (ho *DefaultHandlerObj) Register(f HandlerFunc, filters ...Filter) {
	ht := HandlerType{
		Callback: &f,
		Filters:  filters,
	}

	ho.handlers = append(ho.handlers, ht)
}

// Unregister checkout to memory address
// and cut up it if find something, with same address
func (ho *DefaultHandlerObj) Unregister(handler *HandlerFunc) {
	var index int
	for i, h := range ho.handlers {
		if h.Callback == handler {
			// deleting from slice
			index = i - 1
			if index < 0 {
				index = 0
			}
			ho.handlers = append(ho.handlers[:index], ho.handlers[i:]...)
		}
	}
}

// RegisterMiddleware looks like a bad code
// for example, you want to register every user which writed to you bot
// You can registerMiddleware for MessageHandler, not for all handlers
// Or maybe want to make throttling middleware, just Registers middleware
//
// Example of middlware see in handler_test.go
func (ho *DefaultHandlerObj) RegisterMiddleware(f ...MiddlewareFunc) {
	ho.Middleware.Register(f...)
}

// Notify is from aiogram framework
// Notify is notify all callbacks in handler
// when middlewares activates, middleware calls a handler
// Just triggers one, you must call Handler in Middleware,
func (ho *DefaultHandlerObj) Notify(upd *objects.Update) {
	for _, cb := range ho.handlers {
		be := ho.Middleware.Trigger(upd)
		if !be {
			continue
		}
		cb.Call(upd)
	}
}
