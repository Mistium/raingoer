package main

import (
	"os"
	"fmt"
	"strings"
	"github.com/mistium/raingoer/tokens"
	"github.com/mistium/raingoer/runner"
)

func parseMain(block string) []string {
	var lines = tokens.Tokenise(block, '\n')
}

func main() {

	fi, err := os.ReadFile(os.Args[1])
    if err != nil {
        panic(err)
    }

    var code = fmt.Sprintf("%s\n", fi)
	var lines = tokenise(code, '\n')
	
	for i, line := range lines {
		var split_line = tokens.Tokenise(line, ' ')
		if len(split_line) == 0 { continue }
		runner.Run(split_line[0], split_line[1:]...)
	}
}
