package filters

import (
	"regexp"

	"github.com/pikoUsername/tgp/objects"
)

type Regexp struct {
	RegexpPattern *regexp.Regexp
}

func (r *Regexp) Check(u *objects.Update) bool {
	var content string
	if u.Message != nil {
		content = u.Message.Text
	} else if u.CallbackQuery != nil {
		content = u.CallbackQuery.Text
	} else if u.Poll != nil {
		content = u.Poll.Question
	} else {
		return false
	}

	match := string(r.RegexpPattern.Find([]byte(content)))
	return match != ""
}

func NewRegexp(re string) (*Regexp, error) {
	rex, err := regexp.Compile(re)
	if err != nil {
		return &Regexp{}, err
	}

	return &Regexp{
		RegexpPattern: rex,
	}, nil
}