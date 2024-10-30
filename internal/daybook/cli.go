package daybook

import (
	"flag"
	"fmt"
	"os"
)

func Main() int {
	if err := CLI(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	return 0
}

func CLI(args []string) error {
	var cfg appConfig
	if err := cfg.fromArgs(args); err != nil {
		return err
	}

	scan := NewScanner(os.Stdin)
	entries, err := scan.Entries()
	if err != nil {
		return err
	}

	if cfg.showFields {
		fields, err := headerFields(entries)
		if err != nil {
			return err
		}
		if cfg.showCounts {
			for _, field := range fields {
				fmt.Printf("%3d %s\n", field.count, field.name)
			}
		} else {
			for _, field := range fields {
				fmt.Printf("%s\n", field.name)
			}
		}
	}

	return nil
}

type appConfig struct {
	showFields bool
	showCounts bool
	summary    bool
}

func (c *appConfig) fromArgs(args []string) error {
	fset := flag.NewFlagSet(args[0], flag.ExitOnError)
	fset.BoolVar(&c.showFields, "fields", false, "list fields")
	fset.BoolVar(&c.showCounts, "c", false, "include counts")
	fset.BoolVar(&c.summary, "summary", false, "print summary")

	return fset.Parse(args[1:])
}
