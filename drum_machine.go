package main

import (
	"fmt"
	"os"
)

func main() {
	p := NewParser(os.Stdin)
	beats, err := p.Parse()
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}
	fmt.Printf("beats: %+v\n", beats)
	fmt.Printf("len beats: %d\n", len(beats))
}
