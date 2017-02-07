package main

import "github.com/dutchcoders/swiss/cmd"

func main() {
	app := cmd.New()
	app.RunAndExitOnError()
}
