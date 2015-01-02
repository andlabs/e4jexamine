// generated by stringer -type ChecksumType; DO NOT EDIT

package main

import "fmt"

const _ChecksumType_name = "CRC32MD5SHA1CRC32C"

var _ChecksumType_index = [...]uint8{0, 5, 8, 12, 18}

func (i ChecksumType) String() string {
	i -= 1
	if i+1 >= ChecksumType(len(_ChecksumType_index)) {
		return fmt.Sprintf("ChecksumType(%d)", i+1)
	}
	return _ChecksumType_name[_ChecksumType_index[i]:_ChecksumType_index[i+1]]
}
