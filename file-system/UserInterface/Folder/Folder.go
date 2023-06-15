package Folder

import (
	"LSF/AppFSLayer"
	"LSF/DiskLayer"
	"LSF/Setting"
	//"fmt"
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
	//maxI := 0
	for i, v := range folderINode.Pointers {
		if v >= 0 {
			//maxI = i
			////fmt.Println("Is it ok?")
			co := afs.ReadFile(folderIN, i).(FolderBlock).fileEntrys
			////fmt.Println("Ok!")
			for _, v := range co {
				if v.valid {
					re = append(re, v)
				}
			}
			//re = append(re, co[:]...)
		}
	}
	//fmt.Println("Max i to:", maxI)
	return re
}

func rebuildFolder(fEsO []fileEntry) ([]int, []DiskLayer.Block) {
	returnIndex := []int{}
	returnBlock := []DiskLayer.Block{}
	fEs := []fileEntry{}
	for _, v := range fEsO {
		if v.valid {
			fEs = append(fEs, v)
		}
	}
	i := 0
	for {
		if len(fEs) > 0 {
			returnIndex = append(returnIndex, i)
			i++
			fEB := FolderBlock{}
			copy(fEB.fileEntrys[:], fEs)
			returnBlock = append(returnBlock, fEB)

			if len(fEs) > MaxFilePerFolderBlock {
				fEs = fEs[MaxFilePerFolderBlock:]
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
	//fmt.Println("Before adding:", fEs)
	//fEs
	fEs = append(fEs, fE)

	indexs, bs := rebuildFolder(fEs)

	deleteI := []int{}
	for i, _ := range folderINode.Pointers {
		deleteI = append(deleteI, i)
	}
	afs.DeleteBlockInFile(folderIN, deleteI)

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

	deleteI := []int{}
	for i, _ := range folderINode.Pointers {
		deleteI = append(deleteI, i)
	}

	afs.DeleteBlockInFile(folderIN, deleteI)

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

// ////////////////////////////////////////////////////////////////////////////////
// / FOR TEST
type FolderBlock2 struct {
	F [MaxFilePerFolderBlock]FileEntry2
} // 1 per block
type FileEntry2 struct {
	Name  string
	Inode int
	Valid bool
}

func ConcatFolderUnsafe(afs *AppFSLayer.AppFS, folderIN int) []FolderBlock2 {
	folderINode := afs.GetFileINfo(folderIN)
	re := []FolderBlock2{}
	//maxI := 0
	for i, v := range folderINode.Pointers {
		re = append(re, FolderBlock2{})
		if v > 0 {
			//maxI = i
			////fmt.Println("Is it ok?")
			co := afs.ReadFile(folderIN, i).(FolderBlock).fileEntrys
			////fmt.Println("Ok!")
			for x, v := range co {

				re[i].F[x].Name = v.name
				re[i].F[x].Inode = v.inode
				re[i].F[x].Valid = v.valid
			}
		}
		//re = append(re, co[:]...)
	}

	//fmt.Println("Max i to:", maxI)
	return re
}
