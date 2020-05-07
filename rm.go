package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

var commandRm = &cli.Command{
	Name:    "remove",
	Aliases: []string{"rm"},
	Usage:   "remove the registered point",
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
		var deleteAlias string
		view.SetSelectedFunc(func(idx int, mainText string, secondaryText string, shortcut rune) {
			text = mainText
			deleteAlias = secondaryText
			app.Stop()
		})
		if err := app.SetRoot(view, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}

		// If a user tries to register an alias named "Quit", a bug will occur.
		if text == "Quit" {
			return nil
		}

		// Make new file for deleting configuration.
		f, err = os.Create(fpath)
		if err != nil {
			return err
		}
		defer f.Close()
		for _, p := range paths {
			if p.alias == deleteAlias {
				continue
			}
			if _, err = f.WriteString(fmt.Sprintf("%s\t%s\n", p.alias, p.path)); err != nil {
				return err
			}
		}
		Logf("info", "alias (%s) has been successfully deleted.", deleteAlias)

		return nil
	},
}
