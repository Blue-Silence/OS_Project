package Folder

import (
	"LSF/AppFSLayer"
	"LSF/DiskLayer"
	"LSF/Setting"
)

const (
	MaxFilePerFolderBlock = (Setting.BlockSize / 512) //this can be changed.
	MaxNameLen            = 256
)

type FolderBlock struct {
	fileEntrys [MaxFilePerFolderBlock]fileEntry
} // 1 per block
type fileEntry struct {
	name  string
	inode int
	valid bool
}

func (s FolderBlock) CanBeBlock() {
}

func concatFolder(afs *AppFSLayer.AppFS, folderIN int) []fileEntry {
	folderINode := afs.GetFileINfo(folderIN)
	re := []fileEntry{}
	for i, _ := range folderINode.Pointers {
		co := afs.ReadFile(folderIN, i).(FolderBlock).fileEntrys
		re = append(re, co[:]...)
	}
	return re
}

func rebuildFolder(fEs []fileEntry) ([]int, []DiskLayer.Block) {
	returnIndex := []int{}
	returnBlock := []DiskLayer.Block{}
	i := 0
	for {
		if len(fEs) > 0 {
			returnIndex = append(returnIndex, i)
			i++
			fEB := FolderBlock{}
			copy(fEB.fileEntrys[:], fEs)
			returnBlock = append(returnBlock, fEB)

			if len(fEs) > MaxFilePerFolderBlock {
				fEs = fEs[MaxFilePerFolderBlock-1:]
			} else {
				break
			}
		} else {
			break
		}
	}
	return returnIndex, returnBlock
}

func AddFileToFolder(afs *AppFSLayer.AppFS, folderIN int, fileIN int) {
	fEs := concatFolder(afs, folderIN)
	fileINode := afs.GetFileINfo(fileIN)
	folderINode := afs.GetFileINfo(folderIN)
	fE := fileEntry{name: fileINode.Name, inode: fileINode.InodeN, valid: true}
	fEs = append(fEs, fE)
	indexs, bs := rebuildFolder(fEs)
	afs.DeleteBlockInFile(folderIN, folderINode.Pointers[:])
	afs.WriteFile(folderIN, indexs, bs)
}

func DeleteFileToFolder(afs *AppFSLayer.AppFS, folderIN int, fileIN int) {
	fEs := concatFolder(afs, folderIN)
	folderINode := afs.GetFileINfo(folderIN)
	newFe := []fileEntry{}
	for _, v := range fEs {
		if v.inode != fileIN {
			newFe = append(newFe, v)
		}
	}
	indexs, bs := rebuildFolder(newFe)
	afs.DeleteBlockInFile(folderIN, folderINode.Pointers[:])
	afs.WriteFile(folderIN, indexs, bs)
	afs.DeleteFile(fileIN)
}

func GetFolderContent(afs *AppFSLayer.AppFS, inode int) map[string]int {
	nameMapping := make(map[string]int)
	for _, v := range concatFolder(afs, inode) {
		if v.valid {
			nameMapping[v.name] = v.inode
		}
	}
	return nameMapping
}
