package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

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
			f, err := os.Create(fpath)
			if err != nil {
				return err
			}
			f.Close()
			f.WriteString("# This is the ccd-generated config file of the TSV.")
			f.WriteString("# Be sure to use the tab to separate when editing.")
			Logf("info", "The config file has been successfully initialized.")
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
