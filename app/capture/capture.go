package capture

import (
	"auto-record/app/template"
	"auto-record/app/window"
	"auto-record/config"
	"fmt"
	"github.com/go-vgo/robotgo"
	"path/filepath"
)

// CaptureCount 每当有一个键盘鼠标得输入事件的时候， 按照时间戳截图一张当前窗口句柄。。。
var CaptureCount int = 0

const ext = "png"

func WindowCapture(input chan string) {
	for _ = range input {
		captureName := fmt.Sprintf("capture_%d.%s", CaptureCount, ext)
		imgFile := filepath.Join(config.Settings.FilePath.Record, template.CurrentTemplate, captureName)
		win := window.CurrentWindow
		robotgo.SetActive(robotgo.GetHandPid(win.Pid))
		x, y, w, h := robotgo.GetBounds(win.Pid)
		img := robotgo.CaptureImg(x, y, w, h)
		if err := robotgo.Save(img, imgFile); err != nil {
			fmt.Println("==========>", err)
		}
		CaptureCount++
	}
}
