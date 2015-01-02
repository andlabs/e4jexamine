// 2 january 2015
package main

import (
	"fmt"
	"bytes"
)

func descdump(pos int) {
	p := make([]byte, blocksize)
	// TODO may not work on pipes
	_, err := r.Seek(int64(pos), 0)
	if err != nil {
		panic(err)
	}
	if readblock(r, p) == false {
		// TODO really?
		return
	}
	rr := bytes.NewReader(p)
	h, b := getblock(rr)
	if err, ok := b.(error); ok && err != nil {
		panic(err)
	}
	if h == nil || h.BlockType != DescriptorBlock {
		// TODO
		panic("block at pos does not represent a descriptor")
	}
	d, _ := readDescriptors(rr, *u64, *v3)
	for _, dd := range d {
		pos += blocksize
		fmt.Printf("0x%08X -> block %d", pos, dd.TargetBlock())
		if (dd.Flags() & Escaped) != 0 {
			fmt.Printf(" [escaped]")
		}
		fmt.Printf("\n")
	}
}
