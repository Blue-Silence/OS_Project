package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("Test FS format:")
	testFS := AppFS{}
	testFS.formatFS(VirtualDisk{})
	printSuperBlock(testFS)
	printSingleBlock(testFS, 0)
	printSingleBlock(testFS, 1)
	//printBlock(testFS, 100)
	node1N := testFS.createFile(NormalFile, "test1")

	testFS.logCommit()
	//printBlock(testFS, 100)
	fmt.Println("\n\n\n\n\nAnd the new superBlock\n")
	printSuperBlock(testFS)
	printSingleBlock(testFS, 130)
	printSingleBlock(testFS, 131)
	printSingleBlock(testFS, 132)
	printSingleBlock(testFS, 133)
	fmt.Println("My dear world!")

	testFS.deleteFile(node1N)
	testFS.logCommit()
	printSuperBlock(testFS)
	printSingleBlock(testFS, 133)
	printSingleBlock(testFS, 134)
	printSingleBlock(testFS, 135)
	printSingleBlock(testFS, 136)
	printSingleBlock(testFS, 2)

	node2N := testFS.createFile(NormalFile, "test2")
	//testFS.logCommit()
	node3N := testFS.createFile(NormalFile, "test3")
	testFS.logCommit()
	testFileWrite(testFS)
	fmt.Println(node2N)
	fmt.Println(node3N)

}

func printBlock(afs AppFS, length int) {
	for i := 0; i < length; i++ {
		printSingleBlock(afs, i)
	}
}

func createDataBlock() DataBlock {
	var dataB DataBlock
	for i := 0; i < BlockSize; i++ {
		dataB.data[i] = uint8(rand.Uint32())
	}
	return dataB
}

func testFileWrite(afs AppFS) {
	indexL := []int{1, 2, 3, 6, 9}
	fmt.Println("-------------------Test write start----------------------------")
	inodeN := afs.createFile(NormalFile, "TestWrite")
	afs.logCommit()
	ds := []Block{createDataBlock(), createDataBlock(), createDataBlock(), createDataBlock(), createDataBlock()}
	afs.writeFile(inodeN, indexL, ds)
	fmt.Println("-------------------Log after write-----------------------")
	fmt.Println(afs.fLog.inodeByImap)
	afs.logCommit()
	inode := afs.getFileINfo(inodeN)
	fmt.Println("-------------------Inode after write-----------------------")
	fmt.Println(inode)
	fmt.Println("-------------------Data check-----------------------")
	printSingleBlock(afs, 145)
	printSingleBlock(afs, 146)
	printSingleBlock(afs, 147)
	printSingleBlock(afs, 148)
	for iD, v := range indexL {
		testF := true
		arr := afs.readFile(inodeN, v).(DataBlock).data
		for i := 0; i < BlockSize; i++ {
			if ds[iD].(DataBlock).data[i] != arr[i] {
				testF = false
			}
		}
		fmt.Println("Index:", v, "  test passed?: ", testF)
	}
	fmt.Println("-------------------Test write done----------------------------")
}

func printSingleBlock(afs AppFS, index int) {
	fmt.Print("\n\n", index, "th block:\n")
	switch afs.fs.VD.blocks[index].(type) {
	case INodeBlock:
		fmt.Print("This is a inode block\n")
	case SegHead:
		fmt.Print("This is a seg header block\n")
	case INodeMap:
		fmt.Print("This is a inode map block\n")
	case DataBlock:
		fmt.Print("This is a data block\n")
	default:
		fmt.Print("What is this?\n")
	}
	fmt.Print(afs.fs.VD.blocks[index], "\n")
}

func printSuperBlock(afs AppFS) {
	fmt.Println("Now the super!\n")
	//fmt.Println("Now the super!")
	fmt.Println("Now the inodeMap location:\n", afs.fs.VD.superBlock.iNodeMaps)
	//fmt.Println("Now the bitmap location:\n", afs.fs.VD.superBlock.bitMap)
}
