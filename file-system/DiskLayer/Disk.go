package DiskLayer

import (
	"LSF/Setting"
	"log"
)

type VirtualDisk struct {
	blocks     [Setting.BlockN]Block
	superBlock Block
}

type Block interface {
	CanBeBlock()
}

func (d *VirtualDisk) ReadBlock(index int) Block {
	if index < 0 || index > len(d.blocks) {
		log.Fatal("Invalid disk read access at ", index)
	}
	return d.blocks[index]
}

func (d *VirtualDisk) WriteBlock(index int, b Block) {
	if index < 0 || index > len(d.blocks) {
		log.Fatal("Invalid disk write access at ", index)
	}
	d.blocks[index] = b
}

func (d *VirtualDisk) ReadSuperBlock() Block {
	return d.superBlock
}

func (d *VirtualDisk) WriteSuperBlock(b Block) {
	d.superBlock = b
}
