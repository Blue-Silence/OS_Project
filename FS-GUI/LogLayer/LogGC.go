package LogLayer

import (
	"LSF/BlockLayer"
	"LSF/DiskLayer"

	//"LSF/Setting"
	//"fmt"
	"log"
)

func SegLenFromHead(s BlockLayer.SegHead) int {
	return 1 + s.InodeMapN + s.InodeBlockN + s.DataBlockN
}

func ReConstructLog(start int, segB []DiskLayer.RealBlock) (map[int]BlockLayer.INodeMap, map[int]([]BlockLayer.INode), map[int]DataBlockMem) {
	inodeMap := make(map[int]BlockLayer.INodeMap)
	inodes := make(map[int]([]BlockLayer.INode))
	dataBs := make(map[int]DataBlockMem)
	head := BlockLayer.SegHead{}.FromBlock(segB[0]).(BlockLayer.SegHead)
	if SegLenFromHead(head) != len(segB) {
		log.Fatal("Seg len mismatch!")
	}
	i := 1
	for c := 0; c < head.InodeMapN; c++ {
		b := BlockLayer.INodeMap{}.FromBlock(segB[i]).(BlockLayer.INodeMap)
		inodeMap[i+start] = b
		i++
	}
	for c := 0; c < head.InodeBlockN; c++ {
		b := BlockLayer.INodeBlock{}.FromBlock(segB[i]).(BlockLayer.INodeBlock)
		for _, v := range b.NodeArr {
			if v.Valid {
				inodes[i+start] = append(inodes[i+start], v)
			}
		}
		i++
	}

	for c := 0; c < head.DataBlockN; c++ {
		b := segB[i]
		blockInfo := head.DataBlockSummary[c]
		dataBs[i+start] = DataBlockMem{Inode: blockInfo.InodeN, Index: blockInfo.Index, Data: b}
	}

	return inodeMap, inodes, dataBs

}
