package daybook

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type isoDate time.Time

func (d isoDate) String() string {
	return time.Time(d).Format(time.DateOnly)
}

type RawEntry struct {
	header []byte
	body   []byte
}

func (re RawEntry) String() string {
	return fmt.Sprintf("%s\n\n%s", string(re.header), string(re.body))
}

type Entry struct {
	Date       isoDate
	Topic      string
	Location   string
	References []string
	Body       string
}

func (e Entry) String() string {
	return fmt.Sprintf("date: %s\ntopic: %s\nlocation: %s\nreferences: %v\n\n%s",
		e.Date, e.Topic, e.Location, strings.Join(e.References, ", "), e.Body,
	)
}

func transform(raw RawEntry) (*Entry, error) {
	head := struct {
		Date       time.Time
		Topic      string
		Location   string
		References []string
	}{}

	if err := yaml.Unmarshal(raw.header, &head); err != nil {
		return nil, err
	}

	return &Entry{
		Date:       isoDate(head.Date),
		Topic:      head.Topic,
		Location:   head.Location,
		References: head.References,
		Body:       string(raw.body),
	}, nil

}
