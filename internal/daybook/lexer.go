package daybook

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type RawEntry struct {
	header  []byte
	content []byte
}

func (re *RawEntry) String() string {
	return fmt.Sprintf("header:\n%s\n\ncontent:\n%s\n\n", string(re.header), string(re.content))
}

type lexerState int

const (
	preamble lexerState = 1 << iota
	sectionHead
	sectionBody
)

func (s lexerState) String() string {
	switch s {
	case preamble:
		return "PREAMBLE"
	case sectionHead:
		return "SECTIONHEAD"
	case sectionBody:
		return "SECTIONBODY"
	default:
		return "UNKNOWN"
	}
}

type lexer struct {
	source    *bufio.Reader
	state     lexerState
	preamble  []byte
	currBlock []byte
	count     int
}

func (l *lexer) String() string {
	return fmt.Sprintf(`
	state:    %s
	count:    %d
	preamble: %s
	block:    %s
	`, l.state, l.count, string(l.preamble), string(l.currBlock))
}

func NewLexer(r io.Reader) *lexer {
	return &lexer{
		source: bufio.NewReader(r),
		state:  preamble,
	}
}

func (l *lexer) newBlock() {
	l.currBlock = nil
}

func (l *lexer) readPreamble() error {
	if l.state != preamble {
		return fmt.Errorf("unexpected call to readPreamble(), state: %s", l.state)
	}

	l.newBlock()
	for {
		data, err := l.source.ReadBytes('\n')
		if err != nil {
			return err
		}

		if bytes.Equal([]byte("---\n"), data) {
			l.state = sectionHead
			return nil
		}

		l.preamble = append(l.preamble, data...)
	}
}

func (l *lexer) readSectionHead() error {
	if l.state != sectionHead {
		return fmt.Errorf("unexpected call to readSectionHead(), state: %s", l.state)
	}

	l.newBlock()
	for {
		data, err := l.source.ReadBytes('\n')
		if err != nil {
			return err
		}

		if bytes.Equal([]byte("\n"), data) {
			l.state = sectionBody
			return nil
		}

		l.currBlock = append(l.currBlock, data...)
	}
}

func (l *lexer) readSectionBody() error {
	if l.state != sectionBody {
		return fmt.Errorf("unexpected call to readSectionBody(), state: %s", l.state)
	}

	l.newBlock()
	for {
		data, err := l.source.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if bytes.Equal([]byte("---\n"), data) {
			l.state = sectionHead
			return nil
		}

		l.currBlock = append(l.currBlock, data...)
	}
}

func (l *lexer) NextEntry() (*RawEntry, error) {
	if err := l.readPreamble(); err != nil {
		return nil, err
	}

	if err := l.readSectionHead(); err != nil {
		return nil, err
	}
	hd := l.currBlock

	if err := l.readSectionBody(); err != nil {
		return nil, err
	}

	l.count++

	return &RawEntry{
		header:  hd,
		content: l.currBlock,
	}, nil
}
