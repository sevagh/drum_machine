#!/usr/bin/env bash

set -euxo pipefail

go-fuzz-build -func FuzzDrummachine -o FuzzDrummachine.zip github.com/sevagh/drum_machine/internal

go-fuzz -bin FuzzDrummachine.zip -workdir harmonixset-corpus/
