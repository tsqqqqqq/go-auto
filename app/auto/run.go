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
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "HookEnabled") {
			continue
		}
		// FIXME 这里的协程在for循环里面启动的太多了，会导致系统卡顿。测试行数多的时候。应该设置一个锁或者协程group。限制协程的数量。
		go event.MouseEventFormat(scanner.Text())
		go event.KeyboardEventFormat(scanner.Text())
	}

}
