package Folder



const (
	MaxFilePerFolderBlock = (BlockSize / 512) //this can be changed.
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

func (s FolderBlock) canBeBlock() {
}

func (afs *AppFS) concatFolder(floderIN int) []fileEntry {
	folderINode := afs.getFileINfo(folderIN)
	re := []fileEntry{}
	for _,v := folderINode.
}


func (afs *AppFS) AddFileToFolder(floderIN int, fileIN int) {
	fileINode := afs.getFileINfo(fileIN)
	folderINode := afs.getFileINfo(folderIN)
	fE := fileEntry{name: fileINode.name, inode: fileINode.inodeN, valid: true}
	p := -1
	for i,v := range folderINode.pointers {
		if
	}
}
