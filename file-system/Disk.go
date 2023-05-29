package main

import "log"

type VirtualDisk struct {
	blocks     [BlockN]Block
	superBlock SuperBlock
}

func (d *VirtualDisk) readBlock(index int) Block {
	if index < 0 || index > len(d.blocks) {
		log.Fatal("Invalid disk read access at ", index)
	}
	return d.blocks[index]
}

func (d *VirtualDisk) writeBlock(index int, b Block) {
	if index < 0 || index > len(d.blocks) {
		log.Fatal("Invalid disk write access at ", index)
	}
	d.blocks[index] = b
}

func (d *VirtualDisk) readSuperBlock() SuperBlock {
	return d.superBlock
}

func (d *VirtualDisk) writeSuperBlock(b SuperBlock) {
	d.superBlock = b
}
