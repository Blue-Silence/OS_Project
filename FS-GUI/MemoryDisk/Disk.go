package MemoryDisk

import (
	"LSF/DiskLayer"
	"LSF/Setting"
	"log"
	"runtime/debug"
)

type RamDisk struct {
	//blocks     [Setting.BlockN]Block
	blocks [Setting.BlockN]*DiskLayer.RealBlock
	//superBlock Block
	//superBlock RealBlock
	superBlock []*DiskLayer.RealBlock
}

func (d *RamDisk) ReadBlock(index int) *DiskLayer.RealBlock {
	if index < 0 || index > len(d.blocks) {
		debug.PrintStack()
		log.Fatal("Invalid disk read access at ", index)
	}
	return d.blocks[index]
}

func (d *RamDisk) WriteBlock(index int, b *DiskLayer.Block) {
	if index < 0 || index > len(d.blocks) {
		debug.PrintStack()
		log.Fatal("Invalid disk write access at ", index)
	}
	d.blocks[index] = (*b).ToBlock()
}

func (d *RamDisk) ReadSuperBlock() []*DiskLayer.RealBlock {
	return d.superBlock
}

func (d *RamDisk) WriteSuperBlock(b []*DiskLayer.RealBlock) {
	d.superBlock = b
}
