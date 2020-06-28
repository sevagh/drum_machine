package internal

import (
	"strings"
	"testing"
)

func TestTokenToString(t *testing.T) {
	t1 := Token(0)
	t2 := Token(1)
	t3 := Token(2)
	t4 := Token(3)
	t5 := Token(4)

	if t1 != ILLEGAL || t1.String() != "illegal" {
		t.Fatalf("bad")
	}
	if t2 != EOF || t2.String() != "eof" {
		t.Fatalf("bad")
	}
	if t3 != WS || t3.String() != "whitespace" {
		t.Fatalf("bad")
	}
	if t4 != IDENT || t4.String() != "identifier" {
		t.Fatalf("bad")
	}
	if t5.String() != "fatal: undefined" {
		t.Fatalf("bad")
	}
}

func TestCharDelimiterComparions(t *testing.T) {
	raw := []rune{' ', '\t', '\n', eof}
	for i, c := range raw {
		if i < 3 {
			if !isWhitespace(c) {
				t.Fatalf("bad")
			}
		}
		if i == 4 {
			if isWhitespace(c) || isDigit(c) || isDecimalPoint(c) {
				t.Fatalf("bad")
			}
		}
	}
}

func TestCharIdentComparions(t *testing.T) {
	if !isDecimalPoint('.') {
		t.Fatalf("bad")
	}
	if isDigit('a') {
		t.Fatalf("bad")
	}
	if !isDigit('5') {
		t.Fatalf("bad")
	}
	if isDigit('\t') {
		t.Fatalf("bad")
	}
}

func TestScanReadsUntilEof(t *testing.T) {
	testString := `abcdef
123456`

	rdr := strings.NewReader(testString)

	s := NewScanner(rdr)

	ctr := 0
	var next rune

	for {
		next = s.read()
		if next == eof {
			if ctr != 13 {
				t.Fatalf("should've read 13 runes")
			}
			return
		}
		if ctr < 6 && next != rune(97+ctr) {
			t.Fatalf("bad")
		}
		if ctr == 6 && !isWhitespace(next) {
			t.Logf("next is: %c\n", next)
			t.Logf("next should've been: %c\n", rune(97+ctr))
			t.Fatalf("bad")
		}
		if ctr > 6 && next != rune(48-6+ctr) {
			t.Errorf("bad")
		}
		ctr++
	}
}

func TestScanUnreadPutsBackChar(t *testing.T) {
	testString := `abcdef
123456`

	rdr := strings.NewReader(testString)
	s := NewScanner(rdr)

	next := s.read()
	if next != 'a' {
		t.Fatal("bad")
	}
	s.unread()
	next = s.read()
	if next != 'a' {
		t.Fatal("bad")
	}
}

func TestScanWhitespaceAndIdent(t *testing.T) {
	testString := `abcdef    1345   
	foo`

	rdr := strings.NewReader(testString)
	s := NewScanner(rdr)

	var tok Token
	var lit string

	tok, lit = s.scanWhitespace()
	if tok != WS {
		t.Fatal("bad")
	}
	if lit != "a" {
		t.Fatal("bad")
	}

	tok, lit = s.scanIdent()
	if tok != ILLEGAL {
		t.Fatal("bad")
	}

	// skip past the letters 'abcdef'
	for i := 0; i < 4; i++ {
		t.Logf("consuming: %c", s.read())
	}

	tok, lit = s.scanWhitespace()
	tok, lit = s.scanIdent()
	if tok != IDENT {
		t.Fatal("bad")
	}

	if lit != "1345" {
		t.Logf("lit: %s\n", lit)
		t.Fatal("bad")
	}
}

func TestScanAll(t *testing.T) {
	testString := `1345  135.8
	133.7`

	rdr := strings.NewReader(testString)
	s := NewScanner(rdr)

	correct := 0

	var tok Token
	var lit string
	for {
		tok, lit = s.Scan()
		t.Logf("tok: %d, %s, lit: %s\n", tok, tok.String(), lit)
		if tok == ILLEGAL {
			break
		} else {
			switch lit {
			case "1345", "135.8", "133.7":
				correct++
			}
		}
	}

	if correct != 3 {
		t.Fatalf("bad")
	}
}
