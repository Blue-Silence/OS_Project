package BlockLayer

import "LSF/Setting"

const MaxEditBlcokN int = 100
const DirectPointerPerINode = 10

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
	Index        int
	InodeMapPart [Setting.InodePerInodemapBlock]int
} // 1 per block

type SegHead struct {
	//SegLen           int // total len in block size.
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
	Pointers [DirectPointerPerINode]int

	PointerToNextINode int
	CurrentLevel       int
	IsRoot             bool
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
