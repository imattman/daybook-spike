package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/imattman/daybook-spike"
)

func main() {
	os.Exit(Main())
}

func Main() int {
	config, err := configure()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return 1
	}

	if err := Run(config); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return 1
	}
	return 0
}

type config struct {
	extractFields bool
	showCounts    bool
}

func configure() (*config, error) {
	extractFields := flag.Bool("fields", false, "extract header fields")
	showCounts := flag.Bool("counts", false, "show counts")
	flag.Parse()

	return &config{
		extractFields: *extractFields,
		showCounts:    *showCounts,
	}, nil
}

func Run(cfg *config) error {
	var in io.Reader
	if len(flag.Args()) > 0 {
		fname := flag.Arg(0)
		fin, err := os.Open(fname)
		if err != nil {
			return err
		}
		defer fin.Close()

		in = fin
	}

	if cfg.extractFields {
		scan := daybook.NewScanner(in)
		ents, err := scan.Entries()
		if err != nil {
			return err
		}
		_ = ents
	}

	return nil
}
