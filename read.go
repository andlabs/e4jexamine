// 2 january 2015
package main

import (
	"fmt"
	"io"
	"encoding/binary"
)

func read(r io.Reader, s interface{}) {
	err := binary.Read(r, Endian, s)
	if err != nil {
		panic(err)
	}
}

func getblock(r io.Reader) (header *Header, block interface{}) {
	header = new(Header)
	read(r, header)
	if header.Magic != BlockMagic {
		return nil, nil
	}
	switch header.BlockType {
	case DescriptorBlock:
		return header, nil
	case CommitRecordBlock:
		cb := new(CommitRecord)
		read(r, cb)
		return header, cb
	case SuperblockV1Block:
		sb := new(SuperblockV1)
		read(r, sb)
		return header, sb
	case SuperblockV2Block:
		sb := new(SuperblockV2)
		read(r, sb)
		return header, sb
	case RevocationRecordsBlock:
		rr := new(RevocationRecord)
		read(r, rr)
		return header, rr
	}
	return nil, fmt.Errorf("unrecognized block type %d", header.BlockType)
}

func (r *RevocationRecord) ReadAll(rd io.Reader, u64 bool) []uint64 {
	rr := make([]uint64, r.Count)
	if u64 {
		n := uint64(0)
		for i, _ := range rr {
			read(rd, &n)
			rr[i] = n
		}
	} else {
		n := uint32(0)
		for i, _ := range rr {
			read(rd, &n)
			rr[i] = uint64(n)
		}
	}
	return rr
}

func readDescriptors(r io.Reader, u64 bool, v3 bool) (d []Descriptor, uuids []*[16]byte) {
	var dv Descriptor

	for {
		switch {
		case !u64 && !v3:
			dd := new(DescriptorNonV3)
			read(r, dd)
			dv = dd
		case !u64 && v3, u64 && v3:
			dd := new(DescriptorV3)
			read(r, dd)
			dv = dd
		case u64 && !v3:
			dd := new(DescriptorNonV364)
			read(r, dd)
			dv = dd
		}
		d = append(d, dv)
		if (dv.Flags() & SameUUID) == 0 {
			uuid := new([16]byte)
			read(r, uuid)
			uuids = append(uuids, uuid)
		} else {
			uuids = append(uuids, nil)
		}
		if (dv.Flags() & LastTag) != 0 {
			break
		}
	}
	return d, uuids
}

func readDataBlock(r io.Reader, d Descriptor) []byte {
	p := make([]byte, blocksize)
	n, err := r.Read(p)
	if err != nil {
		panic(err)
	} else if n != blocksize {
		panic(fmt.Errorf("errorless short read of data block (expected %d, got %d)", blocksize, n))
	}
	if (d.Flags() & Escaped) != 0 {
		Endian.PutUint32(p, BlockMagic)
	}
	return p
}
