# ccd [![Go Report Card](https://goreportcard.com/badge/github.com/d-tsuji/ccd)](https://goreportcard.com/report/github.com/d-tsuji/ccd)

Custom Change Directory Manager

You can add aliases to directories and operate them intuitively with TUI. But the actual `cd` command **cannot** be executed. This is because the environment variables of the parent process cannot be modified from the child process.

## Usage

The following four operations are supported.

- add

Register by giving an alias to the directory. The path can be optional. The default is to register the current directory.

```
$ ccd add <alias> <path>
```

- ls

Get a list of registered aliases and copy the selected alias to the clipboard.

```
$ ccd ls
```

- rm

Removes the selected alias.

```
$ ccd rm
```

- edit

Edit the registered aliases and directories with an editor. By default, vi is used.

```
$ ccd edit
```

### Config

The registered aliases and directory settings can be found in `~/.config/ccd/config.tsv`.

## Author

Tsuji Daishiro

## LICENSE

This software is licensed under the MIT license.
