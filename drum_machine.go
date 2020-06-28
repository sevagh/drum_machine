package main

// #cgo LDFLAGS: /usr/local/lib/libmetro.a /usr/local/lib64/libsoundio.a /usr/local/lib/libstk.a -lstdc++ -lasound -lpulse -ljack -lm
// #include <libmetro/cmetro.h>
import "C"

import (
	"fmt"
	"log"
	"math"
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
		fmt.Printf("err: %+v\n", err)
	}

	if len(beats) < 1 {
		fmt.Printf("no beats in input file\n")
		return
	}

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
		if currentBeat < beats[i-1].Beat {
			// check if bars are monotonically increasing
			if beats[i].Bar != currentBar+1 {
				log.Fatal("bars should be monotonically increasing")
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
				log.Fatal("beat should be monotonically increasing within a bar")
			}
			if currentBar != beats[i].Bar {
				// if the beat is still monotonically increasing
				// they should be in the same bar
				log.Fatal("beat skipped bars")
			}
		}
	}

	tempo := song[0].Tempo

	for i := 1; i < len(song); i++ {
		if (math.Abs(song[i].Tempo - song[i-1].Tempo)) > TempoAllowedVarianceBpm {
			log.Printf("[WARN] bpm should be stable across bars, %f and %f are suspect", song[i].Tempo, song[i-1].Tempo)
		}
		tempo += song[i].Tempo
	}

	tempo /= float64(len(song))
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

	firstTimestampUs, err := time.ParseDuration(fmt.Sprintf("%fus", beats[0].Timestamp*1000000))
	if err != nil {
		log.Printf("[WARN] couldn't parse first timestamp duration %f for initial sleep", beats[0].Timestamp)
	}

	// do an initial sleep of the first timestamp
	time.Sleep(firstTimestampUs)

	fmt.Println("Ctrl-C to exit any time")
	if ret = C.metronome_start_and_loop(metro); ret != 0 {
		log.Fatalf("error in libmetro C wrapper: %d\n", ret)
	}
}
