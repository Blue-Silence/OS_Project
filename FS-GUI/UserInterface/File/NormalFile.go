package File

import (
	"LSF/AppFSLayer"
	"LSF/BlockLayer"
	"LSF/DiskLayer"
	"LSF/Setting"
)

func Write(afs *AppFSLayer.AppFS, h FileHandler, index int, data [Setting.BlockSize]uint8) string {
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
	return "", BlockLayer.DataBlock{}.FromBlock(afs.ReadFile(h.inode, index)).(BlockLayer.DataBlock).Data
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
