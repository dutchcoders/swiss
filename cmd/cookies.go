package cmd

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

var (
	_ = Register("cookies", CookiePlugin{})
)

type CookiePlugin struct {
}

func (p CookiePlugin) Do(env *Environment) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Path", "Name", "Value"})

	for _, cookie := range env.Client.Jar.Cookies(env.URL) {
		table.Append([]string{cookie.Path, cookie.Name, cookie.Value})
	}

	table.Render()

	return nil
}
