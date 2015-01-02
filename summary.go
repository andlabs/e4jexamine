// 2 january 2015
package main

import (
	"fmt"
	"os"
	"io"
	"bytes"
)

var r io.Reader

const blocksize = 4096

func main() {
	var pos int64

	if len(os.Args) == 0 {
		r = os.Stdin
	} else {
		f, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer f.Close()
		r = f
	}
	p := make([]byte, blocksize)
	for pos = 0; ; pos += blocksize {
		n, err := r.Read(p)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else if n != blocksize {
			panic(fmt.Errorf("errorless short read reading block from journal file (expected %d got %d)", blocksize, n))
		}
		rr := bytes.NewReader(p)
		h, b := getblock(rr)
		if h == nil {
			if b == nil {
				fmt.Printf("0x%08X - invalid/data block\n", pos)
				continue
			}
			err := b.(error)
			panic(err)
		}
		fmt.Printf("0x%08X - ", pos)
		switch bb := b.(type) {
		case nil:		// descriptors
			fmt.Printf("descriptors\n")
		case *SuperblockV1:
			fmt.Printf("v1 superblock\n")
		case *SuperblockV2:
			fmt.Printf("v2 superblock\n")
		case *CommitRecord:
			fmt.Printf("commit record\n")
		case *RevocationRecord:
			fmt.Printf("revocation record\n")
			_ = bb
		}
	}
}
