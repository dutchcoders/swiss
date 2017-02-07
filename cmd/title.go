package cmd

import (
	"github.com/puerkitobio/goquery"

	"fmt"
	"github.com/kr/pretty"
)

var (
	_ = Register("title", TitlePlugin{})
)

type TitlePlugin struct {
}

func (p TitlePlugin) Do(env *Environment) error {
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

	title := ""
	if v := doc.Find("title"); v != nil {
		title = v.Text()
	}

	fmt.Println("Title:", title)

	description := ""
	if s := doc.Find("meta"); s != nil {
		if val, ok := s.Attr("description"); ok {
			description = val
		}
	}

	fmt.Println("Description:", description)

	pretty.Print(resp.Header)
	return nil
}
