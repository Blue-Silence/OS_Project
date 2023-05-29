package main

import "log"

type FS struct {
	superBlock SuperBlock
	VD         VirtualDisk
} //This is the file system data structure in memory

func (fs *FS) formatFS(VD VirtualDisk) {
	fs.VD = VD
	var s SuperBlock

	var rootINode INode
	rootINode.valid = true
	rootINode.fileType = Folder
	rootINode.name = ""
	rootINode.inodeN = 0
	for i, _ := range rootINode.pointers {
		rootINode.pointers[i] = -1
	}

	var initSegHead SegHead
	initSegHead.inodeMapN = MaxINodemapPartN
	initSegHead.inodeBlockN = MaxInodeN / INodePerBlock
	initSegHead.dataBlockN = 0
	initSegHead.segLen = initSegHead.inodeMapN + initSegHead.inodeBlockN + initSegHead.dataBlockN + 1 // All inode map and inode is in this seg

	for i := 0; i < MaxINodemapPartN; i++ {
		var m INodeMap
		m.offset = i * InodePerInodemapBlock
		for i, _ := range m.inodeMapPart {
			iNodeN := m.offset + i
			m.inodeMapPart[i] = iNodeN/INodePerBlock + 1 + MaxINodemapPartN
		}
		fs.VD.writeBlock(1+i, m)
	} // Init inode map part

	for i := 0; i < initSegHead.inodeBlockN; i++ {
		var b INodeBlock
		for j, _ := range b.nodeArr {
			b.nodeArr[j] = INode{valid: false, inodeN: i*INodePerBlock + j}
		}
		fs.VD.writeBlock(1+initSegHead.inodeMapN+i, b)
	}

	//s.bitMap =
	for i, _ := range s.iNodeMaps {
		s.iNodeMaps[i] = 1 + i
	}
	fs.superBlock = s
	fs.VD.writeSuperBlock(s)

}

func (fs *FS) readFile(inodeN int, index int) Block {
	inode := fs.iNodeN2iNode(inodeN)
	if inode.valid == false {
		log.Fatal("No valid inode found for:", inodeN)
	}
	return fs.VD.readBlock(inode.pointers[index])
}

func (fs *FS) iNodeN2iNode(n int) INode {
	iNodemapN := fs.superBlock.iNodeMaps[n/InodePerInodemapBlock]
	var iNodemap INodeMap = fs.VD.readBlock(iNodemapN)
	iNodeBlockN := iNodemap.inodeMapPart[n-iNodemap.offset]
	if iNodemap.offset != n/InodePerInodemapBlock {
		log.Fatal("Warning!Mistmatch!")
	}

	var nB INodeBlock = fs.VD.readBlock(iNodeBlockN)
	for _, v := range nB.nodeArr {
		if v.valid && v.inodeN == n {
			return v
		}
	}
	return INode{valid: false}
}

func (fs *FS) findSpaceForSeg(len int) int {

	segStart := -1
	count := 0
	for i, v := range fs.superBlock.bitMap {
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

func (fs *FS) writeFile(inodeN int, index int, data Block) {
	inode := fs.iNodeN2iNode(inodeN)
	if inode.valid == false {
		log.Fatal("No valid inode found for:", inodeN)
	}
	var segHead SegHead
	segHead.dataBlockN = 1
	segHead.inodeMapN = 1
	segHead.inodeBlockN = 1
	segStart := fs.findSpaceForSeg(segHead.len())
	if segStart == -1 {
		log.Fatal("No space!")
	} // We will do gc later.
	segHead.dataBlockSummary[0] = FileIndexInfo{inodeN: inodeN, index: index}

}
