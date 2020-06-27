package main

// #cgo LDFLAGS: /usr/local/lib/libmetro.a /usr/local/lib64/libsoundio.a /usr/local/lib/libstk.a -lstdc++ -lasound -lpulse -ljack -lm
// #include <libmetro/cmetro.h>
import "C"

import (
	"log"
	"fmt"
)

func LibMetroTest() {
	bpm := C.int(120)
	measureLen := C.int(4)

	metro := C.metronome_create(bpm)
	defer C.metronome_destroy(metro)

	downbeat := C.note_create_drum_downbeat_1()
	defer C.note_destroy(downbeat)

	beat := C.note_create_drum_beat_1()
	defer C.note_destroy(beat)

	measure := C.measure_create(measureLen)
	defer C.measure_destroy(measure)

	var ret C.int

	if ret = C.measure_set_note(measure, C.int(0), downbeat); ret != 0 {
		log.Fatalf("error in libmetro C wrapper: %d\n", ret)
	}

	for i := 1; i < 4; i++ {
		if ret = C.measure_set_note(measure, C.int(i), beat); ret != 0 {
			log.Fatalf("error in libmetro C wrapper: %d\n", ret)
		}
	}

	if ret = C.metronome_add_measure(metro, measure); ret != 0 {
		log.Fatalf("error in libmetro C wrapper: %d\n", ret)
	}

	fmt.Println("Ctrl-C to exit any time")
	if ret = C.metronome_start_and_loop(metro); ret != 0 {
		log.Fatalf("error in libmetro C wrapper: %d\n", ret)
	}
}
