package internal

import "bytes"

const (
	fuzzInteresting = 1
	fuzzNormal      = 0
	fuzzDiscard     = -1
)

func FuzzDrummachine(in []byte) int {
	p := NewParser(bytes.NewReader(in))

	_, err := p.Parse()
	if err == nil {
		return fuzzInteresting
	}

	return fuzzNormal
}
