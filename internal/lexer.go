package internal

import (
	"bufio"
	"bytes"
	"io"
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
	eof             = rune(0)
	CURR_LINE_NO    = 1
	CURR_COL_NO     = 0
	PREV_COL_NO     = -1
	LAST_CH_NEWLINE = false
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isNewline(ch rune) bool {
	return ch == '\n'
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
	CURR_COL_NO++
	if ch == '\n' {
		CURR_LINE_NO++
		PREV_COL_NO = CURR_COL_NO - 1
		CURR_COL_NO = 0
		LAST_CH_NEWLINE = true
	} else {
		LAST_CH_NEWLINE = false
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
	CURR_COL_NO--
	if LAST_CH_NEWLINE {
		CURR_LINE_NO--
		CURR_COL_NO = PREV_COL_NO
		LAST_CH_NEWLINE = false
	}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a digit then consume as an ident.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isDigit(ch) { // an ident shouldn't start with a decimal point e.g. timestamp = .356
		s.unread()
		return s.scanIdent()
	} else {
		// there's no other possible correct first character
		return ILLEGAL, string(ch)
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

	firstCh := s.read()
	if !(isDigit(firstCh) || isDecimalPoint(firstCh)) {
		return ILLEGAL, ""
	}

	buf.WriteRune(firstCh)

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !(isDigit(ch) || isDecimalPoint(ch)) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return IDENT, buf.String()
}
