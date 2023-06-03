package Setting

// Block size = 4KB
const BlockSize int = 4 * 1024

const INodeSize int = 512

const MaxInodeN int = 2048
const MaxINodemapPartN int = MaxInodeN/InodePerInodemapBlock + 1
const BlockN int = 1024 * 1024

const (
	BitPerBitmapBlock     int = BlockSize * 8
	INodePerBlock         int = BlockSize / INodeSize
	InodePerInodemapBlock int = BlockSize / 4
)
