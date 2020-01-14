#!/bin/sh

# Set development mode
export DEVMODE=1

# build and run (requires silver searcher and entr)
ag -l | entr -drs "echo --- && golint && go build -race && ./matlistan"