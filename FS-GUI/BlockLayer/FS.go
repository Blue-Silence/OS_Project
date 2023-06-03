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

func (fs *BlockFS) ReadFile(inodeN int, index int) DiskLayer.RealBlock {
	inode := fs.INodeN2iNode(inodeN)
	if inode.Valid == false {
		log.Fatal("No Valid inode found for:", inodeN)
	}
	if inode.Pointers[index] < 0 {
		var empty DiskLayer.RealBlock
		return empty
	}
	return fs.VD.ReadBlock(inode.Pointers[index])
}

func (fs *BlockFS) INodeN2iNodeAndPointer(n int) (INode, int) {
	// return Inode itself and the pointer to its block
	if n < 0 {
		return INode{Valid: false}, -1
	}
	iNodemapN := fs.superBlock.INodeMaps[n/Setting.InodePerInodemapBlock]
	//fmt.Println(fs.superBlock.INodeMaps)
	var iNodemap INodeMap = (INodeMap{}.FromBlock((fs.VD.ReadBlock(iNodemapN)))).(INodeMap)
	//fmt.Println("n:", n, "  offset:", iNodemap.Index)
	iNodeBlockN := iNodemap.InodeMapPart[n-iNodemap.Index]
	if iNodemap.Index != n/Setting.InodePerInodemapBlock {
		log.Fatal("Warning!Mistmatch!")
	}

	var nB INodeBlock = INodeBlock{}.FromBlock((fs.VD.ReadBlock(iNodeBlockN))).(INodeBlock)
	for _, v := range nB.NodeArr {
		if v.Valid && v.InodeN == n {
			return v, iNodeBlockN
		}
	}
	return INode{Valid: false}, -1
}

func (fs *BlockFS) INodeN2iNode(n int) INode {
	node, _ := fs.INodeN2iNodeAndPointer(n)
	return node
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
	fs.VD.WriteSuperBlock(super.ToBlocks())
}

func (fs *BlockFS) ReclaimBlock(start int, len int) {
	super := fs.superBlock
	for i := 0; i < len; i++ {
		super.BitMap[start+i] = true
	}
	fs.superBlock = super
	fs.VD.WriteSuperBlock(super.ToBlocks())
}

func (fs *BlockFS) GetIMapPointer(index int) int {
	if index >= 0 && index < Setting.MaxINodemapPartN {
		return fs.superBlock.INodeMaps[index]
	} else {
		return -1
	}
}

func (fs *BlockFS) GetDataBPointer(inodeN int, index int) int {
	inode := fs.INodeN2iNode(inodeN)
	if !inode.Valid {
		return -1
	}
	return inode.Pointers[index]
}

func (fs *BlockFS) GetOneSegHeadStartFrom(start int) int {
	r := -1
	/*for i := len(fs.superBlock.BitMap) - 1; i > 0; i-- {
		if fs.superBlock.BitMap[i] && !(fs.superBlock.BitMap[i-1]) {
			r = i
			break
		}
	}*/
	for i, v := range fs.superBlock.BitMap {
		if i >= start && v {
			r = i
			break
		}
	}
	return r
}

func (fs *BlockFS) Init(VD DiskLayer.VirtualDisk) {
	fs.VD = VD
	for i, _ := range fs.superBlock.INodeMaps {
		fs.superBlock.INodeMaps[i] = -1
	}
	fs.VD.WriteSuperBlock(fs.superBlock.ToBlocks())
}

func (fs *BlockFS) Load(VD DiskLayer.VirtualDisk) {
	fs.VD = VD
	fs.superBlock = fs.superBlock.FromBlocks(VD.ReadSuperBlock())
}

/////////////////////////////////////////////////

func (fs *BlockFS) SuperBlockUNsafe() SuperBlock {
	return fs.superBlock
}
