package daybook

import (
	"bufio"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var listItemPat = regexp.MustCompile(`^\s*- `)

type headerField struct {
	name  string
	count int
}

func (f headerField) String() string {
	return fmt.Sprintf("%s: %d", f.name, f.count)
}

func fieldNames(e RawEntry) ([]string, error) {

	scan := bufio.NewScanner(strings.NewReader(string(e.header)))
	scan.Split(bufio.ScanLines)

	fields := make([]string, 0)
	for scan.Scan() {
		line := scan.Text()
		if listItemPat.Find([]byte(line)) != nil {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) > 2 {
			return nil, fmt.Errorf("unrecognized header line %q", line)
		}

		fields = append(fields, strings.TrimSpace(parts[0]))
	}

	if err := scan.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}

func headerFields(entries []RawEntry) ([]headerField, error) {
	fmap := make(map[string]int)
	for _, e := range entries {
		fnames, err := fieldNames(e)
		if err != nil {
			return nil, err
		}

		for _, fname := range fnames {
			fmap[fname]++
		}
	}

	fcounts := make([]headerField, 0, len(fmap))
	for f, cnt := range fmap {
		fcounts = append(fcounts, headerField{f, cnt})
	}

	sort.Slice(fcounts, func(i, j int) bool {
		return fcounts[i].name < fcounts[j].name
	})

	return fcounts, nil
}
