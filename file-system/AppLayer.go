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

func (afs *AppFS) createFile(fType int, name string) int {
	newInodeN := afs.findFreeINode()
	//MAYBE WE SHOULD CLEAR THE LOG BEFORE CONTINUE?
	if newInodeN == -1 {
		log.Fatal("No inode number available.") //Maybe later we should check the log? Maybe later. TO BE DONE
	}
	afs.fLog.constructLog([]INode{createInode(fType, name, true, newInodeN)}, []DataBlockMem{})
}

func (afs *AppFS) findFreeINode() int {
	for i := 0; i < MaxInodeN; i++ {
		if afs.fs.iNodeN2iNode(i).valid == false {
			return i
		}
	}
	return -1
}
