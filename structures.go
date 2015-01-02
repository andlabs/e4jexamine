// 2 january 2015
package main

import (
	"fmt"
	"encoding/binary"
	"hash"
	"hash/crc32"
	"crypto/md5"
	"crypto/sha1"
	"time"
)

var Endian = binary.BigEndian

type Header struct {
	Magic		uint32
	BlockType	BlockType
	TransactionID	uint32
}

const BlockMagic uint32 = 0xC03B3998

//go:generate stringer -type BlockType
type BlockType uint32
const (
	DescriptorBlock BlockType = 1
	CommitRecordBlock BlockType = 2
	SuperblockV1Block BlockType = 3
	SuperblockV2Block BlockType = 4
	RevocationRecordsBlock BlockType = 5
)

type SuperblockV1 struct {
	BlockSize			uint32
	BlockCount		uint32
	FirstLogBlock		uint32
	FirstCommitID		uint32
	StartingLogBlock	uint32
	Errno			uint32
}

type SuperblockV2 struct {
	SuperblockV1
	FeatureCompat			FeatureCompat
	FeatureIncompat		FeatureIncompat
	FeatureReadOnlyCompat	FeatureReadOnlyCompat
	UUID				[16]byte
	FilesystemCount		uint32
	DynamicSuperblockPos	uint32
	MaxTransactions		uint32
	MaxBlocksPerTransaction		uint32
	ChecksumType				ChecksumType
	_						[3]byte
	_						[42]uint32
	Checksum				uint32
	Filesystems				[16 * 48]byte
}

//go:generate stringer -type FeatureCompat
type FeatureCompat uint32
const (
	MaintainsChecksums FeatureCompat = 1
)

//go:generate stringer -type FeatureIncompat
type FeatureIncompat uint32
const (
	HasBlockRevocationRecords FeatureIncompat = 1
	Handles64BitBlockNumbers FeatureIncompat = 2
	CommitsAsynchronously FeatureIncompat = 4
	UsesV2Checksums FeatureIncompat = 8
	UsesV3Checksums FeatureIncompat = 0x10
)

//go:generate stringer -type FeatureReadOnlyCompat
type FeatureReadOnlyCompat uint32
// no defined values yet
const MakeStringerHappy FeatureReadOnlyCompat = 0

//go:generate stringer -type ChecksumType
type ChecksumType byte
const (
	CRC32 ChecksumType = 1
	MD5 ChecksumType = 2
	SHA1 ChecksumType = 3
	CRC32C ChecksumType = 4
)

// via https://ext4.wiki.kernel.org/index.php/Ext4_Metadata_Checksums#Algorithm
const CRC32Polynomial = 0x04C11DB7
var CRC32Table = crc32.MakeTable(CRC32Polynomial)
const CRC32CPolynomial = 0x1EDC6F41
var CRC32CTable = crc32.MakeTable(CRC32CPolynomial)

func (c ChecksumType) New() hash.Hash {
	switch c {
	case CRC32:
		return crc32.New(CRC32Table)
	case MD5:
		return md5.New()
	case SHA1:
		return sha1.New()
	case CRC32C:
		return crc32.New(CRC32CTable)
	}
	panic(fmt.Errorf("unrecognized checksum type %d", c))
}

type DescriptorNonV3 struct {
	TargetBlockLow32	uint32
	Checksum		uint16
	DescriptorFlags	DescriptorFlags
}

func (d *DescriptorNonV3) TargetBlock() uint64 {
	return uint64(d.TargetBlockLow32)
}

func (d *DescriptorNonV3) Flags() DescriptorFlags {
	return d.DescriptorFlags
}

//go:generate stringer -type DescriptorFlags
type DescriptorFlags uint16
const (
	Escaped DescriptorFlags = 1		// block actually begins with BlockMagic so we need to watch for that when dumping
	SameUUID = 2
	Deleted = 4
	LastTag = 8
)

type DescriptorNonV364 struct {
	DescriptorNonV3
	TargetBlockHigh32	uint32
}

func (d *DescriptorNonV364) TargetBlock() uint64 {
	return (uint64(d.TargetBlockHigh32) << 32) | uint64(d.TargetBlockLow32)
}

func (d *DescriptorNonV364) Flags() DescriptorFlags {
	return d.DescriptorFlags
}

type DescriptorV3 struct {
	TargetBlockLow32	uint32
	DescriptorFlags	DescriptorFlags
	TargetBlockHigh32	uint32			// must be zero if no 64-bit support, so we're good to always use this
	Checksum		uint32
}

func (d *DescriptorV3) TargetBlock() uint64 {
	return (uint64(d.TargetBlockHigh32) << 32) | uint64(d.TargetBlockLow32)
}

func (d *DescriptorV3) Flags() DescriptorFlags {
	return d.DescriptorFlags
}

type Descriptor interface {
	TargetBlock() uint64
	Flags() DescriptorFlags
}

type RevocationRecord struct {
	Count	uint32
}

type CommitRecord struct {
	ChecksumType		ChecksumType
	ChecksumSize		byte
	_				[2]byte
	Checksums		[8]uint32		// TODO constant
	Timestamp		uint64
	TimestampNano	uint32
}

func (c *CommitRecord) Time() time.Time {
	return time.Unix(int64(c.Timestamp), int64(c.TimestampNano))
}
