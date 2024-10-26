package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/imattman/daybook-spike/internal/daybook"
)

func main() {

	test := `
# example data in preamble

# more preamble
---
date: 2024-10-26
topic: Go
location: library
instructor: Professor Squirrel

This is the start of the data block

---
date: 2024-10-26
topic: Linux

Some cool linux CLI tricks
	`

	lex := daybook.NewLexer(strings.NewReader(test))
	ent, err := lex.NextEntry()
	if err != nil {
		fail("lex error: %v", err)
	}

	fmt.Printf("Lexer info: %v\n", lex)
	fmt.Printf("entry: %s\n", ent)
}

func fail(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}
