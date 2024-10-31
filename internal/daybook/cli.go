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
	var cli cliConfig
	if err := cli.fromArgs(args); err != nil {
		return err
	}

	rawEntries, err := cli.loadRawData()
	if err != nil {
		return err
	}

	switch {
	case cli.flagFields:
		return cli.processFields(rawEntries)
	case cli.flagSummary:
		return cli.processSummary(rawEntries)
	case cli.flagOut:
		return cli.processOut(rawEntries)
	}

	return nil
}

type cliConfig struct {
	flagFields  bool
	flagCounts  bool
	flagSummary bool
	flagOut     bool
	args        []string
}

func (cli *cliConfig) fromArgs(args []string) error {
	fset := flag.NewFlagSet(args[0], flag.ExitOnError)
	fset.BoolVar(&cli.flagFields, "fields", false, "list fields")
	fset.BoolVar(&cli.flagCounts, "c", false, "include counts")
	fset.BoolVar(&cli.flagSummary, "summary", false, "show summary")
	fset.BoolVar(&cli.flagOut, "out", false, "print in canonical form")

	if err := fset.Parse(args[1:]); err != nil {
		return err
	}

	cli.args = fset.Args()

	return nil
}

func (cli *cliConfig) loadRawData() ([]RawEntry, error) {
	scan := NewScanner(os.Stdin)
	entries, err := scan.Entries()
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (cli *cliConfig) processFields(entries []RawEntry) error {
	fields, err := headerFields(entries)
	if err != nil {
		return err
	}
	if cli.flagCounts {
		for _, field := range fields {
			fmt.Printf("%3d %s\n", field.count, field.name)
		}
	} else {
		for _, field := range fields {
			fmt.Printf("%s\n", field.name)
		}
	}

	return nil
}

func (cli *cliConfig) processSummary(entries []RawEntry) error {
	return nil
}

func (cli *cliConfig) processOut(entries []RawEntry) error {
	for i, raw := range entries {
		e, err := transform(raw)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", e)

		if i > 0 {
			fmt.Printf("---\n")
		}
	}
	return nil
}
