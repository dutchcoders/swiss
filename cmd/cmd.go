package cmd

import (
	"fmt"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	_ "github.com/fatih/color"
	"github.com/minio/cli"
	"github.com/op/go-logging"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/user"
	"path"
	"strings"
)

// censys
// sqlmap
// nmap
// phantom
// identify
// plugins
// dns
// brute / automatic form filler
// tor
// trackers (analytics)
// lua ?
// pageshot?

var Version = "0.1"
var helpTemplate = `NAME:
{{.Name}} - {{.Usage}}

DESCRIPTION:
{{.Description}}

USAGE:
{{.Name}} {{if .Flags}}[flags] {{end}}command{{if .Flags}}{{end}} [arguments...]

COMMANDS:
{{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
{{end}}{{if .Flags}}
FLAGS:
{{range .Flags}}{{.}}
{{end}}{{end}}
VERSION:
` + Version +
	`{{ "\n"}}`

var log = logging.MustGetLogger("tool")

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "url",
		Usage: "url",
	},
	cli.StringFlag{
		Name:  "proxy",
		Usage: "proxy",
		Value: "socks5://127.0.0.1:9150",
	},
}

type Cmd struct {
	*cli.App
}

type Client struct {
	*http.Client
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "")

	return c.Client.Do(req)
}

type Environment struct {
	URL    *url.URL
	Client Client // *http.Client

	Method string
}

var completer = readline.NewPrefixCompleter(
	readline.PcItemDynamic(func(line string) []string {
		names := make([]string, len(plugins))
		for k, _ := range plugins {
			names = append(names, k)
		}
		return names
	}),
	readline.PcItem("quit"),
	readline.PcItem("exit"),
)

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

var format = logging.MustStringFormatter(
	"%{color} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

func New() *Cmd {
	app := cli.NewApp()
	app.Name = "swiss"
	app.Author = "DutchSec"
	app.Usage = "swiss"
	app.Description = `Swiss knife for security analysts`
	app.Flags = globalFlags
	app.CustomAppHelpTemplate = helpTemplate
	app.Commands = []cli.Command{}

	app.Before = func(c *cli.Context) error {
		logBackends := []logging.Backend{}

		backend1 := logging.NewLogBackend(os.Stdout, "", 0)

		backend1Formatter := logging.NewBackendFormatter(backend1, format)

		backend1Leveled := logging.AddModuleLevel(backend1Formatter)

		level, err := logging.LogLevel("debug")
		if err != nil {
			panic(err)
		}

		backend1Leveled.SetLevel(level, "")

		logBackends = append(logBackends, backend1Leveled)

		/*
		   for _, log := range c.Logging {

		           var output io.Writer = os.Stdout
		           switch log.Output {
		           case "stdout":
		           case "stderr":
		                   output = os.Stderr
		           default:
		                   output, err = os.OpenFile(log.Output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		           }

		           if err != nil {
		                   panic(err)
		           }

		           backend1 := logging.NewLogBackend(output, "", 0)

		           backend1Formatter := logging.NewBackendFormatter(backend1, format)

		           backend1Leveled := logging.AddModuleLevel(backend1Formatter)

		           level, err := logging.LogLevel(log.Level)
		           if err != nil {
		                   panic(err)
		           }

		           backend1Leveled.SetLevel(level, "")

		           logBackends = append(logBackends, backend1Leveled)
		   }
		*/

		logging.SetBackend(logBackends...)

		return nil
	}

	app.Action = func(c *cli.Context) {
		fmt.Println(color.YellowString("Swiss, the swiss army knife for security specialists."))

		jar, _ := cookiejar.New(&cookiejar.Options{})

		d := net.Dial

		if c.GlobalString("proxy") == "" {
		} else if u, err := url.Parse(c.GlobalString("proxy")); err != nil {
			panic(err)
		} else if v, err := proxy.FromURL(u, proxy.Direct); err != nil {
			panic(err)
		} else {
			log.Infof("Using proxy %s.", c.GlobalString("proxy"))

			d = v.Dial
		}

		env := &Environment{
			Client: Client{
				&http.Client{
					Transport: &http.Transport{
						Dial: d,
					},
					Jar: jar,
				},
			},
		}

		if u, err := url.Parse(c.GlobalString("url")); err != nil {
			panic(fmt.Sprintf("Could not parse url: %s", err.Error()))
		} else {
			env.URL = u
		}

		// will create .swiss home folder for history

		home := ""
		if usr, err := user.Current(); err != nil {
		} else {
			home = usr.HomeDir
		}

		p := path.Join(home, ".swiss")
		if _, err := os.Stat(p); err == nil {
		} else if !os.IsNotExist(err) {
			log.Errorf("Could not create .swiss folder: %s", err.Error())
		} else if err = os.Mkdir(p, 0700); err != nil {
			log.Errorf("Could not create .swiss folder: %s", err.Error())
		}

		l, err := readline.NewEx(&readline.Config{
			Prompt:          "\033[31m»\033[0m ",
			HistoryFile:     path.Join(p, "history"),
			AutoComplete:    completer,
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",

			HistorySearchFold: true,
		})
		if err != nil {
			panic(err)
		}
		defer l.Close()

		for {
			line, err := l.Readline()
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}

			line = strings.TrimSpace(line)
			switch {
			case strings.HasPrefix(strings.ToUpper(line), "SET "):
				parts := strings.Split(line[4:], "=")

				if len(parts) >= 2 {
					if strings.ToUpper(parts[0]) == "URL" {
						urlstr := strings.Join(parts[1:], "=")

						if rel, err := url.Parse(urlstr); err != nil {
							fmt.Println(color.RedString(fmt.Sprintf("Could not parse url: %s: %s", urlstr, err.Error())))
						} else if u := env.URL.ResolveReference(rel); err != nil {
							fmt.Println(color.RedString(fmt.Sprintf("Could not resolve reference: %s: %s", rel, err.Error())))
						} else if u.Scheme == "" {
							fmt.Println(color.RedString(fmt.Sprintf("URL contains no scheme: %s", u.String())))
						} else if u.Host == "" {
							fmt.Println(color.RedString(fmt.Sprintf("URL contains no host: %s", u.String())))
						} else {
							env.URL = u

							fmt.Printf("URL set to %s\n", env.URL.String())
						}
					} else if strings.ToUpper(parts[0]) == "METHOD" {
						env.Method = strings.ToUpper(parts[1])

						fmt.Printf("Method set to %s\n", env.Method)
					}
				}
			case line == "quit":
				fallthrough
			case line == "exit":
				goto exit
			case line == "":
			default:
				if p, ok := plugins[line]; ok {
					if err := p.Do(env); err != nil {
						fmt.Printf(color.RedString(fmt.Sprintf("Error executing plugin: %s\n", line)))
					}
				} else {
					fmt.Printf(color.RedString(fmt.Sprintf("Command not found: %s\n", line)))
				}
			}
		}
	exit:
	}

	return &Cmd{
		App: app,
	}
}
