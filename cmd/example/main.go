package main

import (
	"fmt"
	"os"

	"github.com/tyler-johnson/minimist"
)

var usage = `
example - an example os minimist

Usage:
  -h, -?, --help     This screen
  -v, --version      Print version
`

func main() {
	argm := minimist.Parse(&minimist.Config{
		Alias: map[string]string{
			"h": "help", "?": "help",
			"v": "version",
		},
	})
	fmt.Printf("%q\n", os.Args)

	// cmd --help || cmd --h || cmd -?
	if v, ok := argm.Get("help").(bool); ok && v {
		fmt.Println(usage)
	}

	// cmd -v || cmd --version
	if v, ok := argm.Get("version").(bool); ok && v {
		fmt.Println("1.0")
	}
}
