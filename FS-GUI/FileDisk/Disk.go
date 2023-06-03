package FileDisk

import (
	"LSF/DiskLayer"
	"LSF/Setting"
	"log"
	"os"
	"runtime/debug"
)

type FileDisk struct {
	//blocks     [Setting.BlockN]Block
	blocksFile *os.File // [Setting.BlockN]DiskLayer.RealBlock
	//superBlock Block
	//superBlock RealBlock
	superBlockFile *os.File
}

func (d *FileDisk) Init(blocksFilePath, superBlockFilePath string) {
	var err1 error
	var err2 error
	d.blocksFile, err1 = os.OpenFile(blocksFilePath, os.O_CREATE|os.O_RDWR, 0666)
	d.superBlockFile, err2 = os.OpenFile(superBlockFilePath, os.O_CREATE|os.O_RDWR, 0666)
	if err1 != nil || err2 != nil {
		log.Fatal("Error when opening the file. Err1:", err1, "  Err2:", err2)
	}
}

func (d *FileDisk) ReadBlock(index int) *DiskLayer.RealBlock {
	if index < 0 || index > Setting.BlockN {
		debug.PrintStack()
		log.Fatal("Invalid disk read access at ", index)
	}
	var b DiskLayer.RealBlock
	d.blocksFile.ReadAt(b[:], int64(index)*int64(Setting.BlockSize))
	//log.Println("Block.Read at", index, " n:", n, "  err:", err)
	return &b
}

func (d *FileDisk) WriteBlock(index int, b DiskLayer.Block) {
	if index < 0 || index > Setting.BlockN {
		debug.PrintStack()
		log.Fatal("Invalid disk write access at ", index)
	}
	//d.blocks[index] = b.ToBlock()
	data := b.ToBlock()
	d.blocksFile.WriteAt(data[:], int64(index)*int64(Setting.BlockSize))
	//log.Println("Block.Write at", index, " n:", n, "  err:", err)
}

func (d *FileDisk) ReadSuperBlock() []*DiskLayer.RealBlock {
	data := make([]byte, 4096*Setting.BlockSize)
	d.superBlockFile.ReadAt(data, 0)
	//log.Println("Superblock.Read n:", n, "  err:", err)

	return (DiskLayer.BytesToBlocks(data))
}

func (d *FileDisk) WriteSuperBlock(b []*DiskLayer.RealBlock) {
	bs := []byte{}
	for _, v := range b {
		bs = append(bs, DiskLayer.BlockToBytes(v)...)
	}
	_, _ = d.superBlockFile.WriteAt(bs, 0)
	//log.Println("Superblock.Write n:", n, "  err:", err)
}

func (d *FileDisk) Close() {
	d.blocksFile.Close()
	d.superBlockFile.Close()
	log.Println("Closing!")
}
