package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

type InstallCommand struct {
}

func (*InstallCommand) Help() string {
	return "install dependency"
}

func (*InstallCommand) Run(args []string) int {
	fmt.Printf("install, %v\n", args)
	return 0
}

func (command *InstallCommand) Synopsis() string {
	return command.Help()
}

func main() {
	c := cli.NewCLI("godpkg", "1.0.0")
	c.Args = os.Args[1:]

	install := func() (cli.Command, error) {
		return &InstallCommand{}, nil
	}

	c.Commands = map[string]cli.CommandFactory{
		"install": install,
	}

	status, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(status)
}
