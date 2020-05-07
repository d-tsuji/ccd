package main

import (
	"errors"
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
			ecd.Logf("error", "%+v", err)
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
		fpath := filepath.Join(home, ".config", "ecd", "config.tsv")
		os.MkdirAll(filepath.Dir(fpath), 0755)
		f, err := os.Create(fpath)
		if err != nil {
			return err
		}
		defer f.Close()

		alias := c.Args().Get(0)
		if alias == "" {
			return errors.New("alias is required")
		}
		dir := c.Args().Get(1)
		if dir == "" {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			dir = wd
		}
		if _, err = f.WriteString(fmt.Sprintf("%s\t%s\n", alias, dir)); err != nil {
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
