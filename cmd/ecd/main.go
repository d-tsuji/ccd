package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/rivo/tview"

	"github.com/d-tsuji/ecd"
	"github.com/urfave/cli/v2"
)

var errCommandHelp = fmt.Errorf("command help shown")

func main() {
	app := cli.NewApp()

	// Default action when the following aliases are not matched.
	app.Action = commandCd
	app.UsageText = fmt.Sprintf("%s <alias or path> | %s command [command options]",
		filepath.Base(os.Args[0]),
		filepath.Base(os.Args[0]),
	)
	app.Usage = "Extend change directory tool"

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
			// TODO: Print an easy-to-understand help message.
			return fmt.Errorf("alias is required")
		}
		dirpath := c.Args().Get(1)
		if dirpath == "" {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			dirpath = wd
		}
		if _, err = f.WriteString(fmt.Sprintf("%s\t%s\n", alias, dirpath)); err != nil {
			return err
		}
		return nil
	},
}

type Item struct {
	alias string
	path  string
}

var characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

var commandLs = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "show list for registered points",
	Action: func(c *cli.Context) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		fpath := filepath.Join(home, ".config", "ecd", "config.tsv")
		f, err := os.Open(fpath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		sc := bufio.NewScanner(f)
		var pathList []*Item
		for sc.Scan() {
			text := sc.Text()
			ecd.Logf("debug", "text: %s", text)
			s := strings.Split(text, "\t")
			if len(s) != 2 {
				fmt.Errorf("invalid config")
			}
			alias, path := s[0], s[1]
			pathList = append(pathList, &Item{alias, path})
		}

		app := tview.NewApplication()
		view := tview.NewList()
		for i, v := range pathList {
			// TODO: Convert an index to a-zA-Z rune.
			view.AddItem(fmt.Sprintf("%s -> %s\n", v.alias, v.path), v.alias, characters[i], nil)
		}
		view.AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
		view.ShowSecondaryText(false)
		var targetAlias string
		view.SetSelectedFunc(func(idx int, mainText string, secondaryText string, shortcut rune) {
			targetAlias = secondaryText
			app.Stop()
		})
		if err := app.SetRoot(view, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}

		if err := clipboard.WriteAll(targetAlias); err != nil {
			return err
		}
		ecd.Logf("info", `"%s" is copied to clipboard!`, targetAlias)
		return nil
	},
}

// TODO: Need implementation.
var commandRm = &cli.Command{
	Name:    "remove",
	Aliases: []string{"rm"},
	Usage:   "remove the registered point",
	Action: func(c *cli.Context) error {
		return nil
	},
}

var commandCd = func(c *cli.Context) error {
	targetAlias := c.Args().Get(0)
	var targetPath string

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	fpath := filepath.Join(home, ".config", "ecd", "config.tsv")
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		text := sc.Text()
		s := strings.Split(text, "\t")
		if len(s) != 2 {
			fmt.Errorf("invalid config")
		}
		alias, path := s[0], s[1]
		if alias == targetAlias {
			targetPath = path
		}
	}

	if targetPath == "" {
		fmt.Errorf(`alias "%s" could not be found in config (%s)`, targetAlias, fpath)
	}

	// TODO: Change the current directory of the terminal itself. Probably difficult.
	if err := os.Chdir(targetPath); err != nil {
		return err
	}
	return nil
}
