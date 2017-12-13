package cmd

import (
	"bufio"
	"os"

	"github.com/puerkitobio/goquery"

	"fmt"
)

var (
	_ = Register("title", TitlePlugin{})
)

type TitlePlugin struct {
}

func (p TitlePlugin) Do(env *Environment) error {
	r, err := os.OpenFile("/tmp/urls_2.txt", os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		func() error {
			u := scanner.Text() // Println will add back the final '\n'
			u = "http://" + u + "/"
			fmt.Println(u)

			resp, err := env.Client.Get(u)
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
			return nil
			// pretty.Print(resp.Header)

		}()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return nil
}
