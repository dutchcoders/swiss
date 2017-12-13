package cmd

import (
	"github.com/puerkitobio/goquery"

	"fmt"
)

var (
	_ = Register("form", FormPlugin{})
)

type FormPlugin struct {
}

func (p FormPlugin) Do(env *Environment) error {
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

	// no form found

	doc.Find("form").Each(func(i int, v *goquery.Selection) {
		fmt.Println("[+] Form")
		fmt.Println("========")

		name := ""
		if val, ok := v.Attr("name"); !ok {
		} else {
			name = val
		}
		fmt.Printf("Name: %s\n", name)

		action := ""
		if val, ok := v.Attr("action"); !ok {
		} else {
			action = val
		}

		fmt.Printf("Action: %s\n", action)

		v.Find("select").Each(func(i int, s *goquery.Selection) {
			name := ""
			if val, ok := s.Attr("name"); !ok {
			} else {
				name = val
			}

			type_ := ""
			if val, ok := s.Attr("type"); !ok {
			} else {
				type_ = val
			}

			fmt.Printf("Found input type %s with name %s", type_, name)
		})

		v.Find("input").Each(func(i int, s *goquery.Selection) {
			name := ""
			if val, ok := s.Attr("name"); !ok {
			} else {
				name = val
			}

			type_ := ""
			if val, ok := s.Attr("type"); !ok {
			} else {
				type_ = val
			}

			fmt.Printf("Found input type %s with name %s\n", type_, name)
		})

		// method
		// action

		// find all input s
	})

	return nil
}
