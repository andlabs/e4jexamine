// 2 january 2015
package main

import (
	"fmt"
	"os"
	"io"
	"flag"
)

var r io.ReadSeeker

var blocksize int

var (
	u64 = flag.Bool("64", false, "use 64-bit block numbers")
	v3 = flag.Bool("3", false, "use version 3 checksums")
)
func init() {
	flag.IntVar(&blocksize, "bs", 4096, "block size (bytes)")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] command [file]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "commands:\n")
		fmt.Fprintf(os.Stderr, "  summary - print a summary of all journal blocks\n")
	}
	flag.Parse()
	if flag.NArg() != 1 && flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	command := flag.Arg(0)
	if flag.NArg() == 1 {
		r = os.Stdin
	} else {
		f, err := os.Open(flag.Arg(1))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		r = f
	}
	switch {
	case command == "summary":
		summary()
	default:
		fmt.Fprintf(os.Stderr, "unrecognized command %q\n", command)
		flag.Usage()
		os.Exit(1)
	}
}
