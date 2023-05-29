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

func (L *FSLog) constructLog(inodes []INode, ds []DataBlockMem) {
	for _, v := range inodes {
		L.inodeByImap[v.inodeN/InodePerInodemapBlock] = append(L.inodeByImap[v.inodeN/InodePerInodemapBlock], v)
	}
	for _, v := range ds {
		L.data[v.inode] = append(L.data[v.inode], v)
	}

	/*length := 0
	for _, v := range L.inodeByImap {
		length++
		length += len(v)
	}
	for _, v := range L.data {
		length += len(v)
	}
	L.len = length*/
}

func (L *FSLog) log2DiskBlock(start int, inodeMap map[int]INodeMap) ([]Block, map[int]int) {
	re := []Block{}
	var segHead SegHead

	for _, v := range L.inodeByImap {
		segHead.inodeBlockN += len(v)
		segHead.inodeMapN++
	}
	segHead.inodeBlockN = segHead.inodeBlockN / INodePerBlock

	for _, v := range L.data {
		segHead.dataBlockN += len(v)
	}

	var dataBlock []Block

	for _, v := range L.inodeByImap {
		for _, n := range v {
			for _, dataB := range L.data[n.inodeN] {
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
	//var returnBlock []Block
	//returnBlock[1] = segHead
	return append(append(append([]Block{SegHead}, imapBlock...), nodesByBlock...), dataBlock...), returnMap
}
