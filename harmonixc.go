package main

import (
	"fmt"
	"os"
)

func main() {
	p := NewParser(os.Stdin)
	_, err := p.Parse()
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}
}
