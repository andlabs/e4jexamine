// 2 january 2015
package main

import (
	"fmt"
	"bytes"
	"os"
)

func blockdump(pos int) {
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
	pout := readDataBlock(rr, *escaped)
	n, err := os.Stdout.Write(pout)
	if err != nil {
		panic(err)
	} else if n != blocksize {
		panic(fmt.Errorf("errorless short write to stdout in blockdump (expected %d, wrote %d)", blocksize, n))
	}
}
