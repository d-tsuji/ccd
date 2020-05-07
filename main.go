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
	app.Version = fmt.Sprintf("%s", Version)
	err := app.Run(os.Args)
	if err != nil {
		if err != errCommandHelp {
			Logf("error", "%+v", err)
		}
	}
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

var commandAdd = &cli.Command{
	Name:  "add",
	Usage: "add directory to your list",
	Action: func(c *cli.Context) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		fpath := filepath.Join(home, ".config", "ccd", "config.tsv")
		os.MkdirAll(filepath.Dir(fpath), 0755)
		if !exists(fpath) {
			Logf("info", "The config file has been successfully initialized.")
			f, err := os.Create(fpath)
			if err != nil {
				return err
			}
			f.Close()
		}

		f, err := os.Open(fpath)
		if err != nil {
			return err
		}
		defer f.Close()

		sc := bufio.NewScanner(f)

		var paths []config
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

		targetAlias := c.Args().Get(0)
		if targetAlias == "" {
			// TODO: Print an easy-to-understand help message.
			return fmt.Errorf("alias is required")
		}
		targetPath := c.Args().Get(1)
		if targetPath == "" {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			targetPath = wd
		}
		overwrite := false
		for i, p := range paths {
			if p.alias == targetAlias {
				sc := bufio.NewScanner(os.Stdin)
				fmt.Fprintln(os.Stdout, fmt.Sprintf(`Alias (%s: %s) has already been registered. Do you want me to override it?`, p.alias, p.path))
				if sc.Scan() {
					t := sc.Text()
					if t == "Y" || t == "y" || t == "YES" || t == "yes" || t == "Yes" {
						Logf("info", "Path overwrite: %s -> %s.", paths[i].path, targetPath)
						paths[i].path = targetPath
						overwrite = true
					} else {
						Logf("warn", "Alias (%s) was not registered.", targetAlias)
						return nil
					}
				}
			}
		}

		if !overwrite {
			paths = append(paths, config{targetAlias, targetPath})
		}

		// Make new file for updating configuration.
		f, err = os.Create(fpath)
		if err != nil {
			return err
		}
		defer f.Close()
		for _, p := range paths {
			if _, err = f.WriteString(fmt.Sprintf("%s\t%s\n", p.alias, p.path)); err != nil {
				return err
			}
		}
		Logf("info", "alias (%s: %s) has been successfully registered.", targetAlias, targetPath)

		return nil
	},
}

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

// TODO: Need implementation. Remove the config line and reconfigure the file.
var commandRm = &cli.Command{
	Name:    "remove",
	Aliases: []string{"rm"},
	Usage:   "remove the registered point",
	Action: func(c *cli.Context) error {
		return nil
	},
}

// TODO: Need implementation.
var commandEdit = &cli.Command{
	Name:  "edit",
	Usage: "edit the registered point",
	Action: func(c *cli.Context) error {
		return nil
	},
}
