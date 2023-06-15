package File

import "time"

type FileHandler struct {
	inode int
}

type FileInfo struct {
	Name           string
	FileType       int
	SizeInBlock    int
	AllocatedBlock []int
	Handler        FileHandler
	ModTime        time.Time
}
