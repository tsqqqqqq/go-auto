package capture

import (
	"auto-record/app/template"
	"auto-record/app/window"
	"auto-record/config"
	"context"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"image"
	"path/filepath"
)

type Capture struct {
	ctx context.Context
}

// CaptureCount 每当有一个键盘鼠标得输入事件的时候， 按照时间戳截图一张当前窗口句柄。。。
var CaptureCount int = 0
var ImageChan chan image.Image

const ext = "png"

func init() {
	if ImageChan == nil {
		ImageChan = make(chan image.Image, 2)
	}
}

func NewCapture(ctx context.Context) *Capture {
	return &Capture{ctx}
}

func WindowCapture(input chan string) {
	for _ = range input {
		captureName := fmt.Sprintf("capture_%d.%s", CaptureCount, ext)
		imgFile := filepath.Join(config.Settings.FilePath.Record, template.CurrentTemplate, captureName)
		win := window.CurrentWindow
		robotgo.SetActive(robotgo.GetHandPid(win.Pid))
		x, y, w, h := robotgo.GetBounds(win.Pid)
		img := robotgo.CaptureImg(x, y, w, h)
		ImageChan <- img
		if err := robotgo.Save(img, imgFile); err != nil {
			fmt.Println("==========>", err)
		}
		CaptureCount++
	}
}

func (c *Capture) CaptureEvent(input chan image.Image) {
	for img := range input {
		runtime.EventsEmit(c.ctx, "test", img)
	}
}
