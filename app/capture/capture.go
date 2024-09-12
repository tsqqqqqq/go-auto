package capture

import (
	"auto-record/app/window"
	"context"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
	re "github.com/tsqqqqqq/GoRecord"
	"github.com/tsqqqqqq/GoRecord/winapi"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"os/signal"
	"syscall"
)

type Capture struct {
	ctx context.Context
}

var ImageChan chan string

func init() {
	if ImageChan == nil {
		ImageChan = make(chan string, 2)
	}
}

func NewCapture(ctx context.Context) *Capture {
	return &Capture{ctx}
}

func WindowCapture(input chan string) {
	var rdHwnd win.HWND
	for _ = range input {
		// TODO 这里还需要优化以下 最后在根据操作生成视频
		currentWindow := window.CurrentWindow
		robotgo.SetActive(robotgo.GetHandPid(currentWindow.Pid))
		rdHwnd = winapi.FindWindow(nil, winapi.MustUTF16PtrFromString(currentWindow.Title))
		if rdHwnd == 0 {
			win.MessageBox(0, winapi.MustUTF16PtrFromString("Could not find window"), winapi.MustUTF16PtrFromString("RDP Relative Input"), win.MB_ICONERROR)
		}
		handler := new(re.CaptureHandler)
		err := handler.StartCapture(rdHwnd)
		defer handler.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
	}

}

func (c *Capture) CaptureEvent(input chan string) {
	for imgFile := range input {

		// FIXME 这里可能最佳解决方案还是调用window capture 来录屏 周末搞一下
		// 如果不行，那就使用wails的动态资产获取文件路径
		runtime.EventsEmit(c.ctx, "capture", imgFile)
	}
}
