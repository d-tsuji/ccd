package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mattn/go-tty"
	"github.com/urfave/cli/v2"
)

var commandEdit = &cli.Command{
	Name:  "edit",
	Usage: "edit the registered point",
	Action: func(c *cli.Context) error {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		tty, err := tty.Open()
		if err != nil {
			return err
		}
		defer tty.Close()
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
		cmd := exec.Command(editor, f.Name())
		cmd.Stdin = tty.Input()
		cmd.Stdout = tty.Output()
		cmd.Stderr = tty.Output()
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("abort renames: %s", err)
		}
		Logf("info", "Edit has been successfully completed.")
		return nil
	},
}
