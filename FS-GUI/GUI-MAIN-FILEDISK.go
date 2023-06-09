package main

import (
	"LSF/AppFSLayer"
	"LSF/BlockLayer"
	"LSF/FileDisk"
	"LSF/Setting"
	"LSF/UserInterface/Helper"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	FF "LSF/UserInterface/File"
)

type File struct {
	Name    string
	IsDir   bool
	ModTime time.Time
}

var afs AppFSLayer.AppFS

var (
	currentPath     string
	files           []File
	selectedFileIdx int
)

func main() {

	disk := FileDisk.FileDisk{}
	nF := ""
	sF := ""
	fmt.Printf("Please enter the file as the normal block storage[Default:./NormalBlock]:")
	fmt.Scanf("%v", nF)
	fmt.Printf("Please enter the file as the super block storage[Default:./SuperBlock]:")
	fmt.Scanf("%v", sF)
	if nF == "" {
		nF = "NormalBlock"
	}
	if sF == "" {
		sF = "SuperBlock"
	}
	disk.Init(nF, sF)
	fmt.Printf("Do you want to reformat the disk?[Y/N][Default:N]:")
	c := 'a'
	_, _ = fmt.Scanf("%c", &c)

	if c == 'Y' || c == 'y' {
		afs.FormatFS(&disk)
	} else {
		afs.LoadFS(&disk)
	}

	defer disk.Close()

	a := app.New()
	w := a.NewWindow("File Manager")
	w.Resize(fyne.NewSize(800, 600))

	// 创建当前路径标签
	currentPathLabel := widget.NewLabel("Current Path: " + currentPath)

	// 创建主布局
	list := widget.NewList(
		func() int {
			return len(files)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(index int, obj fyne.CanvasObject) {
			file := files[index]
			obj.(*widget.Label).SetText(file.Name)
		},
	)

	// 设置文件列表的双击事件处理函数
	list.OnSelected = func(index int) {
		selectedFileIdx = index
	}

	// 创建顶部工具栏

	toolbar := container.NewHBox(
		widget.NewButtonWithIcon("Parent Folder", theme.FolderOpenIcon(), func() {
			// 打开文件夹
			currentPath = filepath.Dir(currentPath)
			log.Println("Open folder:", currentPath)
			updateFileList(list, currentPathLabel)
		}),
		widget.NewButtonWithIcon("Create File", theme.ContentAddIcon(), func() {
			// 创建文件
			showCreateFileDialog(w, list, currentPathLabel)
		}),
		widget.NewButtonWithIcon("Create Folder", theme.FolderNewIcon(), func() {
			// 创建文件夹
			showCreateFolderDialog(w, list, currentPathLabel)
		}),
		widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
			// 删除文件或文件夹
			if selectedFileIdx >= 0 && selectedFileIdx < len(files) {
				file := files[selectedFileIdx]
				showDeleteFileDialog(w, file, list, currentPathLabel)
			}
		}),
		widget.NewButtonWithIcon("Open", theme.ConfirmIcon(), func() {
			// 打开或查看文件
			if selectedFileIdx >= len(files) {
				return
			}
			file := files[selectedFileIdx]
			if file.IsDir {
				// 如果选中的是文件夹，则进入文件夹
				//log.Println("before", currentPath)
				currentPath = filepath.Join(currentPath, file.Name)
				//log.Println("After", currentPath)
				updateFileList(list, currentPathLabel)
			} else {
				// 如果选中的是文本文件，则打开文件
				if true {
					err1, Handler := FF.GetHandler(&afs, filepath.Join(currentPath, file.Name))
					if err1 != "" {
						log.Fatal("Err1:", err1)
					}
					content := getFullContent(Handler)
					showFileInfoDialog(w, file, string(content), func(newContent string) {
						indexs, blocks := bytes2Block([]byte(newContent))
						for i, _ := range indexs {
							err := FF.Write(&afs, Handler, indexs[i], blocks[i])
							log.Println("Err:", err)
						}

					})
				}
			}
		}),
		widget.NewButtonWithIcon("File Info", theme.InfoIcon(), func() {
			// 查看文件属性
			updateFileList(list, currentPathLabel)
			if selectedFileIdx >= 0 && selectedFileIdx < len(files) {
				file := files[selectedFileIdx]
				showFilePropertiesDialog(w, file)
			}
		}),
		widget.NewButtonWithIcon("Flush", theme.CheckButtonCheckedIcon(), func() {
			// Flush
			FF.Flush(&afs)
		}),
	)

	// 创建底部状态栏
	statusBar := widget.NewLabel("")

	topW := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), toolbar, currentPathLabel)
	// 创建主布局容器
	content := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(topW, nil, statusBar, nil),
		topW,
		list,
		statusBar,
	)

	w.SetContent(content)

	// 初始化文件列表
	currentPath = "/"
	updateFileList(list, currentPathLabel)

	w.ShowAndRun()
	log.Println("OKKKKKKKKKKKKKKKKKKKK")
}

// 更新文件列表
func updateFileList(list *widget.List, currentPathLabel *widget.Label) {
	files = nil

	/*fileInfos, err := ioutil.ReadDir(currentPath)
	if err != nil {
		dialog.ShowError(err, nil)
		return
	}*/
	err1, h := FF.GetHandler(&afs, currentPath)
	err2, fileInfos := FF.GetFolderContent(&afs, h)
	if err1 != "" || err2 != "" {
		log.Fatal("Err1:", err1, "  Err2:", err2)
	}
	for _, fileInfo := range fileInfos {
		file := File{
			Name: fileInfo.Name,
			//IsDir: fileInfo.FileType,
			ModTime: fileInfo.ModTime,
		}
		if fileInfo.FileType == BlockLayer.Folder {
			file.IsDir = true
		}
		if fileInfo.FileType == BlockLayer.Folder {
			file.IsDir = true
		}
		files = append(files, file)
	}

	// 更新列表显示
	list.Refresh()
	currentPathLabel.SetText("Current Path: " + currentPath)
}

