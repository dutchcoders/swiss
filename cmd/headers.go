package cmd

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	_ = Register("headers", HeadersPlugin{})
)

type HeadersPlugin struct {
}

func (p HeadersPlugin) Do(env *Environment) error {
	req, err := http.NewRequest(env.Method, env.URL.String(), nil)
	if err != nil {
		fmt.Println("Coud not read input: %s", err.Error())
		return err
	}

	resp, err := env.Client.Do(req)
	if err != nil {
		fmt.Println("Coud not read input: %s", err.Error())
		return err
	}

	defer resp.Body.Close()

	fmt.Println("[+] Statuscode")
	fmt.Println("===========")
	fmt.Printf("%s\n", resp.Status)
	fmt.Println("")

	fmt.Println("[+] Headers")
	fmt.Println("===========")
	for k, v := range resp.Header {
		fmt.Printf("%s: %s\n", k, strings.Join(v, `\n`))
	}
	fmt.Println("")

	return nil
}
