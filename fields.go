package daybook

import (
	"bufio"
	"fmt"
	"strings"
)

type field struct {
	name  string
	count int
}

func extractFields(e RawEntry) ([]string, error) {
	scan := bufio.NewScanner(strings.NewReader(string(e.header)))
	scan.Split(bufio.ScanLines)

	fields := make([]string, 0)
	for scan.Scan() {
		line := scan.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("unrecognized header line %q", line)
		}

		fields = append(fields, strings.TrimSpace(parts[0]))
	}

	if err := scan.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}
