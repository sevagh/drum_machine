package main

import (
	"fmt"
	"io"
	"log"
	"math"
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
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string, done bool) {
	tok, lit = p.scan()
	if tok == WS {
		if lit[len(lit)-1] == '\n' {
			done = true
		}
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
		tok, lit, done := p.scanIgnoreWhitespace()
		if done {
			break
		}

		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected field on line %d, col %d", lit, CURR_LINE_NO, CURR_COL_NO)
		}

		switch pos {
		case 0:
			// we know from harmonixset that the float timestamp is in s
			if currentBeat.Timestamp, parseErr = strconv.ParseFloat(lit, 64); parseErr != nil {
				return nil, fmt.Errorf("parse error %v on line %d, col %d", parseErr, CURR_LINE_NO, CURR_COL_NO)
			}
		case 1:
			if currentBeat.Beat, parseErr = strconv.Atoi(lit); parseErr != nil {
				return nil, fmt.Errorf("parse error %v on line %d, col %d", parseErr, CURR_LINE_NO, CURR_COL_NO)
			}
		case 2:
			if currentBeat.Bar, parseErr = strconv.Atoi(lit); parseErr != nil {
				return nil, fmt.Errorf("parse error %v on line %d, col %d", parseErr, CURR_LINE_NO, CURR_COL_NO)
			}
		default:
			return nil, fmt.Errorf("too many tokens on line %d, col %d", CURR_LINE_NO, CURR_COL_NO)
		}

		pos++

		if pos == 3 {
			beats = append(beats, *currentBeat)
			currentBeat = &Beat{}
			pos = 0
		}

		// we're done here
		if tok, _, done := p.scanIgnoreWhitespace(); tok != IDENT {
			if done {
				break
			}
			if tok == ILLEGAL {
				return nil, fmt.Errorf("too many tokens on line %d, col %d", CURR_LINE_NO, CURR_COL_NO)
			}
			break
		}
		// otherwise rewind
		p.unscan()
	}

	return beats, nil
}

func validateBeats(beats []Beat) (float64, time.Duration, Song, error) {
	song := Song{}

	barTempos := make(map[int][]float64)

	// acumulate tempo via timestamp increments
	currentBar := beats[0].Bar

	for i := 1; i < len(beats); i++ {
		if _, ok := barTempos[currentBar]; !ok {
			barTempos[currentBar] = []float64{}
		}
		barTempos[currentBar] = append(barTempos[currentBar], beats[i].Timestamp-beats[i-1].Timestamp)

		currentBeat := beats[i].Beat

		// check if cycled to a new bar
		// or it's the last beat
		if currentBeat < beats[i-1].Beat || i == len(beats)-1 {
			// check if bars are monotonically increasing
			if i != len(beats)-1 {
				if beats[i].Bar != currentBar+1 {
					return -1, -1, nil, fmt.Errorf("bars should be monotonically increasing")
				}
			} else {
				if beats[i].Bar != currentBar {
					return -1, -1, nil, fmt.Errorf("final note is in different bar")
				}
			}

			// assume 4/4 for now
			bar := Bar{NBeats: beats[i-1].Beat}

			thisBarsTempos := barTempos[currentBar]

			bar.Tempo = thisBarsTempos[0]
			for i := 1; i < len(thisBarsTempos); i++ {
				if math.Abs(thisBarsTempos[i]-thisBarsTempos[i-1]) > TempoAllowedVarianceSecs {
					log.Printf("[WARN] timestamp increments should be similar in a bar, %f and %f are suspect", thisBarsTempos[i], thisBarsTempos[i-1])
				}
			}

			for _, barTempoJump := range barTempos[currentBar] {
				bar.Tempo += barTempoJump
			}
			// average the tempo
			bar.Tempo /= float64(bar.NBeats)

			// transform it to bpm
			bar.Tempo = 60 / bar.Tempo

			song = append(song, bar)

			currentBar++
		} else {
			if currentBeat != beats[i-1].Beat+1 {
				return -1, -1, nil, fmt.Errorf("beat should be monotonically increasing within a bar")
			}
			if currentBar != beats[i].Bar {
				// if the beat is still monotonically increasing
				// they should be in the same bar
				return -1, -1, nil, fmt.Errorf("beat skipped bars")
			}
		}
	}

	if len(song) == 0 {
		return -1, -1, nil, nil
	}

	tempo := song[0].Tempo

	for i := 1; i < len(song); i++ {
		if (math.Abs(song[i].Tempo - song[i-1].Tempo)) > TempoAllowedVarianceBpm {
			log.Printf("[WARN] bpm should be stable across bars, %f and %f are suspect", song[i].Tempo, song[i-1].Tempo)
		}
		tempo += song[i].Tempo
	}

	tempo /= float64(len(song))

	firstTimestampUs, err := time.ParseDuration(fmt.Sprintf("%fus", beats[0].Timestamp*1000000))
	if err != nil {
		log.Printf("[WARN] couldn't parse first timestamp duration %f for initial sleep", beats[0].Timestamp)
	}

	return tempo, firstTimestampUs, song, nil
}
