package app

import (
	"os"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) OSUsername() string {
	if u := os.Getenv("USER"); u != "" {
		return u
	}
	if u := os.Getenv("USERNAME"); u != "" {
		return u
	}
	return ""
}

// PickDirectory 打开系统原生文件夹选择对话框，返回所选路径。
func (a *App) PickDirectory() (string, error) {
	return wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: "选择工作目录",
	})
}
