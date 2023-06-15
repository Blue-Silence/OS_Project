package BlockLayer

import (
	"LSF/DiskLayer"
	"LSF/Setting"
	"log"
)

type BlockFS struct {
	superBlock SuperBlock
	VD         DiskLayer.VirtualDisk
} //This is the file system data structure in memory

func (fs *BlockFS) ReadFile(inodeN int, index int) DiskLayer.Block {
	inode := fs.INodeN2iNode(inodeN)
	if inode.Valid == false {
		log.Fatal("No Valid inode found for:", inodeN)
	}
	return fs.VD.ReadBlock(inode.Pointers[index])
}

func (fs *BlockFS) INodeN2iNode(n int) INode {
	iNodemapN := fs.superBlock.INodeMaps[n/Setting.InodePerInodemapBlock]
	var iNodemap INodeMap = (fs.VD.ReadBlock(iNodemapN)).(INodeMap)
	iNodeBlockN := iNodemap.InodeMapPart[n-iNodemap.Offset]
	if iNodemap.Offset != n/Setting.InodePerInodemapBlock {
		log.Fatal("Warning!Mistmatch!")
	}

	var nB INodeBlock = (fs.VD.ReadBlock(iNodeBlockN)).(INodeBlock)
	for _, v := range nB.NodeArr {
		if v.Valid && v.InodeN == n {
			return v
		}
	}
	return INode{Valid: false}
}

func (fs *BlockFS) FindSpaceForSeg(len int) int {
	segStart := -1
	count := 0
	for i, v := range fs.superBlock.BitMap {
		if !v {
			count++
			if count == len {
				segStart = i - (len - 1)
				break
			}
		} else {
			count = 0
		}
	}
	return segStart
}

func (fs *BlockFS) ApplyUpdate(start int, bs []DiskLayer.Block, newIMapLocation map[int]int) {
	super := fs.superBlock
	for i, v := range newIMapLocation {
		super.INodeMaps[i] = v
	} // Update imap address
	for i, v := range bs {
		fs.VD.WriteBlock(i+start, v)
		super.BitMap[start+i] = true
	}
	fs.superBlock = super
	fs.VD.WriteSuperBlock(super)
}

func (fs *BlockFS) ReclaimBlock(start int, len int) {
	super := fs.superBlock
	for i := 0; i < len; i++ {
		super.BitMap[start+i] = true
	}
	fs.superBlock = super
	fs.VD.WriteSuperBlock(super)
}
