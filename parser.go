package main

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"time"
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

func (p *Parser) Parse() ([]Beat, error) {
	beats := []Beat{}
	currentBeat := &Beat{}

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
			// we know from harmonixset that the float timestamp is in s
			tstampFloatSeconds, err := strconv.ParseFloat(lit, 64)
			if err != nil {
				log.Fatalf("parse error %v on line number TODO", err)
			}
			us := fmt.Sprintf("%dus", int(tstampFloatSeconds*1000000.0))
			if currentBeat.Timestamp, parseErr = time.ParseDuration(us); parseErr != nil {
				log.Fatalf("parse error %v on line number TODO", parseErr)
			}
		case 1:
			if currentBeat.Beat, parseErr = strconv.Atoi(lit); parseErr != nil {
				log.Fatalf("parse error %v on line number TODO", parseErr)
			}
		case 2:
			if currentBeat.Bar, parseErr = strconv.Atoi(lit); parseErr != nil {
				log.Fatalf("parse error %v (line/char no TODO)", parseErr)
			}
		default:
			log.Fatal("too many tokens")
		}

		pos++

		if pos == 3 {
			beats = append(beats, *currentBeat)
			currentBeat = &Beat{}
			pos = 0
		}

		// we're done here
		if tok, _ := p.scanIgnoreWhitespace(); tok != IDENT {
			break
		}
		// otherwise rewind
		p.unscan()
	}

	return beats, nil
}
