package cmd

import (
	"github.com/puerkitobio/goquery"

	"fmt"
)

var (
	_ = Register("body", BodyPlugin{})
)

type BodyPlugin struct {
}

func (p BodyPlugin) Do(env *Environment) error {
	resp, err := env.Client.Get(env.URL.String())
	if err != nil {
		fmt.Println("Coud not read input: %s", err.Error())
		return err
	}

	defer resp.Body.Close()

	fmt.Printf("Statuscode: %d\n", resp.StatusCode)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Coud not read input: %s", err.Error())
		return err
	}

	fmt.Println(doc.Text())

	return nil
}
