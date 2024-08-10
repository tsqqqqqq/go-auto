package window

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type Window struct {
	Pid   int
	Title string
	HWND  int
}

func NewWindow() *Window {
	return &Window{}
}

func (w *Window) GetWindows() []*Window {
	pids, err := robotgo.Pids()
	if err != nil {
		panic(err)
	}
	wins := make([]*Window, 0)
	for _, pid := range pids {
		// 这里会出现中文乱码 用[]rune 转换一下
		name := robotgo.GetTitle(pid)
		title, err := simplifiedchinese.GB18030.NewDecoder().String(name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		hwnd := robotgo.GetHWNDByPid(pid)
		if hwnd == 0 || title == "" {
			continue
		}
		wins = append(wins, &Window{
			Pid:   pid,
			Title: title,
			HWND:  hwnd,
		})
	}
	return wins
}
