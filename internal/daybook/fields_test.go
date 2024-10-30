package daybook

import (
	"slices"
	"testing"
)

func TestHeaderFields(t *testing.T) {
	tests := []struct {
		name   string
		header string
		fields []string
	}{
		{"empty", "", nil},
		{"one field", "topic: git", []string{"topic"}},
		{"two fields", "topic: git\ndate: 2024-10-28", []string{"topic", "date"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := RawEntry{header: []byte(tt.header)}
			fields, err := fieldNames(e)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if slices.Compare(tt.fields, fields) != 0 {
				t.Errorf("unexpected fields, want: %v, got: %v", tt.fields, fields)
			}
		})
	}

}
