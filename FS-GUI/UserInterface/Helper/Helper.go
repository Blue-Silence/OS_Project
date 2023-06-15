package Helper

import (
	"LSF/AppFSLayer"
	"LSF/UserInterface/File"
	"fmt"
	"strings"
)

func getParentHandle(afs *AppFSLayer.AppFS, path string) (string, File.FileHandler) {
	names := strings.Split(path, "/")
	if len(names) < 1 {
		return "Not allowed.", File.FileHandler{}
	}
	pathParent := ""
	for _, v := range names[0 : len(names)-1] {
		pathParent = fmt.Sprint(pathParent, "/", v)
	}

	err, hanP := File.GetHandler(afs, pathParent)
	return err, hanP
}

func CreateByPath(afs *AppFSLayer.AppFS, path string, fileType int) (string, File.FileHandler) {
	names := strings.Split(path, "/")
	err, hanP := getParentHandle(afs, path)
	if err != "" {
		return err, File.FileHandler{}
	}
	newName := names[len(names)-1]
	_, lts := File.GetFolderContent(afs, hanP)
	for _, v := range lts {
		if v.Name == newName {
			return fmt.Sprint("File already exist:", newName), File.FileHandler{}
		}
	}

	return File.Create(afs, hanP, newName, fileType)

}

func DeleteByPath(afs *AppFSLayer.AppFS, path string) string {
	err, p := File.GetHandler(afs, path)
	if err != "" {
		return err
	}
	err, hanP := getParentHandle(afs, path)
	if err != "" {
		return err
	}
	return File.Delete(afs, hanP, p)

}
