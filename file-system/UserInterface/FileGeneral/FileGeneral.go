package FileGeneral

import (
	"LSF/AppFSLayer"
	"LSF/BlockLayer"
	"LSF/DiskLayer"
	"LSF/Setting"
	"LSF/UserInterface/Folder"
	"fmt"
	"strings"
)

type FileHandler struct {
	inode int
}

type FileInfo struct {
	name           string
	fileType       int
	sizeInBlock    int
	allocatedBlock []int
}

func GetHandler(afs *AppFSLayer.AppFS, path string) (string, FileHandler) {
	//defer fmt.Println("EXIT")
	namesO := strings.Split(path, "/")
	names := []string{}
	for _, v := range namesO {
		if len(v) > 0 {
			names = append(names, v)
		}
	}
	inodeNC := 0
	//fmt.Println("Get:", path, "      ------    ", names)
	/*for i, v := range names {
		fmt.Println(i, "th str is:", v, "$ len:", len(v))
	}*/
	for {
		inodeN := afs.GetFileINfo(inodeNC)
		if !inodeN.Valid {
			//fmt.Println("What?")
			return "File not exist.", FileHandler{}
		}
		if len(names) == 0 {
			//fmt.Println("OKKKKK")
			break
		}
		//fmt.Println("Len:", len(names), "      ------    ", names[0][len], "      ------    ", names[1])
		if inodeN.FileType != BlockLayer.Folder {
			//fmt.Println("emmmm?")
			return "Meet file(not folder in the path)", FileHandler{}
		}
		fileMap := Folder.GetFolderContent(afs, inodeNC)
		//fmt.Println("Right")
		inodeNC = fileMap[names[0]]
		names = names[1:]
		if (inodeNC) == 0 {
			//fmt.Println("Fuck")
			return "No such file", FileHandler{}
		}
	}
	//fmt.Println("OK...")
	return "", FileHandler{inodeNC}
}

func getParentHandle(afs *AppFSLayer.AppFS, path string) (string, FileHandler) {
	names := strings.Split(path, "/")
	if len(names) < 1 {
		return "Not allowed.", FileHandler{}
	}
	pathParent := ""
	for _, v := range names[0 : len(names)-1] {
		pathParent = fmt.Sprint(pathParent, "/", v)
	}

	err, hanP := GetHandler(afs, pathParent)
	return err, hanP
}

func Create(afs *AppFSLayer.AppFS, path string, fileType int) (string, FileHandler) {
	names := strings.Split(path, "/")
	err, hanP := getParentHandle(afs, path)
	if err != "" {
		return err, FileHandler{}
	}

	newI := afs.CreateFile(fileType, names[len(names)-1])
	if !(newI > 0) {
		return "Create file fail.", FileHandler{}
	}
	pN := afs.GetFileINfo(hanP.inode)
	if pN.FileType != BlockLayer.Folder {
		return "Not a folder.", FileHandler{}
	}

	Folder.AddFileToFolder(afs, hanP.inode, newI)
	return "", FileHandler{newI}
}

func Delete(afs *AppFSLayer.AppFS, path string) string {
	err, p := GetHandler(afs, path)
	if err != "" {
		return err
	}
	err, hanP := getParentHandle(afs, path)
	if err != "" {
		return err
	}
	//p := afs.GetFileINfo(hanP.inode)
	if afs.GetFileINfo(p.inode).FileType == BlockLayer.Folder {
		fileL := Folder.GetFolderContent(afs, p.inode)
		for childN, _ := range fileL {
			Delete(afs, fmt.Sprint(path, "/", childN))
			//afs.DeleteFile(in)
		}
		//return "Not a folder."
	}

	Folder.DeleteFileToFolder(afs, hanP.inode, p.inode)

	return ""
}

func Write(afs *AppFSLayer.AppFS, h FileHandler, index int, data [Setting.BlockSize]uint8) string {
	//var dataB DiskLayer.Block = BlockLayer.DataBlock{data}
	hN := afs.GetFileINfo(h.inode)
	if !hN.Valid {
		return "No such file."
	}
	if hN.FileType != BlockLayer.NormalFile {
		return "Can't read this type of file."
	}
	afs.WriteFile(h.inode, []int{index}, []DiskLayer.Block{BlockLayer.DataBlock{data}})
	return ""
}

func Read(afs *AppFSLayer.AppFS, h FileHandler, index int) (string, [Setting.BlockSize]uint8) {
	hN := afs.GetFileINfo(h.inode)
	if !hN.Valid {
		return "No such file.", [Setting.BlockSize]uint8{}
	}
	if hN.FileType != BlockLayer.NormalFile {
		return "Can't read this type of file.", [Setting.BlockSize]uint8{}
	}
	return "", afs.ReadFile(h.inode, index).(BlockLayer.DataBlock).Data
}

func DeleteBlock(afs *AppFSLayer.AppFS, h FileHandler, index []int) string {
	hN := afs.GetFileINfo(h.inode)
	if !hN.Valid {
		return "No such file."
	}
	if hN.FileType != BlockLayer.NormalFile {
		return "Can't mod this type of file."
	}
	afs.DeleteBlockInFile(hN.InodeN, index)
	return ""
}

func GetFolderContent(afs *AppFSLayer.AppFS, h FileHandler) (string, []string) {
	hN := afs.GetFileINfo(h.inode)
	re := []string{}
	if !hN.Valid {
		return "No such file.", re
	}
	if hN.FileType != BlockLayer.Folder {
		return "Not a folder.", re
	}

	for name, _ := range Folder.GetFolderContent(afs, hN.InodeN) {
		re = append(re, name)
	}
	return "", re
}

func GetInfo(afs *AppFSLayer.AppFS, h FileHandler) (string, FileInfo) {
	info := FileInfo{}
	hN := afs.GetFileINfo(h.inode)
	if !hN.Valid {
		return "No such file.", info
	}

	info.name = hN.Name
	info.fileType = hN.FileType
	for i, v := range hN.Pointers {
		if v > 0 {
			info.allocatedBlock = append(info.allocatedBlock, i)
			info.sizeInBlock++
		}
	}
	return "", info
}
