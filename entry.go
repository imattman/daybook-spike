package daybook

import "fmt"

type RawEntry struct {
	header []byte
	body   []byte
}

func (re RawEntry) String() string {
	return fmt.Sprintf("%s\n\n%s", string(re.header), string(re.body))
}
