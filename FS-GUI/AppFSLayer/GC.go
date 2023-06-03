package AppFSLayer

import (
	"LSF/BlockLayer"
	"LSF/DiskLayer"
	"LSF/LogLayer"
	"fmt"
	"log"
	//"fmt"
)

func (afs *AppFS) ExtractNeeded(inodeMapO map[int]BlockLayer.INodeMap, inodesO map[int]([]BlockLayer.INode), dataBsO map[int]LogLayer.DataBlockMem) (map[int]BlockLayer.INodeMap, []BlockLayer.INode, []LogLayer.DataBlockMem) {
	inodeMap := make(map[int]BlockLayer.INodeMap)
	inodes := make([]BlockLayer.INode, 0)
	dataBs := make([]LogLayer.DataBlockMem, 0)

	//superB := afs.blockFs.VD.ReadSuperBlock().(BlockLayer.SuperBlock)
	for p, inm := range inodeMapO {
		if p == afs.blockFs.GetIMapPointer(inm.Index) {
			inodeMap[inm.Index] = inm
		}
	}

	for oP, inB := range inodesO {
		for _, in := range inB {
			inN, inP := afs.blockFs.INodeN2iNodeAndPointer(in.InodeN)
			if inP == oP {
				inodes = append(inodes, inN)
			}
		}
	}

	for p, dB := range dataBs {
		if p == afs.blockFs.GetDataBPointer(dB.Inode, dB.Index) {
			dataBs = append(dataBs, dB)
		}
	}

	return inodeMap, inodes, dataBs
}

func (afs *AppFS) GC(maxSegCount int) int {
	fmt.Println("\n\nBefore:", afs.blockFs.SuperBlockUNsafe().INodeMaps)
	afs.LogCommit()
	scanStart := 0
	imapFinal := make(map[int]BlockLayer.INodeMap)
	count := 0
	for ; ; count++ {
		//fmt.Println("OK???????")
		if maxSegCount > 0 && count > maxSegCount {
			//fmt.Println("OKKKKKK")
			break
		}
		if afs.fLog.NeedCommit() {
			afs.LogCommitWithINMap(imapFinal)
			imapFinal = make(map[int]BlockLayer.INodeMap)
		}
		hP := afs.blockFs.GetOneSegHeadStartFrom(scanStart)
		if hP < 0 {
			//fmt.Println("BREAKING!")
			break
		}
		//fmt.Println("hP:", hP)

		segHead := BlockLayer.SegHead{}.FromBlock(afs.blockFs.VD.ReadBlock(hP)).(BlockLayer.SegHead)
		segLen := LogLayer.SegLenFromHead(segHead)
		segBs := []DiskLayer.RealBlock{}
		scanStart = hP + segLen
		for i := 0; i < segLen; i++ {
			segBs = append(segBs, afs.blockFs.VD.ReadBlock(hP+i))
		}
		inodeMO, inodesO, dataBsO := LogLayer.ReConstructLog(hP, segBs)
		inodeM, inodes, dataBs := afs.ExtractNeeded(inodeMO, inodesO, dataBsO)
		afs.blockFs.ReclaimBlock(hP, segLen)

		fmt.Println("inodem: ", inodeM)
		imapFinal = inmapMerge(imapFinal, inodeM)
		conSuccess := afs.fLog.ConstructLog(inodes, dataBs)
		if !conSuccess {
			//fmt.Println("AOW!!!!!!")
			afs.LogCommitWithINMap(imapFinal)
			imapFinal = make(map[int]BlockLayer.INodeMap)
			conSuccess = afs.fLog.ConstructLog(inodes, dataBs)
		}
		if !conSuccess {
			log.Fatal("Bug here!!!Reconstruction fail!")
		}
		//fmt.Println()
		fmt.Println("After:", afs.blockFs.SuperBlockUNsafe().INodeMaps)
	}
	fmt.Println("imapFinal:\n", imapFinal)
	afs.LogCommitWithINMap(imapFinal)
	fmt.Println("After:", afs.blockFs.SuperBlockUNsafe().INodeMaps)
	return count
}

func inmapMerge(a, b map[int]BlockLayer.INodeMap) map[int]BlockLayer.INodeMap {
	for i, v := range b {
		a[i] = v
	}
	return a
}

//func inmapMerge(a, b map[int]BlockLayer.INodeMap) map[int]BlockLayer.INodeMap {
