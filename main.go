package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type (
	Content struct {
		text []byte
	}
)

var decode bool

func init() {
	flag.BoolVar(&decode, "d", false, "Decode mode")
	flag.Parse()
}

func isArgFileRef(a string) bool {
	return strings.HasPrefix(a, "@") && !strings.HasPrefix(a, "\\@")
}

func isArgStdinRef(a string) bool {
	return a == "-"
}

func main() {
	var args = os.Args
	switch {
	case len(args) < 1:
		log.Fatalf("Not enough arguments (expected >1, got %d)", len(args))
	case len(args) == 1 && !isArgStdinRef(args[0]) && !isArgFileRef(args[0]):
		log.Fatal("Must supply data")
	}
	var err error
	var c Content
	var arg = args[1]
	println(arg)
	switch {
	// Handle first argument: can be -, @file, «var path»
	case isArgStdinRef(arg):
		// read the specification into memory from stdin
		if c.text, err = io.ReadAll(os.Stdin); err != nil {
			log.Fatalf("Error reading from stdin: %s", err)
		}
	case isArgFileRef(arg):
		// ArgFileRefs start with "@" so we need to peel that off
		// detect format based on file extension
		specPath := arg[1:]
		if c.text, err = os.ReadFile(specPath); err != nil {
			log.Fatalf(fmt.Sprintf("Error reading %q: %s", specPath, err))
		}
	}

	if decode {
		decodedText, err := base64.StdEncoding.DecodeString(string(c.text))
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't decode due to %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stdout, string(decodedText))
	} else {
		encodedText := base64.StdEncoding.EncodeToString(c.text)
		fmt.Fprintln(os.Stdout, encodedText)
	}
}
