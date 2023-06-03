package File

import (
	"LSF/AppFSLayer"
	"LSF/BlockLayer"
)

func GetFolderContent(afs *AppFSLayer.AppFS, h FileHandler) (string, []FileInfo) {
	hN := afs.GetFileINfo(h.inode)

	if !hN.Valid {
		return "No such file.", []FileInfo{}
	}
	if hN.FileType != BlockLayer.Folder {
		return "Not a folder.", []FileInfo{}
	}
	return "", getFolderContentH(afs, h.inode)
}
