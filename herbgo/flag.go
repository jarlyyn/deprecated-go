package main

import (
	"flag"
	"os"
)

var Args = flag.NewFlagSet(os.Args[0], flag.PanicOnError)

func ParseArgs() {
	if !Args.Parsed() {
		Args.Parse(os.Args[2:])
	}
}
