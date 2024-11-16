#!/bin/bash

if [[ -z ${GOPATH} ]] ; then
  echo 'run this from a gopath with built blueprintread cmd'
  exit 1
fi

if [[ ! -e ${GOPATH}/bin/blueprintread ]] ; then
  echo 'run: go install badc0de.net/pkg/factorioblueprint/blueprintread'
  exit 1
fi

if [[ ! -e ../../../github.com/Redcrafter/verilog2factorio ]] ; then
  mkdir -p ../../../github.com/Redcrafter/
  (
    cd ../../../github.com/Redcrafter/
    git clone github.com/Redcrafter/verilog2factorio
    cd verilog2factorio/
    npm install
  )
fi

if [[ ! -e /usr/bin/yosys ]] ; then
  sudo apt install yosys
  # 0.34 or earlier, according to readme of v2f.
  # current debian, bookworm, has 0.23
fi

# we could pipe directly, but let's keep a debug file around
(cd ../../../github.com/Redcrafter/verilog2factorio ;  ./v2f -g chunkAnnealing ./samples/counter.v | tail -n+4) > /tmp/blueprint.txt
blueprintread "$@" < /tmp/blueprint.txt

