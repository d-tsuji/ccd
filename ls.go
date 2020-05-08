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
	"github.com/urfave/cli/v2"
)

var commandLs = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "show list for registered points",
	Action: func(c *cli.Context) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		fpath := filepath.Join(home, ".config", "ccd", "config.tsv")
		f, err := os.Open(fpath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		sc := bufio.NewScanner(f)

		paths := []config{}
		for sc.Scan() {
			text := sc.Text()
			if text == "" {
				continue
			}
			s := strings.Split(text, "\t")
			if len(s) != 2 {
				fmt.Errorf("invalid config, please check the config.tsv")
			}
			alias, path := s[0], s[1]
			fmt.Println(alias, path)
			paths = append(paths, config{alias, path})
		}

		app := tview.NewApplication()
		view := tview.NewList()
		for i, p := range paths {
			view.AddItem(fmt.Sprintf("%s -> %s\n", p.alias, p.path), p.alias, characters[i], nil)
		}
		view.AddItem("Quit", "Press to exit", 'q', func() { app.Stop() })
		view.ShowSecondaryText(false)
		var text string
		var targetAlias string
		var index int
		view.SetSelectedFunc(func(idx int, mainText string, secondaryText string, shortcut rune) {
			index = idx
			text = mainText
			targetAlias = secondaryText
			app.Stop()
		})
		if err := app.SetRoot(view, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}

		// If a user tries to register an alias named "Quit", a bug will occur.
		if text == "Quit" {
			return nil
		}

		if err := clipboard.WriteAll(paths[index].path); err != nil {
			return err
		}
		Logf("info", `"%s" is copied to clipboard!`, targetAlias)
		return nil
	},
}
