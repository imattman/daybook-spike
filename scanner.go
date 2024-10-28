package daybook

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

type scannerState int

const (
	preamble scannerState = 1 << iota
	entryHead
	entryBody
)

var yamlFence = regexp.MustCompile(`^\s*---\s*\n$`)
var emptyLine = regexp.MustCompile(`^\s*\n$`)

type Scanner struct {
	r     *bufio.Reader
	state scannerState
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r:     bufio.NewReader(r),
		state: preamble,
	}
}

func (s *Scanner) scanUntil(pat *regexp.Regexp, nextState scannerState) ([]byte, bool, error) {
	chunk := make([]byte, 0)

	for {
		line, err := s.r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return chunk, true, nil // special EOF case
			}
			return nil, true, err // regular error
		}
		if pat.Find(line) != nil {
			s.state = nextState
			break
		}
		chunk = append(chunk, line...)
	}

	return chunk, false, nil
}

func (s *Scanner) scanPreamble() ([]byte, bool, error) {
	if s.state != preamble {
		return nil, false, nil
	}

	return s.scanUntil(yamlFence, entryHead)
}

func (s *Scanner) scanEntryHead() ([]byte, bool, error) {
	if s.state != entryHead {
		return nil, false, nil
	}

	return s.scanUntil(emptyLine, entryBody)
}

func (s *Scanner) scanEntryBody() ([]byte, bool, error) {
	if s.state != entryBody {
		return nil, false, nil
	}

	return s.scanUntil(yamlFence, entryHead)
}

func (s *Scanner) Entries() ([]RawEntry, error) {
	preamble, done, err := s.scanPreamble()
	if err != nil {
		return nil, err
	}
	if done {
		return nil, nil
	}
	fmt.Printf("preamble:\n%s\n", preamble)

	entries := make([]RawEntry, 0)
	for {
		head, _, err := s.scanEntryHead()
		if err != nil {
			return nil, err
		}
		body, done, err := s.scanEntryBody()
		if err != nil {
			return nil, err
		}

		entries = append(entries, RawEntry{
			header: head,
			body:   body,
		})

		if done {
			break
		}
	}

	return entries, nil
}
