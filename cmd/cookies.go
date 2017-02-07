package cmd

import (
	"github.com/kr/pretty"
)

var (
	_ = Register("cookies", CookiePlugin{})
)

type CookiePlugin struct {
}

func (p CookiePlugin) Do(env *Environment) error {
	pretty.Print(env.Client.Jar)

	return nil
}
