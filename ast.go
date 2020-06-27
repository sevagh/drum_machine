package main

import "time"

type Beat struct {
	Timestamp time.Duration
	Beat      int
	Bar       int
}

type Bar []*Beat

type Song []*Bar
