package main

import "log"

type AppFS struct {
	fs   FS
	fLog FSLog
}

func (afs *AppFS) formatFS(VD VirtualDisk) {
	afs.fs.VD = VD
	afs.fLog.initLog()
	initINodes := []INode{createInode(Folder, "", true, 0)} //Adding root
	for i := 1; i < MaxInodeN; i++ {
		initINodes = append(initINodes, createInode(NormalFile, "", false, i)) //Adding invalid inodes to init imap
	}
	afs.fLog.constructLog(initINodes, []DataBlockMem{})
	_, _, _, initSegLen := afs.fLog.lenInBlock()
	initStart := afs.fs.findSpaceForSeg(initSegLen)
	blocks, imapLs := afs.fLog.log2DiskBlock(initStart, make(map[int]INodeMap))
	afs.fs.applyUpdate(initStart, blocks, imapLs)
}

func createInode(fType int, name string, valid bool, inodeN int) INode {
	return INode{valid: valid, fileType: fType, name: name, inodeN: inodeN}
}

func (afs *AppFS) findFreeINode() int {
	for i := 0; i < MaxInodeN; i++ {
		if afs.fs.iNodeN2iNode(i).valid == false {
			return i
		}
	}
	return -1
}

func (afs *AppFS) logCommit() {
	imapNeeded := make(map[int]INodeMap)
	for _, v := range afs.fLog.imapNeeded() {
		imapNeeded[v] = (afs.fs.VD.readBlock(afs.fs.superBlock.iNodeMaps[v])).(INodeMap)
	} //Get inaodmap needed
	_, _, _, logSegLen := afs.fLog.lenInBlock()
	start := afs.fs.findSpaceForSeg(logSegLen)
	if start < 0 {
		//WE will add GC later. TO BE DONE
		log.Fatal("No space!")
	}
	bs, newIMap := afs.fLog.log2DiskBlock(start, imapNeeded)
	afs.fs.applyUpdate(start, bs, newIMap)
	afs.fLog.initLog()
}

func (afs *AppFS) isINodeInLog(n int) bool {
	for _, v := range afs.fLog.inodeByImap[n/InodePerInodemapBlock] {
		if v.inodeN == n {
			return true
		}
	}
	return false
}

func (afs *AppFS) createFile(fType int, name string) int {
	newInodeN := afs.findFreeINode()
	if afs.isINodeInLog(newInodeN) {
		afs.logCommit()
		newInodeN = afs.findFreeINode()
	} //Avoid reallocating a inode.

	if newInodeN == -1 {
		log.Fatal("No inode number available.") //Maybe later we should check the log? Maybe later. TO BE DONE
	}
	if afs.fLog.constructLog([]INode{createInode(fType, name, true, newInodeN)}, []DataBlockMem{}) {
	} else {
		afs.logCommit()
		afs.fLog.constructLog([]INode{createInode(fType, name, true, newInodeN)}, []DataBlockMem{})
	}
	return newInodeN
}
