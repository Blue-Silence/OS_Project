package File

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
	FileEntrys [MaxFilePerFolderBlock]FileEntry
} // 1 per block
type FileEntry struct {
	Name  string
	Inode int
	Valid bool
}

func (s FolderBlock) CanBeBlock() {
}

func concatFolder(afs *AppFSLayer.AppFS, folderIN int) []FileEntry {
	folderINode := afs.GetFileINfo(folderIN)
	re := []FileEntry{}
	//maxI := 0
	for i, v := range folderINode.Pointers {
		if v >= 0 {
			//maxI = i
			////fmt.Println("Is it ok?")
			co := FolderBlock{}.FromBlock(afs.ReadFile(folderIN, i)).(FolderBlock).FileEntrys
			////fmt.Println("Ok!")
			for _, v := range co {
				if v.Valid {
					re = append(re, v)
				}
			}
			//re = append(re, co[:]...)
		}
	}
	//fmt.Println("Max i to:", maxI)
	return re
}

func rebuildFolder(fEsO []FileEntry) ([]int, []DiskLayer.Block) {
	returnIndex := []int{}
	returnBlock := []DiskLayer.Block{}
	fEs := []FileEntry{}
	for _, v := range fEsO {
		if v.Valid {
			fEs = append(fEs, v)
		}
	}
	i := 0
	for {
		if len(fEs) > 0 {
			returnIndex = append(returnIndex, i)
			i++
			fEB := FolderBlock{}
			copy(fEB.FileEntrys[:], fEs)
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

func addFileToFolder(afs *AppFSLayer.AppFS, folderIN int, fileIN int) {
	fEs := concatFolder(afs, folderIN)
	fileINode := afs.GetFileINfo(fileIN)
	folderINode := afs.GetFileINfo(folderIN)
	fE := FileEntry{Name: fileINode.Name, Inode: fileINode.InodeN, Valid: true}
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

func deleteFileToFolder(afs *AppFSLayer.AppFS, folderIN int, fileIN int) {
	fEs := concatFolder(afs, folderIN)
	folderINode := afs.GetFileINfo(folderIN)
	newFe := []FileEntry{}
	for _, v := range fEs {
		if v.Inode != fileIN {
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

func getFolderContentH(afs *AppFSLayer.AppFS, inode int) []FileInfo {
	fileLt := []FileInfo{}
	for _, v := range concatFolder(afs, inode) {
		if v.Valid {
			_, h := GetInfo(afs, FileHandler{v.Inode})
			fileLt = append(fileLt, h)
		}
	}
	return fileLt
}

// ////////////////////////////////////////////////////////////////////////////////
// / FOR TEST
type folderBlock2 struct {
	F [MaxFilePerFolderBlock]FileEntry2
} // 1 per block
type FileEntry2 struct {
	Name  string
	Inode int
	Valid bool
}

func ConcatFolderUnsafe(afs *AppFSLayer.AppFS, folderIN int) []folderBlock2 {
	folderINode := afs.GetFileINfo(folderIN)
	re := []folderBlock2{}
	//maxI := 0
	for i, v := range folderINode.Pointers {
		re = append(re, folderBlock2{})
		if v > 0 {
			//maxI = i
			////fmt.Println("Is it ok?")
			co := FolderBlock{}.FromBlock(afs.ReadFile(folderIN, i)).(FolderBlock).FileEntrys
			////fmt.Println("Ok!")
			for x, v := range co {

				re[i].F[x].Name = v.Name
				re[i].F[x].Inode = v.Inode
				re[i].F[x].Valid = v.Valid
			}
		}
		//re = append(re, co[:]...)
	}

	//fmt.Println("Max i to:", maxI)
	return re
}
