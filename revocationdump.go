// 2 january 2015
package main

import (
	"fmt"
	"bytes"
)

func revocationdump(pos int) {
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
	// TODO should this be RevocationRecordBlock?
	if h == nil || h.BlockType != RevocationRecordsBlock {
		// TODO
		panic("block at pos does not represent a revocation record")
	}
	rv := b.(*RevocationRecord)
	blocks := rv.ReadAll(rr, *u64)
	for _, block := range blocks {
		// TODO provide a way to map journal file blocks to device block numbers
		fmt.Printf("%d\n", block)
	}
}
