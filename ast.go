package main

type Beat struct {
	Timestamp float64
	Beat      int
	Bar       int
}

type Bar struct {
	Tempo  float64
	NBeats int
}

type Song []Bar
