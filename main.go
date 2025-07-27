package main

import (
    "fmt"
    "os"
    "github.com/mistium/raingoer/runner"
    "github.com/mistium/raingoer/tokens"
)

func parseMain(block string) []string {
	var lines = tokens.Tokenise(block, '\n')
	return lines
}

func main() {

	fi, err := os.ReadFile(os.Args[1])
    if err != nil {
        panic(err)
    }

    var code = string(fi)
	var lines = tokens.Tokenise(code, '\n')

	for i, line := range lines {
		var split_line = tokens.Tokenise(line, ' ')
		if len(split_line) == 0 || strings.TrimSpace(split_line[0]) == "" { continue }
		runner.Run(split_line[0], split_line[1:]...)
	}
}
