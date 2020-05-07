package main

import (
	"fmt"

	"github.com/motemen/go-colorine"
)

var logger = &colorine.Logger{
	Prefixes: colorine.Prefixes{
		"debug": colorine.Verbose,
		"info":  colorine.Info,
		"warn":  colorine.Warn,
		"error": colorine.Error,
		"":      colorine.Verbose,
	},
}

// Logf is a logger that displays logs colorfully.
func Logf(prefix, pattern string, args ...interface{}) {
	logger.Log(prefix, fmt.Sprintf(pattern, args...))
}
