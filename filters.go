package tgp

import (
	"github.com/pikoUsername/tgp/objects"
)

// Filter uses for filtering comming text, data and etc.
// When Filter check is not success handler will not pass
// If not just continue comming update, and etc. etc.
// if you need aleardy ready filters check filters/
//
// and this filter have second type is: func(u *objects.Update) bool
// more simpler choice, but you cannot set variables to it
type Filter interface {
	Check(update *objects.Update) bool
}

// check out for filters
func checkFilters(filters []Filter, upd *objects.Update) bool {
	for _, filter := range filters {
		if !filter.Check(upd) {
			return false
		}
	}
	return true
}
