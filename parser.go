package main

import (
	"io"
	"log"
	"fmt"
	"time"
	"strconv"
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

func (p *Parser) Parse() (*Beat, error) {
	beat := &Beat{}

	var parseErr error
	pos := 0
	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		switch pos {
		case 0:
			if beat.Timestamp, parseErr = time.ParseDuration(lit); parseErr != nil {
				log.Fatal("parse error %v on line number TODO", parseErr)
			}
		case 1:
			if beat.Beat, parseErr = strconv.Atoi(lit); parseErr != nil {
				log.Fatal("parse error %v on line number TODO", parseErr)
			}
		case 2:
			if beat.Bar, parseErr = strconv.Atoi(lit); parseErr != nil {
				log.Fatal("parse error %v (line/char no TODO)", parseErr)
			}
		default:
			log.Fatal("too many tokens (line/char no TODO), %s is extra", lit)
		}

		// If the next token is not an ident then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != IDENT {
			p.unscan()
			break
		}
		pos++
	}

	return beat, nil
}
