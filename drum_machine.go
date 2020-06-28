package main

// #cgo LDFLAGS: /usr/local/lib/libmetro.a /usr/local/lib64/libsoundio.a /usr/local/lib/libstk.a -lstdc++ -lasound -lpulse -ljack -lm
// #include <libmetro/cmetro.h>
import "C"

import (
	"fmt"
	"github.com/sevagh/drum_machine/internal"
	"log"
	"os"
	"time"
)

func main() {
	p := internal.NewParser(os.Stdin)
	beats, err := p.Parse()
	if err != nil {
		log.Fatalf("parse err: %+v\n", err)
	}

	if len(beats) < 1 {
		fmt.Printf("no beats in input file\n")
		return
	}

	tempo, firstTimestampUs, song, err := internal.ValidateBeats(beats)
	if err != nil {
		log.Fatalf("song error: %+v\n", err)
	}

	bpm := C.int(tempo)

	// this is libmetro's C code
	metro := C.metronome_create(bpm)
	defer C.metronome_destroy(metro)

	var ret C.int

	// build a repeating metronome from the first bar
	// TODO: all bars are probably repeating since this is pop
	// but more complicated bars can coexist by discovering their
	// least common multiple and adjusting the bpm and adding silences
	// accordingly. libmetro supports this - see
	// https://github.com/sevagh/libmetro/blob/master/examples/poly_43.cpp#L39

	downbeat := C.note_create_drum_downbeat_1()
	defer C.note_destroy(downbeat)

	beat := C.note_create_drum_beat_1()
	defer C.note_destroy(beat)

	measure := C.measure_create(C.int(song[0].NBeats))
	defer C.measure_destroy(measure)

	C.measure_set_note(measure, C.int(0), downbeat)
	for i := 1; i < song[0].NBeats; i++ {
		if ret = C.measure_set_note(measure, C.int(i), beat); ret != 0 {
			log.Fatalf("error in libmetro C wrapper: %d\n", ret)
		}
	}

	if ret = C.metronome_add_measure(metro, measure); ret != 0 {
		log.Fatalf("error in libmetro C wrapper: %d\n", ret)
	}

	// do an initial sleep of the first timestamp
	time.Sleep(firstTimestampUs)

	fmt.Println("Ctrl-C to exit any time")
	if ret = C.metronome_start_and_loop(metro); ret != 0 {
		log.Fatalf("error in libmetro C wrapper: %d\n", ret)
	}
}
