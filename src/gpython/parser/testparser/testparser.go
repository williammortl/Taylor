// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/williammortl/Taylor/src/gpython/compile"
	"github.com/williammortl/Taylor/src/gpython/parser"
)

var (
	lexFile     = flag.Bool("l", false, "Lex the file only")
	compileFile = flag.Bool("c", false, "Lex, Parse and compile the file")
	debugLevel  = flag.Int("d", 0, "Debug level 0-4")
)

func main() {
	flag.Parse()
	parser.SetDebug(*debugLevel)
	if len(flag.Args()) == 0 {
		log.Printf("Need files to parse")
		os.Exit(1)
	}
	for _, path := range flag.Args() {
		if *lexFile {
			fmt.Printf("Lexing %q\n", path)
		} else if *compileFile {
			fmt.Printf("Compiling %q\n", path)
		} else {
			fmt.Printf("Parsing %q\n", path)
		}
		in, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		if *debugLevel > 0 {
			fmt.Printf("-----------------\n")
		}
		if *lexFile {
			_, err = parser.Lex(in, path, "exec")
		} else if *compileFile {
			var input []byte
			input, err = ioutil.ReadAll(in)
			if err != nil {
				log.Fatalf("Failed to read %q: %v", path, err)
			}
			_, err = compile.Compile(string(input), path, "exec", 0, false)
		} else {
			_, err = parser.Parse(in, path, "exec")
		}
		if *debugLevel > 0 {
			fmt.Printf("-----------------\n")
		}
		closeErr := in.Close()
		if err != nil {
			log.Fatalf("Failed on %q: %v", path, err)
		}
		if closeErr != nil {
			log.Fatalf("Failed to close %q: %v", path, closeErr)
		}
	}
}
