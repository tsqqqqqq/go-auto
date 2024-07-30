package auto

import (
	"auto-record/app/event"
	"auto-record/app/template"
	"auto-record/config"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (ar *AutoRecord) Run() {
	fmt.Println("Run")
	filename := filepath.Join(config.Settings.FilePath.Record, template.CurrentTemplate, "text.log")
	ioReader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(ioReader)

	// FIXME 这里的协程在for循环里面启动的太多了，会导致系统卡顿。测试行数多的时候。应该设置一个锁或者协程group。限制协程的数量。

	// 方案1 用管道管理Text(). 分别启动两个协程，根据管道传递的类型执行不同的协程
	// OPTIMIZE 2024-07-30 使用协程的单向通道， 定义键盘和鼠标的channel消费者，在使用写入协程单向对channel写入，但感觉鼠标移动还是不够流畅，待测试。
	mouseChan := make(chan *event.MouseMoveEvent, 5)
	keyboardChan := make(chan *event.KeyboardEvent, 5)

	go event.MouseEventFormat(mouseChan)
	go event.KeyboardEventFormat(keyboardChan)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Mouse") {
				mouseChan <- event.NewMouseMoveEvent(line)
			} else {
				keyboardChan <- event.NewKeyboardEvent(line)
			}
		}
	}()

	// 方案2 用sync lock 和 unlock 锁住协程， 用group限制协程数量

	//for scanner.Scan() {
	//	if strings.Contains(scanner.Text(), "HookEnabled") {
	//		continue
	//	}
	//	go event.MouseEventFormat(scanner.Text())
	//	go event.KeyboardEventFormat(scanner.Text())
	//}

	// 方案3 不用协程直接运行 和协程一样慢
	//for scanner.Scan() {
	//	line := scanner.Text()
	//	if strings.Contains(line, "Mouse") {
	//		mouse := event.NewMouseMoveEvent(line)
	//		mouse.MouseEventFormat()
	//		//mouseChan <- event.NewMouseMoveEvent(line)
	//	} else {
	//		//keyboardChan <- event.NewKeyboardEvent(line)
	//	}
	//}

}
