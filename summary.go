// 2 january 2015
package main

import (
	"fmt"
	"bytes"
)

func summary() {
	var pos int

	p := make([]byte, blocksize)
	for pos = 0; readblock(r, p); pos += blocksize {
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
