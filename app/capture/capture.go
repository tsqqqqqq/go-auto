package capture

import (
	"auto-record/app/window"
	"fmt"
)

// 每当有一个键盘鼠标得输入事件的时候， 按照时间戳截图一张当前窗口句柄。。。
//

func WindowCapture(input chan string) {
	for e := range input {
		fmt.Println(e, window.CurrentWindow)
	}
}
