package AppFSLayer

import (
	"LSF/BlockLayer"
	"LSF/LogLayer"

	//"fmt"
	"log"
)

func (afs *AppFS) createFileWithSpecINodeType(fType int, name string, level int, isRoot bool) int {
	newInodeN := afs.findFreeINode()
	if afs.isINodeInLog(newInodeN) {
		afs.LogCommit()
		//fmt.Println("Ha?")
		newInodeN = afs.findFreeINode()
	} //Avoid reallocating a inode.

	if newInodeN == -1 {
		log.Fatal("No inode number available.")
	}
	in := createInode(fType, name, true, newInodeN)
	in.CurrentLevel = level
	in.IsRoot = isRoot
	afs.tryLog([]BlockLayer.INode{in}, []LogLayer.DataBlockMem{})
	return newInodeN
}

type InodeTrace struct {
	inode  BlockLayer.INode
	offset int
}

func (afs *AppFS) findBlockFromStart(allocateWhenNeed bool, inodeN int, index int) (bool, []InodeTrace, []InodeTrace) {
	b, root, tree, _ := afs.findBlockInTree(allocateWhenNeed, inodeN, index, 0)
	return b, root, tree
}

func (afs *AppFS) findBlockInTree(allocateWhenNeed bool, inodeN int, index int, level int) (bool, []InodeTrace, []InodeTrace, int) {
	inode := afs.GetFileINfo(inodeN)
	if !inode.Valid {
		if allocateWhenNeed {
			inodeN = afs.createFileWithSpecINodeType(BlockLayer.NormalFile, "//", level, true)
			inode = afs.GetFileINfo(inodeN)
		} else {
			return false, []InodeTrace{}, []InodeTrace{}, -1
		}
	}
	//fmt.Println("inode:", inode, "   index:", index)
	if index < blockInInodeLevel(level) {
		b, trs, _ := afs.findBlockInTreeLeaf(allocateWhenNeed, inodeN, index, level)
		if b {
			return true, []InodeTrace{}, trs, inode.InodeN
		} else {
			return false, []InodeTrace{}, []InodeTrace{}, -1
		}
	} else {
		b, rootTrs, treeTrs, np := afs.findBlockInTree(allocateWhenNeed, inode.PointerToNextINode, index-blockInInodeLevel(level), level+1)
		inode := afs.GetFileINfo(inodeN)
		inode.PointerToNextINode = np
		afs.fLog.ConstructLog([]BlockLayer.INode{inode}, []LogLayer.DataBlockMem{})
		return b, append([]InodeTrace{InodeTrace{inode, -1}}, rootTrs...), treeTrs, inode.InodeN
	}

}

func (afs *AppFS) findBlockInTreeLeaf(allocateWhenNeed bool, inodeN int, index int, level int) (bool, []InodeTrace, int) {
	inode := afs.GetFileINfo(inodeN)
	if !inode.Valid {
		if allocateWhenNeed {
			inodeN = afs.createFileWithSpecINodeType(BlockLayer.NormalFile, "///////", level, false)
			inode = afs.GetFileINfo(inodeN)
		} else {
			//fmt.Println("WHAT???")
			return false, []InodeTrace{}, -1
		}
	}
	//fmt.Println("inode leaf:", inode, "   index:", index)
	if inode.CurrentLevel == 0 {
		//if allocateWhenNeed && inode.Pointers[index] <0 {}
		//fmt.Println("AAAAA")
		return true, []InodeTrace{InodeTrace{inode, index}}, inode.InodeN
	} else {
		offset := index / blockInInodeLevel(level-1)
		trace := InodeTrace{inode: inode, offset: offset}

		b, trs, np := afs.findBlockInTreeLeaf(allocateWhenNeed, inode.Pointers[offset], index%blockInInodeLevel(level-1), level-1)
		inode = afs.GetFileINfo(inodeN)
		if inode.Pointers[offset] != np {
			inode.Pointers[offset] = np
			afs.fLog.ConstructLog([]BlockLayer.INode{inode}, []LogLayer.DataBlockMem{})
			//fmt.Println("BB1111B1B1")
		}
		//fmt.Println("BBBBB:", inode, "  inodeN:", inodeN)
		return b, append([]InodeTrace{trace}, trs...), inode.InodeN
	}
}

func blockInInodeLevel(level int) int {
	r := 1
	for i := 0; i <= level; i++ {
		r = r * BlockLayer.DirectPointerPerINode
	}
	return r
}
