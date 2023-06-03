package DiskLayer

import (
	"LSF/Setting"
	"log"
)

type VirtualDisk interface {
	ReadBlock(index int) RealBlock
	WriteBlock(index int, b Block)
	ReadSuperBlock() []RealBlock
	WriteSuperBlock(b []RealBlock)
}

type Block interface {
	CanBeBlock()
	ToBlock() RealBlock
	FromBlock(RealBlock) Block
}

type RealBlock = [Setting.BlockSize]byte

func BytesToBlock(d []byte) RealBlock {
	var b RealBlock
	if len(d) > Setting.BlockSize {
		log.Fatal("Too big to be a block.")
	}
	copy(b[:], d)
	return b
}

func BytesToBlocks(d []byte) []RealBlock {
	dN := d[:]
	bs := []RealBlock{}
	for len(dN) > Setting.BlockSize {
		bs = append(bs, BytesToBlock(dN[:Setting.BlockSize]))
		dN = dN[Setting.BlockSize:]
	}
	if len(dN) > 0 {
		bs = append(bs, BytesToBlock(dN[:Setting.BlockSize]))
	}
	return bs
}

func BlockToBytes(b RealBlock) []byte {
	return b[:]
}