// 显示文件信息对话框
func showFileInfoDialog(w fyne.Window, file File, content string, onSave func(newContent string)) {
	entry := widget.NewMultiLineEntry()
	entry.SetText(content)

	form := widget.NewForm()
	form.Append("Content", entry)

	dialog.ShowCustomConfirm(file.Name, "Save", "Cancel", form, func(confirmed bool) {
		if confirmed {
			newContent := entry.Text
			onSave(newContent)
		}
	}, w)
}

// 显示创建文件对话框
func showCreateFileDialog(w fyne.Window, list *widget.List, currentPathLabel *widget.Label) {
	entry := widget.NewEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: entry},
		},
	}

	dialog.ShowForm("Create File", "Create", "Cancel", []*widget.FormItem{form.Items[0]}, func(confirmed bool) {
		if confirmed {
			name := entry.Text
			if name != "" {
				path := filepath.Join(currentPath, name)
				/*err := ioutil.WriteFile(path, []byte{}, 0644)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}*/
				err, _ := Helper.CreateByPath(&afs, path, BlockLayer.NormalFile)
				log.Println("Create file:", path)
				if err != "" {
					dialog.NewError(errors.New(err), w).Show()

				}
				updateFileList(list, currentPathLabel)
			}
		}
	}, w)
}

// 显示创建文件夹对话框
func showCreateFolderDialog(w fyne.Window, list *widget.List, currentPathLabel *widget.Label) {
	entry := widget.NewEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: entry},
		},
	}

	dialog.ShowForm("Create Folder", "Create", "Cancel", []*widget.FormItem{form.Items[0]}, func(confirmed bool) {
		if confirmed {
			name := entry.Text
			if name != "" {
				path := filepath.Join(currentPath, name)
				/*err := os.Mkdir(path, 0755)
				err :=
				if err != nil {
					dialog.ShowError(err, w)
					return
				}*/
				err, _ := Helper.CreateByPath(&afs, path, BlockLayer.Folder)
				log.Println("Create folder:", path)
				if err != "" {
					dialog.ShowError(nil, w)
					return
				}
				updateFileList(list, currentPathLabel)
			}
		}
	}, w)
}

// 显示删除文件对话框
func showDeleteFileDialog(w fyne.Window, file File, list *widget.List, currentPathLabel *widget.Label) {
	dialog.ShowConfirm("Delete File", fmt.Sprintf("Are you sure you want to delete '%s'?", file.Name), func(confirmed bool) {
		if confirmed {
			path := filepath.Join(currentPath, file.Name)
			/*err := os.RemoveAll(path)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}*/
			Helper.DeleteByPath(&afs, path)
			log.Println("Delete file:", path)
			updateFileList(list, currentPathLabel)
		}
	}, w)
}

// 显示文件属性对话框
func showFilePropertiesDialog(w fyne.Window, file File) {
	var properties string

	properties += fmt.Sprintf("Name: %s\n", file.Name)
	//properties += fmt.Sprintf("Size: %d bytes\n", file.FileInfo.Size())
	properties += fmt.Sprintf("Modified: %s\n", file.ModTime)
	if file.IsDir {
		properties += fmt.Sprintf("Type: Dir\n")
	} else {
		properties += fmt.Sprintf("Type: Normal File\n")
	}

	dialog.ShowInformation("File Properties", properties, w)
}

// 打开文件夹对话框
func FolderDialog(w fyne.Window, list *widget.List, currentPathLabel *widget.Label) {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err == nil && uri != nil {
			//log.Println("before", currentPath)
			currentPath = uri.String()
			//log.Println("after", currentPath)
			updateFileList(list, currentPathLabel)
		}
	}, w)
}

func getFullContent(h FF.FileHandler) []byte {
	_, info := FF.GetInfo(&afs, h)
	//log.Println("Info:", info)
	if len(info.AllocatedBlock) == 0 {
		content := make([]byte, 0)
		return content
	}
	content := make([]byte, Setting.BlockSize*(info.AllocatedBlock[len(info.AllocatedBlock)-1]+1))
	for _, v := range info.AllocatedBlock {
		_, b := FF.Read(&afs, h, v)
		//log.Println("b:", b)
		copy(content[v*Setting.BlockSize:], b[:])
		//log.Println("b:", b)
		//log.Println("content:", content)
	}
	return content
}

func bytes2Block(data []byte) ([]int, [][Setting.BlockSize]byte) {
	re := make([][Setting.BlockSize]byte, 0)
	reI := make([]int, 0)
	c := 0
	for len(data) >= Setting.BlockSize {
		var b [Setting.BlockSize]byte
		copy(b[:], data)
		re = append(re, b)
		data = data[Setting.BlockSize:]
		reI = append(reI, c)
		c++
	}
	if len(data) > 0 {
		var b [Setting.BlockSize]byte
		copy(b[:], data)
		re = append(re, b)
		reI = append(reI, c)
	}
	return reI, re
}

type Wrap struct {
	a fyne.CanvasObject
}

func (s Wrap) ToolbarObject() fyne.CanvasObject {
	return s.a
}
