package main

import (
	"LSF/AppFSLayer"
	"LSF/BlockLayer"
	"LSF/Setting"
	"fmt"
	"math/rand"
)

func main() {
	test1()

}

func printBlock(afs AppFSLayer.AppFS, length int) {
	for i := 0; i < length; i++ {
		printSingleBlock(afs, i)
	}
}

func createDataBlock() BlockLayer.DataBlock {
	var dataB BlockLayer.DataBlock
	for i := 0; i < Setting.BlockSize; i++ {
		dataB.Data[i] = uint8(rand.Uint32())
	}
	return dataB
}

func printSingleBlock(afs AppFSLayer.AppFS, index int) {
	fmt.Print("\n\n", index, "th block:\n")
	switch afs.ReadBlockUnsafe(index).(type) {
	case BlockLayer.INodeBlock:
		fmt.Print("This is a inode block\n")
	case BlockLayer.SegHead:
		fmt.Print("This is a seg header block\n")
	case BlockLayer.INodeMap:
		fmt.Print("This is a inode map block\n")
	case BlockLayer.DataBlock:
		fmt.Print("This is a data block\n")
	default:
		fmt.Print("What is this?\n")
	}
	fmt.Print(afs.ReadBlockUnsafe(index), "\n")
}

func printSuperBlock(afs AppFSLayer.AppFS) {
	fmt.Println("Now the super!\n")
	//fmt.Println("Now the super!")
	fmt.Println("Now the inodeMap location:\n", afs.ReadSuperUnsafe())
	//fmt.Println("Now the bitmap location:\n", afs.fs.VD.superBlock.bitMap)
}
