package capture

import (
	"auto-record/app/template"
	"auto-record/app/window"
	"auto-record/config"
	"fmt"
	"github.com/go-vgo/robotgo"
	"path/filepath"
)

// 每当有一个键盘鼠标得输入事件的时候， 按照时间戳截图一张当前窗口句柄。。。
var CaptureCount int = 0

const ext = "Png"

// WindowCapture
func WindowCapture(input chan string) {
	for e := range input {
		fmt.Println(e)
		captureName := fmt.Sprintf("capture_%d.%s", CaptureCount, ext)
		imgFile := filepath.Join(config.Settings.FilePath.Record, template.CurrentTemplate, captureName)
		win := window.CurrentWindow
		x, y, w, h := robotgo.GetBounds(win.Pid)
		img := robotgo.CaptureImg(x, y, w, h)
		if err := robotgo.Save(img, imgFile); err != nil {
			fmt.Println(err)
		}
		CaptureCount++
		//x, y, w, h := robotgo.GetDisplayBounds(mainid)
		//name := "test.png"
		////name1 := "test_001.png"
		//img := robotgo.CaptureImg(x, y, w, h)
		//if err := robotgo.Save(img, name); err != nil {
		//	panic(err)
		//}
	}
}
