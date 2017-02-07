package cmd

import (
	"github.com/puerkitobio/goquery"

	"fmt"
	"net/url"
	"sort"
)

var (
	_ = Register("links", LinksPlugin{})
)

type LinksPlugin struct {
}

func (p LinksPlugin) Do(env *Environment) error {
	fmt.Println("Links", env.URL.String())

	resp, err := env.Client.Get(env.URL.String())
	if err != nil {
		fmt.Println("Could not read input: %s", err.Error())
		return err
	}

	defer resp.Body.Close()

	fmt.Printf("Statuscode: %d\n", resp.StatusCode)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Coud not read input: %s", err.Error())
		return err
	}

	fmt.Println("Loaded page")

	links := map[string]map[string]interface{}{}

	extractLinks := func(elem string, attr string) {
		doc.Find(elem).Each(func(i int, s *goquery.Selection) {
			if val, ok := s.Attr(attr); !ok {
				log.Warningf("Could not find attribute %s within element %s.", attr, elem)
			} else if rel, err := url.Parse(val); err != nil {
			} else if u := env.URL.ResolveReference(rel); err != nil {
			} else {
				if _, ok := links[u.Host]; !ok {
					links[u.Host] = map[string]interface{}{}
				}

				links[u.Host][u.String()] = u.String()
			}
		})
	}

	extractLinks("script", "src")
	extractLinks("img", "src")
	extractLinks("form", "action")
	extractLinks("link", "href")
	extractLinks("base", "href")
	extractLinks("a", "href")

	for k, v := range links {
		fmt.Printf("[+] %s\n", k)
		fmt.Println("===============================")

		sortedLinks := make([]string, 0, len(v))
		for k, _ := range v {
			sortedLinks = append(sortedLinks, k)
		}

		sort.Strings(sortedLinks)
		for _, v := range sortedLinks {
			fmt.Println(v)
		}

		fmt.Println("")
	}

	return nil
}
