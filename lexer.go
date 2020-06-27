package main

import (
	"bufio"
	"io"
	"bytes"
	"fmt"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	WS

	IDENT
)

func (t *Token) String() string {
	switch *t {
	case ILLEGAL:
		return "illegal"
	case EOF:
		return "eof"
	case WS:
		return "whitespace"
	case IDENT:
		return "identifier"
	default:
		return "fatal: undefined"
	}
}

var (
	eof = rune(0)
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isNewline(ch rune) bool {
	return ch == '\n' || ch == '\r'
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isDecimalPoint(ch rune) bool {
	return ch == '.'
}

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }


// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a digit then consume as an ident.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isDigit(ch) || isDecimalPoint(ch) {
		s.unread()
		return s.scanIdent()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		ch := s.read()
		if ch == eof {
			break
		}

		if isDigit(ch) || isDecimalPoint(ch) {
			fmt.Printf("DEBUG: writing %c for IDENT\n", ch)
			_, _ = buf.WriteRune(ch)
		} else {
			s.unread()
			break
		}
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}
