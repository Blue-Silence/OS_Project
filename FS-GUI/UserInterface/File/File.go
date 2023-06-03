package File

import (
	"LSF/AppFSLayer"
	"LSF/BlockLayer"
	"strings"
)

func GetHandler(afs *AppFSLayer.AppFS, path string) (string, FileHandler) {
	namesO := strings.Split(path, "/")
	names := []string{}
	for _, v := range namesO {
		if len(v) > 0 {
			names = append(names, v)
		}
	}
	inodeNC := 0
	for {
		inodeN := afs.GetFileINfo(inodeNC)
		if !inodeN.Valid {
			return "File not exist.", FileHandler{}
		}
		if len(names) == 0 {
			break
		}
		if inodeN.FileType != BlockLayer.Folder {
			return "Meet file(not folder in the path)", FileHandler{}
		}
		fileMap := make(map[string]int)
		_, fI := GetFolderContent(afs, FileHandler{inodeNC})
		for _, v := range fI {
			fileMap[v.Name] = v.Handler.inode
		}
		inodeNC = fileMap[names[0]]
		names = names[1:]
		if (inodeNC) == 0 {
			return "No such file", FileHandler{}
		}
	}
	return "", FileHandler{inodeNC}
}

func GetInfo(afs *AppFSLayer.AppFS, h FileHandler) (string, FileInfo) {
	info := FileInfo{}
	hN := afs.GetFileINfo(h.inode)
	if !hN.Valid {
		return "No such file.", info
	}

	info.Name = hN.Name
	info.FileType = hN.FileType
	info.Handler = FileHandler{hN.InodeN}
	for i, v := range hN.Pointers {
		if v > 0 {
			info.AllocatedBlock = append(info.AllocatedBlock, i)
			info.SizeInBlock++
		}
	}
	return "", info
}

func Create(afs *AppFSLayer.AppFS, parentH FileHandler, name string, fileType int) (string, FileHandler) {

	newI := afs.CreateFile(fileType, name)
	if !(newI > 0) {
		return "Create file fail.", FileHandler{}
	}
	err, pInfo := GetInfo(afs, parentH)
	if err != "" {
		return err, FileHandler{}
	}
	if pInfo.FileType != BlockLayer.Folder {
		return "Not a folder.", FileHandler{}
	}

	addFileToFolder(afs, parentH.inode, newI)
	return "", FileHandler{newI}
}

func Delete(afs *AppFSLayer.AppFS, parentH FileHandler, selfH FileHandler) string {
	err, pInfo := GetInfo(afs, parentH)
	if err != "" {
		return err
	}
	if pInfo.FileType != BlockLayer.Folder {
		return "Not a folder."
	}
	err, selfInfo := GetInfo(afs, selfH)
	if err != "" {
		return err
	}
	if selfInfo.FileType == BlockLayer.Folder {
		_, fileLt := GetFolderContent(afs, selfH)
		for _, v := range fileLt {
			Delete(afs, selfH, v.Handler)
		}
	}

	deleteFileToFolder(afs, parentH.inode, selfH.inode)

	return ""
}

func Flush(afs *AppFSLayer.AppFS) {
	afs.LogCommit()
}
