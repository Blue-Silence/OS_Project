package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type File struct {
	Name     string
	IsDir    bool
	FileInfo os.FileInfo
}

var (
	currentPath     string
	files           []File
	selectedFileIdx int
)

func main() {
	a := app.New()
	w := a.NewWindow("File Manager")
	w.Resize(fyne.NewSize(800, 600))

	// 创建当前路径标签
	currentPathLabel := widget.NewLabel("Current Path: " + currentPath)
	log.Println("1", currentPath)

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
	log.Println("2", currentPath)

	// 创建顶部工具栏
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			// 打开文件夹
			FolderDialog(w, list, currentPathLabel)
		}),
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			// 打开文件夹
			currentPath = filepath.Dir(currentPath)
			updateFileList(list, currentPathLabel)
		}),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			// 创建文件
			showCreateFileDialog(w, list, currentPathLabel)
		}),
		widget.NewToolbarAction(theme.FolderNewIcon(), func() {
			// 创建文件夹
			showCreateFolderDialog(w, list, currentPathLabel)
		}),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			// 删除文件或文件夹
			if selectedFileIdx >= 0 && selectedFileIdx < len(files) {
				file := files[selectedFileIdx]
				showDeleteFileDialog(w, file, list, currentPathLabel)
			}
		}),
		widget.NewToolbarAction(theme.ConfirmIcon(), func() {
			// 打开或查看文件
			file := files[selectedFileIdx]
			if file.IsDir {
				// 如果选中的是文件夹，则进入文件夹
				log.Println("before", currentPath)
				currentPath = filepath.Join(currentPath, file.Name)
				log.Println("After", currentPath)
				updateFileList(list, currentPathLabel)
			} else {
				// 如果选中的是文本文件，则打开文件
				if true {
					content, err := ioutil.ReadFile(filepath.Join(currentPath, file.Name))
					if err != nil {
						dialog.ShowError(err, w)
						return
					}
					showFileInfoDialog(w, file, string(content), func(newContent string) {
						err := ioutil.WriteFile(filepath.Join(currentPath, file.Name), []byte(newContent), 0644)
						if err != nil {
							dialog.ShowError(err, w)
							return
						}
					})
				}
			}
		}),
		widget.NewToolbarAction(theme.InfoIcon(), func() {
			// 查看文件属性
			if selectedFileIdx >= 0 && selectedFileIdx < len(files) {
				file := files[selectedFileIdx]
				showFilePropertiesDialog(w, file)
			}
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
	currentPath, _ = os.Getwd()
	//currentPath = "/"
	updateFileList(list, currentPathLabel)

	w.ShowAndRun()
}

// 更新文件列表
func updateFileList(list *widget.List, currentPathLabel *widget.Label) {
	files = nil

	fileInfos, err := ioutil.ReadDir(currentPath)
	if err != nil {
		dialog.ShowError(err, nil)
		return
	}

	for _, fileInfo := range fileInfos {
		file := File{
			Name:     fileInfo.Name(),
			IsDir:    fileInfo.IsDir(),
			FileInfo: fileInfo,
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
				err := ioutil.WriteFile(path, []byte{}, 0644)
				if err != nil {
					dialog.ShowError(err, w)
					return
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
				err := os.Mkdir(path, 0755)
				if err != nil {
					dialog.ShowError(err, w)
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
			err := os.RemoveAll(path)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			updateFileList(list, currentPathLabel)
		}
	}, w)
}

// 显示文件属性对话框
func showFilePropertiesDialog(w fyne.Window, file File) {
	var properties string

	properties += fmt.Sprintf("Name: %s\n", file.Name)
	properties += fmt.Sprintf("Size: %d bytes\n", file.FileInfo.Size())
	properties += fmt.Sprintf("Modified: %s\n", file.FileInfo.ModTime())
	properties += fmt.Sprintf("Permissions: %s\n", file.FileInfo.Mode().Perm())

	dialog.ShowInformation("File Properties", properties, w)
}

// 打开文件夹对话框
func FolderDialog(w fyne.Window, list *widget.List, currentPathLabel *widget.Label) {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err == nil && uri != nil {
			log.Println("before", currentPath)
			currentPath = uri.String()
			log.Println("after", currentPath)
			updateFileList(list, currentPathLabel)
		}
	}, w)
}
