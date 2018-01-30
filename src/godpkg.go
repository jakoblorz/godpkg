package main

import (
	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCli("godpkg", "1.0.0")
	c.Args := os.Args[1:]


	install := func() (cli.Command, error) {

	}

	c.Commands = map[string]cli.CommandFactory{
		"install": install
	}
}
