// 2 january 2015
package main

import (
	"fmt"
	"bytes"
)

func commitdump(pos int) {
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
	if h == nil || h.BlockType != CommitRecordBlock {
		// TODO
		panic("block at pos does not represent a commit record")
	}
	c := b.(*CommitRecord)
	fmt.Printf("commit time: %v\n", c.Time())
}
