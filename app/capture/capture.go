package capture

import (
	"auto-record/app/template"
	"auto-record/app/window"
	"auto-record/config"
	"context"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"path/filepath"
	"time"
)

type Capture struct {
	ctx context.Context
}

// CaptureCount 每当有一个键盘鼠标得输入事件的时候， 按照时间戳截图一张当前窗口句柄。。。
var CaptureCount int = 0

const ext = "png"

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
		if err := robotgo.Save(img, imgFile); err != nil {
			fmt.Println("==========>", err)
		}
		CaptureCount++
	}
}

func (c *Capture) TestEvent() {
	fmt.Println("runtime test")
	count := 0
	for {
		fmt.Printf("current count: %d\n", count)
		runtime.EventsEmit(c.ctx, "test", count)
		time.Sleep(time.Second * 3)
		count++
	}
}
