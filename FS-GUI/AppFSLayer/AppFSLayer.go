package AppFSLayer

import (
	"LSF/BlockLayer"
	"LSF/DiskLayer"
	"LSF/LogLayer"
	"LSF/Setting"

	//"fmt"

	//"fmt"
	"log"
)

type AppFS struct {
	blockFs BlockLayer.BlockFS
	fLog    LogLayer.FSLog
}

func (afs *AppFS) FormatFS(VD DiskLayer.VirtualDisk) {
	//afs.blockFs.VD = VD
	afs.fLog.InitLog()
	//var superInit BlockLayer.SuperBlock
	afs.blockFs.Init(VD)

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
	in.IsRoot = true
	in.PointerToNextINode = -1
	in.CurrentLevel = 0
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

func (afs *AppFS) LogCommitWithINMap(imapNeeded map[int]BlockLayer.INodeMap) {
	for _, v := range afs.fLog.ImapNeeded() {
		imapNeeded[v] = BlockLayer.INodeMap{}.FromBlock(afs.blockFs.VD.ReadBlock(BlockLayer.SuperBlock{}.FromBlocks(afs.blockFs.VD.ReadSuperBlock()).INodeMaps[v])).(BlockLayer.INodeMap)
	} //Get inaodmap needed
	c := 0
	for range imapNeeded {
		c++
	}
	_, _, _, logSegLen := afs.fLog.LenInBlock()
	start := afs.blockFs.FindSpaceForSeg(logSegLen + c)
	if start < 0 {
		//WE will add GC later. TO BE DONE
		afs.GC(-1)
		log.Fatal("No space!")
	}
	bs, newIMap := afs.fLog.Log2DiskBlock(start, imapNeeded)
	afs.blockFs.ApplyUpdate(start, bs, newIMap)
	afs.fLog.InitLog()
}

func (afs *AppFS) LogCommit() {
	afs.LogCommitWithINMap(make(map[int]BlockLayer.INodeMap))
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
	/*if afs.fLog.ConstructLog([]BlockLayer.INode{createInode(fType, name, true, newInodeN)}, []LogLayer.DataBlockMem{}) {
	} else {
		afs.LogCommit()
		//fmt.Println("Oh?")
		afs.fLog.ConstructLog([]BlockLayer.INode{createInode(fType, name, true, newInodeN)}, []LogLayer.DataBlockMem{})
	}*/
	afs.tryLog([]BlockLayer.INode{createInode(fType, name, true, newInodeN)}, []LogLayer.DataBlockMem{})
	return newInodeN
}

func (afs *AppFS) WriteFile(inodeN int, index []int, data []DiskLayer.Block) {
	inode := afs.GetFileINfo(inodeN)
	if inode.Valid == false {
		log.Fatal("Invalid write to non-exsistent inode:", inodeN, "  get inode:", inode)
	}
	for i, ind := range index {
		_, _, traces := afs.findBlockFromStart(true, inodeN, ind)
		inode := afs.GetFileINfo(traces[len(traces)-1].inode.InodeN)
		afs.tryLog([]BlockLayer.INode{inode}, []LogLayer.DataBlockMem{LogLayer.DataBlockMem{Inode: inode.InodeN, Index: traces[len(traces)-1].offset, Data: data[i].ToBlock()}})
	}
}

func (afs *AppFS) ReadFile(inodeN int, index int) DiskLayer.RealBlock {
	b, _, traces := afs.findBlockFromStart(false, inodeN, index)
	if b && traces[len(traces)-1].offset >= 0 {
		inode := afs.GetFileINfo(traces[len(traces)-1].inode.InodeN)
		//fmt.Println("trace:", traces)
		return afs.blockFs.ReadFile(inode.InodeN, traces[len(traces)-1].offset)
	} else {
		var e DiskLayer.RealBlock
		return e
	}
}

func (afs *AppFS) DeleteFile(inodeN int) {
	if afs.isINodeInLog(inodeN) {
		afs.LogCommit()
	}
	inode := BlockLayer.INode{InodeN: inodeN, Valid: false}
	//afs.fLog.ConstructLog([]BlockLayer.INode{inode}, []LogLayer.DataBlockMem{})
	afs.tryLog([]BlockLayer.INode{inode}, []LogLayer.DataBlockMem{})
}

func (afs *AppFS) DeleteBlockInFile(inodeN int, index []int) {
	inode := afs.GetFileINfo(inodeN)
	for _, ind := range index {
		b, _, traces := afs.findBlockFromStart(false, inodeN, ind)
		if b {
			inode = afs.GetFileINfo(traces[len(traces)-1].inode.InodeN)
			inode.Pointers[traces[len(traces)-1].offset] = -1
			afs.tryLog([]BlockLayer.INode{inode}, []LogLayer.DataBlockMem{})
		}
	}
}

func (afs *AppFS) tryLog(inodes []BlockLayer.INode, ds []LogLayer.DataBlockMem) {
	if afs.fLog.NeedCommit() {
		afs.LogCommit()
	}
	if !afs.fLog.ConstructLog(inodes, ds) {
		afs.LogCommit()
		if !afs.fLog.ConstructLog(inodes, ds) {
			log.Fatal("No space!")
		}
	}
}

//////////////////////
/////////////    the function bellow is to get debug info. Don't use these!

func (afs *AppFS) ReadBlockUnsafe(a int) DiskLayer.Block {
	//return afs.blockFs.VD.ReadBlock(a)
	return nil
}
func (afs *AppFS) ReadSuperUnsafe() BlockLayer.SuperBlock {
	return BlockLayer.SuperBlock{}.FromBlocks(afs.blockFs.VD.ReadSuperBlock())
}

func (afs *AppFS) ReadInodeUnsafe(n int) BlockLayer.INode {
	return afs.GetFileINfo(n)
}

/*func (afs *AppFS) PrintLogUnsafe() {
	afs.fLog.PrintLog()
}*/
