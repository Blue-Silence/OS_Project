package BlockLayer

import "LSF/Setting"

const MaxEditBlcokN int = 100

const (
	NormalFile int = 0
	Folder     int = 1
)

//const MaxINodeMapBlockN int = 10

type SuperBlock struct {
	//newestSegHead int
	INodeMaps [Setting.MaxINodemapPartN]int
	BitMap    [Setting.BlockN]bool
} //store at separate location.

type INodeMap struct {
	Offset       int
	InodeMapPart [Setting.InodePerInodemapBlock]int
} // 1 per block

type SegHead struct {
	SegLen           int // total len in block size.
	InodeMapN        int
	InodeBlockN      int
	DataBlockN       int
	DataBlockSummary [MaxEditBlcokN]FileIndexInfo
} // 1 per block

func (s SegHead) Len() int {
	return s.InodeBlockN + s.InodeMapN + s.DataBlockN + 1
}

type FileIndexInfo struct {
	InodeN int
	Index  int
}

type INode struct {
	Valid    bool
	InodeN   int
	Name     string
	FileType int
	Pointers [10]int
}

type INodeBlock struct {
	NodeArr [Setting.INodePerBlock]INode
} // 1 per block

type DataBlock struct {
	Data [Setting.BlockSize]uint8
} // 1 per block

func (s SuperBlock) CanBeBlock() {
}
func (s SegHead) CanBeBlock() {
}
func (s INodeMap) CanBeBlock() {
}
func (s INodeBlock) CanBeBlock() {
}
func (s DataBlock) CanBeBlock() {
}

/*
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
}*/
