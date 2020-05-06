package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/d-tsuji/ecd"
	"github.com/urfave/cli/v2"
)

var errCommandHelp = fmt.Errorf("command help shown")

func main() {
	app := cli.NewApp()
	app.Commands = []*cli.Command{
		commandAdd,
		commandLs,
		commandRm,
	}
	app.Version = fmt.Sprintf("%s", ecd.Version)
	err := app.Run(os.Args)
	if err != nil {
		if err != errCommandHelp {
		}
	}
}

var commandAdd = &cli.Command{
	Name:  "add",
	Usage: "add directory to your list",
	Action: func(c *cli.Context) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		filename := filepath.Join(home, ".config", "ecd", "config")
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		alias := c.Args().Get(0)
		dir := c.Args().Get(1)
		if dir == "" {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			dir = wd
		}
		if _, err = f.WriteString(fmt.Sprintf("%s\t%s\n", alias,dir)); err != nil {
			return err
		}
		return nil
	},
}

var commandLs = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "show list for registered points",
	Action: func(c *cli.Context) error {
		return nil
	},
}

var commandRm = &cli.Command{
	Name:    "remove",
	Aliases: []string{"rm"},
	Usage:   "remove the registered point",
	Action: func(c *cli.Context) error {
		return nil
	},
}
