package LogLayer

import (
	"LSF/BlockLayer"
	"LSF/DiskLayer"
	"LSF/Setting"
	"fmt"

	//"fmt"
	"log"
)

const LogLenth = 200

type DataBlockMem struct {
	Inode int
	Index int
	Data  DiskLayer.RealBlock
}

type FSLog struct {
	data        map[int]([]DataBlockMem)
	inodeByImap map[int]([]BlockLayer.INode)
	//len         int
} // This is in memory

func (L *FSLog) InitLog() {
	L.data = make(map[int][]DataBlockMem)
	L.inodeByImap = make(map[int][]BlockLayer.INode)
	//L.len = 0
}

func (L *FSLog) ConstructLog(inodes []BlockLayer.INode, ds []DataBlockMem) bool {
	_, _, dataBlockN, _ := L.LenInBlock()
	if len(ds)+dataBlockN > BlockLayer.MaxEditBlcokN {
		return false
	}
	//fmt.Println("Test:", len(ds)+dataBlockN)
	for _, v := range inodes {
		L.inodeByImap[v.InodeN/Setting.InodePerInodemapBlock] = append(L.inodeByImap[v.InodeN/Setting.InodePerInodemapBlock], v)
	}
	for _, v := range ds {
		L.data[v.Inode] = append(L.data[v.Inode], v)
		//fmt.Printf("inode:%v index:%v\n", v.Inode, v.Index)
	}
	//_, _, dataBlockN2, _ := L.LenInBlock()
	//fmt.Println("AAA DataBlockN:", dataBlockN, "   new:", dataBlockN2)
	return true
}

func (L *FSLog) LenInBlock() (int, int, int, int) {
	inodeBlockN := 0
	inodeMapN := 0
	dataBlockN := 0
	segLen := 0
	for _, v := range L.inodeByImap {
		inodeBlockN += len(v)
		inodeMapN++
	}

	if inodeBlockN%Setting.INodePerBlock != 0 {
		inodeBlockN = inodeBlockN/Setting.INodePerBlock + 1
	} else {
		inodeBlockN = inodeBlockN / Setting.INodePerBlock
	}

	for _, v := range L.data {
		dataBlockN += len(v)
	}

	segLen = inodeMapN + inodeBlockN + dataBlockN + 1
	return inodeMapN, inodeBlockN, dataBlockN, segLen
}

func (L *FSLog) Log2DiskBlock(start int, inodeMap map[int]BlockLayer.INodeMap) ([]DiskLayer.Block, map[int]int) {
	var segHead BlockLayer.SegHead
	segHead.InodeMapN, segHead.InodeBlockN, segHead.DataBlockN, _ = L.LenInBlock()

	inodeMapC := make(map[int]BlockLayer.INodeMap)
	for i, v := range inodeMap {
		inodeMapC[i] = v
	}
	for _, v := range L.inodeByImap {
		for _, n := range v {
			inodeMapC[n.InodeN/Setting.InodePerInodemapBlock] = BlockLayer.INodeMap{}
		}
	}
	segHead.InodeMapN = 0
	for _, _ = range inodeMapC {
		segHead.InodeMapN++
	} //This way we will get the real length of InodeMap block

	var dataBlock []DiskLayer.Block

	//fmt.Println("segHead.DataBlockN:", segHead.DataBlockN)

	for iv, v := range L.inodeByImap {
		for in, n := range v {
			//fmt.Println("\nInode:", n.InodeN, "    dataBlock num:", len(L.data[n.InodeN]))

			for _, dataB := range L.data[n.InodeN] {
				dataBlock = append(dataBlock, BlockLayer.DataBlock{dataB.Data})
				segHead.DataBlockSummary[len(dataBlock)-1] = BlockLayer.FileIndexInfo{n.InodeN, dataB.Index}
				n.Pointers[dataB.Index] = start + 1 + segHead.InodeMapN + segHead.InodeBlockN + len(dataBlock) - 1
				L.inodeByImap[iv][in] = n
			}
		}
	}

	var nodesByBlock []BlockLayer.INodeBlock

	nodeCount := Setting.INodePerBlock
	for _, v := range L.inodeByImap {
		for _, n := range v {
			if nodeCount == Setting.INodePerBlock {
				nodesByBlock = append(nodesByBlock, BlockLayer.INodeBlock{})
				nodeCount = 0
			}
			nodesByBlock[len(nodesByBlock)-1].NodeArr[nodeCount] = n
			nodeCount++
			//and also do something to change imap next (TO BE DONE)
			iPart := inodeMap[n.InodeN/Setting.InodePerInodemapBlock]
			iPart.Index = n.InodeN / Setting.InodePerInodemapBlock
			(iPart.InodeMapPart)[n.InodeN%Setting.InodePerInodemapBlock] = len(nodesByBlock) - 1 + start + 1 + segHead.InodeMapN
			inodeMap[n.InodeN/Setting.InodePerInodemapBlock] = iPart
		}
	}
	//_, InodeBlockN233, _, _ := L.LenInBlock()
	if len(nodesByBlock) != segHead.InodeBlockN {

		log.Println(len(nodesByBlock))
		log.Println(segHead.InodeBlockN)
		//log.Println(InodeBlockN233)
		log.Println(L.inodeByImap)
		log.Fatal("Mismatch!")
	}

	var imapBlock []BlockLayer.INodeMap
	returnMap := make(map[int]int)
	for _, v := range inodeMap {
		imapBlock = append(imapBlock, v)
		returnMap[v.Index] = start + 1 + len(imapBlock) - 1
	}

	re := []DiskLayer.Block{segHead}
	for _, v := range imapBlock {
		re = append(re, v)
	}
	for _, v := range nodesByBlock {
		re = append(re, v)
	}
	for _, v := range dataBlock {
		re = append(re, v)
	}
	return re, returnMap
}

func (L *FSLog) ImapNeeded() []int {
	var re []int
	for i, _ := range L.inodeByImap {
		re = append(re, i)
	}
	return re
}

func (L *FSLog) IsINodeInLog(n int) bool {
	for _, v := range L.inodeByImap[n/Setting.InodePerInodemapBlock] {
		if v.InodeN == n {
			return true
		}
	}
	return false
}

func (L *FSLog) NeedCommit() bool {
	_, _, dataBlockN, segLen := L.LenInBlock()
	return segLen > LogLenth/3*2 || dataBlockN > BlockLayer.MaxEditBlcokN/3*2
}

// ///////////////////////////////////////////////////////////////////////////
// ///// FOR TEST

func (L *FSLog) PrintLog() {
	//fmt.Print(L.Log2DiskBlock(0, make(map[int]BlockLayer.INodeMap)))
	fmt.Println(L.inodeByImap)
	for _, v := range L.data {
		for _, v1 := range v {
			fmt.Println("INode:", v1.Inode, "   Index:", v1.Index, "  Data:...")
		}
	}

}
