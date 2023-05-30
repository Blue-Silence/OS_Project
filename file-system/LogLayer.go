package main

type DataBlockMem struct {
	inode int
	index int
	data  Block
}

type FSLog struct {
	data        map[int]([]DataBlockMem)
	inodeByImap map[int]([]INode)
	//len         int
} // This is in memory

func (L *FSLog) initLog() {
	L.data = make(map[int][]DataBlockMem)
	L.inodeByImap = make(map[int][]INode)
	//L.len = 0
}

func (L *FSLog) constructLog(inodes []INode, ds []DataBlockMem) bool {
	_, _, dataBlockN, _ := L.lenInBlock()
	if len(ds)+dataBlockN > MaxEditBlcokN {
		return false
	}
	for _, v := range inodes {
		L.inodeByImap[v.inodeN/InodePerInodemapBlock] = append(L.inodeByImap[v.inodeN/InodePerInodemapBlock], v)
	}
	for _, v := range ds {
		L.data[v.inode] = append(L.data[v.inode], v)
	}
	return true
}

func (L *FSLog) lenInBlock() (int, int, int, int) {
	inodeBlockN := 0
	inodeMapN := 0
	dataBlockN := 0
	segLen := 0
	for _, v := range L.inodeByImap {
		inodeBlockN += len(v)
		inodeMapN++
	}
	inodeBlockN = inodeBlockN / INodePerBlock

	for _, v := range L.data {
		dataBlockN += len(v)
	}

	segLen = inodeMapN + inodeBlockN + dataBlockN + 1
	return inodeMapN, inodeBlockN, dataBlockN, segLen
}

func (L *FSLog) log2DiskBlock(start int, inodeMap map[int]INodeMap) ([]Block, map[int]int) {
	var segHead SegHead
	segHead.inodeMapN, segHead.inodeBlockN, segHead.dataBlockN, _ = L.lenInBlock()

	var dataBlock []Block

	for _, v := range L.inodeByImap {
		for _, n := range v {
			for _, dataB := range L.data[n.inodeN] {
				segHead.dataBlockSummary[len(dataBlock)-1] = FileIndexInfo{n.inodeN, dataB.index}
				dataBlock = append(dataBlock, dataB.data)
				n.pointers[dataB.index] = start + 1 + segHead.inodeMapN + segHead.inodeBlockN + len(dataBlock) - 1
			}
		}
	}

	var nodesByBlock []INodeBlock

	nodeCount := INodePerBlock
	for _, v := range L.inodeByImap {
		for _, n := range v {
			if nodeCount == INodePerBlock {
				nodesByBlock = append(nodesByBlock, INodeBlock{})
				nodeCount = 0
			}
			nodesByBlock[len(nodesByBlock)-1].nodeArr[nodeCount] = n
			nodeCount++
			//and also do something to change imap next (TO BE DONE)
			iPart := inodeMap[n.inodeN/InodePerInodemapBlock]
			(iPart.inodeMapPart)[n.inodeN%InodePerInodemapBlock] = len(nodesByBlock) - 1 + start + 1 + segHead.inodeMapN
			inodeMap[n.inodeN/InodePerInodemapBlock] = iPart
		}
	}

	var imapBlock []INodeMap
	returnMap := make(map[int]int)
	for _, v := range inodeMap {
		imapBlock = append(imapBlock, v)
		returnMap[v.offset/InodePerInodemapBlock] = start + 1 + len(imapBlock) - 1
	}

	re := []Block{segHead}
	for _, v := range imapBlock {
		re = append(re, v)
	}
	for _, v := range nodesByBlock {
		re = append(re, v)
	}
	for _, v := range dataBlock {
		re = append(re, v)
	}
	return re, returnMap
}

func (L *FSLog) imapNeeded() []int {
	var re []int
	for i, _ := range L.inodeByImap {
		re = append(re, i)
	}
	return re
}
