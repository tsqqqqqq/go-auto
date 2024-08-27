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
)

type Capture struct {
	ctx context.Context
}

// CaptureCount 每当有一个键盘鼠标得输入事件的时候， 按照时间戳截图一张当前窗口句柄。。。
var CaptureCount int = 0
var ImageChan chan string

const ext = "png"

func init() {
	if ImageChan == nil {
		ImageChan = make(chan string, 2)
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
		ImageChan <- imgFile
		if err := robotgo.Save(img, imgFile); err != nil {
			fmt.Println("==========>", err)
		}
		CaptureCount++
	}
}

func (c *Capture) CaptureEvent(input chan string) {
	for imgFile := range input {
		//base64Buffer := bytes.NewBuffer(nil)
		//err := jpeg.Encode(base64Buffer, img, nil)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		////fmt.Println(base64Buffer.String())
		////fmt.Println(len(base64Buffer.Bytes()))
		//dist := make([]byte, 5000000)                         //开辟存储空间
		//base64.StdEncoding.Encode(dist, base64Buffer.Bytes()) //buff转成base64
		// FIXME 这里可能最佳解决方案还是调用window capture 来录屏 周末搞一下

		// 如果不行，那就使用wails的动态资产获取文件路径
		runtime.EventsEmit(c.ctx, "capture", imgFile)
	}
}
