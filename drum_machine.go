package main

// #cgo LDFLAGS: /usr/local/lib/libmetro.a /usr/local/lib64/libsoundio.a /usr/local/lib/libstk.a -lstdc++ -lasound -lpulse -ljack -lm
// #include <libmetro/cmetro.h>
import "C"

import (
	"fmt"
	"log"
	"os"
	"time"
)

// 125bpm vs 120bpm = 0.48 vs 0.5
const TempoAllowedVarianceSecs float64 = 0.02
const TempoAllowedVarianceBpm float64 = 5

func main() {
	p := NewParser(os.Stdin)
	beats, err := p.Parse()
	if err != nil {
		log.Fatalf("parse err: %+v\n", err)
	}

	if len(beats) < 1 {
		fmt.Printf("no beats in input file\n")
		return
	}

	tempo, firstTimestampUs, song, err := validateBeats(beats)
	if err != nil {
		log.Fatalf("song error: %+v\n", err)
	}

	bpm := C.int(tempo)

	// this is libmetro's C code
	metro := C.metronome_create(bpm)
	defer C.metronome_destroy(metro)

	var ret C.int

	for _, bar := range song {
		downbeat := C.note_create_drum_downbeat_1()
		defer C.note_destroy(downbeat)

		beat := C.note_create_drum_beat_1()
		defer C.note_destroy(beat)

		measure := C.measure_create(C.int(bar.NBeats))
		defer C.measure_destroy(measure)

		if ret = C.measure_set_note(measure, C.int(0), downbeat); ret != 0 {
			log.Fatalf("error in libmetro C wrapper: %d\n", ret)
		}

		for i := 1; i < bar.NBeats; i++ {
			if ret = C.measure_set_note(measure, C.int(i), beat); ret != 0 {
				log.Fatalf("error in libmetro C wrapper: %d\n", ret)
			}
		}

		if ret = C.metronome_add_measure(metro, measure); ret != 0 {
			log.Fatalf("error in libmetro C wrapper: %d\n", ret)
		}
	}

	// do an initial sleep of the first timestamp
	time.Sleep(firstTimestampUs)

	fmt.Println("Ctrl-C to exit any time")
	if ret = C.metronome_start_and_loop(metro); ret != 0 {
		log.Fatalf("error in libmetro C wrapper: %d\n", ret)
	}
}
