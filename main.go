package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	errCommandHelp = fmt.Errorf("command help shown")

	characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

type config struct {
	alias string
	path  string
}

func main() {
	app := cli.NewApp()
	app.Usage = "Custom change directory tool"

	app.Commands = []*cli.Command{
		commandAdd,
		commandLs,
		commandEdit,
		commandRm,
	}
	app.Version = fmt.Sprintf("%s", version)
	err := app.Run(os.Args)
	if err != nil {
		if err != errCommandHelp {
			Logf("error", "%+v", err)
		}
	}
}
