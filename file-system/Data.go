package main

import "fmt"

const MaxEditBlcokN int = 100

const (
	NormalFile int = 0
	Folder     int = 1
)

//const MaxINodeMapBlockN int = 10

type SuperBlock struct {
	//newestSegHead int
	iNodeMaps [MaxINodemapPartN]int
	bitMap    [BlockN]bool
} //store at separate location.

type INodeMap struct {
	offset       int
	inodeMapPart [InodePerInodemapBlock]int
} // 1 per block

type SegHead struct {
	segLen           int // total len in block size.
	inodeMapN        int
	inodeBlockN      int
	dataBlockN       int
	dataBlockSummary [MaxEditBlcokN]FileIndexInfo
} // 1 per block

func (s SegHead) len() int {
	return s.inodeBlockN + s.inodeMapN + s.dataBlockN + 1
}

type FileIndexInfo struct {
	inodeN int
	index  int
}

type INode struct {
	valid    bool
	inodeN   int
	name     string
	fileType int
	pointers [10]int
}

type INodeBlock struct {
	nodeArr [INodePerBlock]INode
} // 1 per block

type DataBlock struct {
	data [BlockSize]uint8
} // 1 per block

type FolderBlock struct {
	names   [BlockSize / 256]string
	inodeNs [BlockSize / 256]int
} // 1 per block

type Block interface {
	toStore() string
	fromStore(string) Block
}

func (s SuperBlock) toStore() string {
	return fmt.Sprintf("{SUPERBLOCK %v %v %v}", s.newestSegHead, s.iNodeMaps, s.bitMap)
}

func (s SuperBlock) fromStore(str string) Block {
	fmt.Sscanf(str, "{SUPERBLOCK %v %v %v}", &s.iNodeMaps, &s.bitMap)
	return s
}

func (s SegHead) toStore() string {
	return fmt.Sprintf("{Seg %v %v %v %v}", s.dataBlockN, s.dataBlockSummary, s.inodeMapStart, s.segLen)
}

func (s SegHead) fromStore(str string) Block {
	fmt.Sscanf(str, "{Seg %v %v %v %v}", &s.dataBlockN, &s.dataBlockSummary, &s.inodeMapStart, &s.segLen)
	return s
}

func (s INodeMap) toStore() string {
	return fmt.Sprintf("{NMAP %v}", s.inodeMapPart)
}

func (s INodeMap) fromStore(str string) Block {
	fmt.Sscanf(str, "{NMAP %v}", &s.inodeMapPart)
	return s
}

func (s INodeBlock) toStore() string {
	return fmt.Sprintf("{IN %v}", s.nodeArr)
}

func (s INodeBlock) fromStore(str string) Block {
	fmt.Sscanf(str, "{IN %v}", &s.nodeArr)
	return s
}

func (s DataBlock) toStore() string {
	return fmt.Sprintf("{DB %v}", s.data)
}

func (s DataBlock) fromStore(str string) Block {
	fmt.Sscanf(str, "{DB %v}", &s.data)
	return s
}
