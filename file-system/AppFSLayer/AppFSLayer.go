package AppFSLayer

import (
	"LSF/BlockLayer"
	"LSF/DiskLayer"
	"LSF/LogLayer"
	"LSF/Setting"

	//"fmt"
	"log"
)

type AppFS struct {
	blockFs BlockLayer.BlockFS
	fLog    LogLayer.FSLog
}

func (afs *AppFS) FormatFS(VD DiskLayer.VirtualDisk) {
	afs.blockFs.VD = VD
	afs.fLog.InitLog()
	initINodes := []BlockLayer.INode{createInode(BlockLayer.Folder, "", true, 0)} //Adding root
	for i := 1; i < Setting.MaxInodeN; i++ {
		initINodes = append(initINodes, createInode(BlockLayer.NormalFile, "", false, i)) //Adding invalid inodes to init imap
	}
	afs.fLog.ConstructLog(initINodes, []LogLayer.DataBlockMem{})
	_, _, _, initSegLen := afs.fLog.LenInBlock()
	initStart := afs.blockFs.FindSpaceForSeg(initSegLen)
	blocks, imapLs := afs.fLog.Log2DiskBlock(initStart, make(map[int]BlockLayer.INodeMap))
	afs.blockFs.ApplyUpdate(initStart, blocks, imapLs)
	afs.fLog.InitLog()
}

func createInode(fType int, name string, valid bool, inodeN int) BlockLayer.INode {
	in := BlockLayer.INode{Valid: valid, FileType: fType, Name: name, InodeN: inodeN}
	for i, _ := range in.Pointers {
		in.Pointers[i] = -1 //Init to invalid pointers
	}
	return in
}

func (afs *AppFS) findFreeINode() int {
	for i := 0; i < Setting.MaxInodeN; i++ {
		if afs.blockFs.INodeN2iNode(i).Valid == false {
			return i
		}
	}
	return -1
}

func (afs *AppFS) LogCommit() {
	imapNeeded := make(map[int]BlockLayer.INodeMap)
	for _, v := range afs.fLog.ImapNeeded() {
		imapNeeded[v] = (afs.blockFs.VD.ReadBlock(afs.blockFs.VD.ReadSuperBlock().(BlockLayer.SuperBlock).INodeMaps[v])).(BlockLayer.INodeMap)
	} //Get inaodmap needed
	_, _, _, logSegLen := afs.fLog.LenInBlock()
	start := afs.blockFs.FindSpaceForSeg(logSegLen)
	if start < 0 {
		//WE will add GC later. TO BE DONE
		log.Fatal("No space!")
	}
	bs, newIMap := afs.fLog.Log2DiskBlock(start, imapNeeded)
	afs.blockFs.ApplyUpdate(start, bs, newIMap)
	afs.fLog.InitLog()
}

func (afs *AppFS) isINodeInLog(n int) bool {
	return afs.fLog.IsINodeInLog(n)
}

func (afs *AppFS) GetFileINfo(inodeN int) BlockLayer.INode {
	if afs.isINodeInLog(inodeN) {
		afs.LogCommit()
	}
	return afs.blockFs.INodeN2iNode(inodeN)
}

func (afs *AppFS) CreateFile(fType int, name string) int {
	newInodeN := afs.findFreeINode()
	if afs.isINodeInLog(newInodeN) {
		afs.LogCommit()
		//fmt.Println("Ha?")
		newInodeN = afs.findFreeINode()
	} //Avoid reallocating a inode.

	if newInodeN == -1 {
		log.Fatal("No inode number available.") //Maybe later we should check the log? Maybe later. TO BE DONE
	}
	if afs.fLog.ConstructLog([]BlockLayer.INode{createInode(fType, name, true, newInodeN)}, []LogLayer.DataBlockMem{}) {
	} else {
		afs.LogCommit()
		//fmt.Println("Oh?")
		afs.fLog.ConstructLog([]BlockLayer.INode{createInode(fType, name, true, newInodeN)}, []LogLayer.DataBlockMem{})
	}
	return newInodeN
}

func (afs *AppFS) WriteFile(inodeN int, index []int, data []DiskLayer.Block) {
	if afs.isINodeInLog(inodeN) {
		afs.LogCommit()
	}
	inode := afs.blockFs.INodeN2iNode(inodeN)
	if inode.Valid == false {
		log.Fatal("Invalid write to non-exsistent inode:", inodeN, "  get inode:", inode)
	}
	ds := []LogLayer.DataBlockMem{}
	for i, v := range index {
		ds = append(ds, LogLayer.DataBlockMem{Inode: inodeN, Index: v, Data: data[i]})
	}
	afs.fLog.ConstructLog([]BlockLayer.INode{inode}, ds)
}

func (afs *AppFS) ReadFile(inodeN int, index int) DiskLayer.Block {
	if afs.isINodeInLog(inodeN) {
		afs.LogCommit()
	}
	return afs.blockFs.ReadFile(inodeN, index)
}

func (afs *AppFS) DeleteFile(inodeN int) {
	if afs.isINodeInLog(inodeN) {
		afs.LogCommit()
	}
	inode := BlockLayer.INode{InodeN: inodeN, Valid: false}
	afs.fLog.ConstructLog([]BlockLayer.INode{inode}, []LogLayer.DataBlockMem{})
}

func (afs *AppFS) DeleteBlockInFile(inodeN int, index []int) {
	if afs.isINodeInLog(inodeN) {
		afs.LogCommit()
	}
	inode := afs.blockFs.INodeN2iNode(inodeN)
	for _, v := range index {
		inode.Pointers[v] = -1
	}
	afs.fLog.ConstructLog([]BlockLayer.INode{inode}, []LogLayer.DataBlockMem{})
}

//////////////////////
/////////////    the function bellow is to get debug info. Don't use these!

func (afs *AppFS) ReadBlockUnsafe(a int) DiskLayer.Block {
	return afs.blockFs.VD.ReadBlock(a)
}
func (afs *AppFS) ReadSuperUnsafe() BlockLayer.SuperBlock {
	return afs.blockFs.VD.ReadSuperBlock().(BlockLayer.SuperBlock)
}

/*func (afs *AppFS) PrintLogUnsafe() {
	afs.fLog.PrintLog()
}*/
