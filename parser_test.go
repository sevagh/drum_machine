package main

import (
	"strings"
	"testing"
	"time"
)

func TestParserInvalidBeat1(t *testing.T) {
	testString := `1234
5678`
	rdr := strings.NewReader(testString)
	p := NewParser(rdr)
	t.Logf("p: %+v\n", p)

	beats, err := p.Parse()
	if err == nil {
		t.Fatalf("bad")
	}
	if len(beats) != 0 {
		t.Fatalf("bad")
	}
}

func TestParserValidBeatsSensicalSong(t *testing.T) {
	testString := `55 1 1
56 2 1
`

	rdr := strings.NewReader(testString)
	p := NewParser(rdr)
	t.Logf("p: %+v\n", p)

	beats, err := p.Parse()
	t.Logf("beats: %+v, err: %+v\n", beats, err)

	if err != nil {
		t.Fatalf("bad: %+v", err)
	}
	if len(beats) != 2 {
		t.Fatalf("bad")
	}

	if beats[0].Timestamp != 55.0 {
		t.Fatal("bad")
	}
	if beats[0].Beat != 1 {
		t.Fatal("bad")
	}
	if beats[0].Bar != 1 {
		t.Fatal("bad")
	}
	if beats[1].Timestamp != 56.0 {
		t.Fatal("bad")
	}
	if beats[1].Beat != 2 {
		t.Fatal("bad")
	}
	if beats[1].Bar != 1 {
		t.Fatal("bad")
	}

	tempo, firstTstamp, song, err := validateBeats(beats)
	if err != nil {
		t.Fatalf("bad")
	}
	if tempo != 30 {
		t.Fatalf("bad")
	}
	if firstTstamp != 55000000*time.Microsecond {
		t.Fatalf("bad")
	}
	// just one bar
	if len(song) != 1 {
		t.Logf("len song: %d\n", len(song))
		t.Fatalf("bad")
	}
}

func TestParserValidBeatsNonsensicalSong(t *testing.T) {
	// jump in bar
	testString := `55 1 1
56 2 3
`
	rdr := strings.NewReader(testString)
	p := NewParser(rdr)
	t.Logf("p: %+v\n", p)

	beats, err := p.Parse()
	t.Logf("beats: %+v, err: %+v\n", beats, err)

	if err != nil {
		t.Fatalf("bad: %+v", err)
	}
	if len(beats) != 2 {
		t.Fatalf("bad")
	}

	tempo, firstTstamp, song, err := validateBeats(beats)
	if err == nil {
		t.Fatalf("bad")
	}
	if tempo != -1 {
		t.Fatalf("bad")
	}
	if firstTstamp != -1 {
		t.Fatalf("bad")
	}
	if len(song) != 0 {
		t.Fatalf("bad")
	}
}
