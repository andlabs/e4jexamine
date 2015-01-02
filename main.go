// 2 january 2015
package main

import (
	"fmt"
	"os"
	"io"
	"flag"
	"strings"
	"strconv"
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

func badline(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n")
	flag.Usage()
	os.Exit(1)
}

func getpos(command string) int64 {
	m := strings.Split(command, ".")
	if len(m) != 2 {
		badline("invalid command specification: %s", command)
	}
	n, err := strconv.ParseInt(m[1], 0, 64)
	if err != nil {
		badline("error parsing command number %q: %v", m[1], err)
	}
	return n
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] command [file]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "commands:\n")
		fmt.Fprintf(os.Stderr, "  summary - print a summary of all journal blocks\n")
		fmt.Fprintf(os.Stderr, "  descdump.nnn - print a summary of the descriptor block at nnn and respective data blocks\n")
		fmt.Fprintf(os.Stderr, "all nnn values are BYTE OFFSETS and may be octal with a leading 0 or hexadecimal with a leading 0x or 0X; decimal otherwise\n")
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
	case strings.HasPrefix(command, "descdump."):
		pos := getpos(command)
		fmt.Println(pos)
	default:
		badline("unrecognized command %q", command)
	}
}
