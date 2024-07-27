package auto

import (
	"auto-record/app/template"
	"auto-record/config"
	"fmt"
	hook "github.com/robotn/gohook"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func (ar *AutoRecord) Listen(isListen chan bool) {
	//var evChan event.Event
	//evChan := hook.Start()
	defer hook.End()
	// 写文件
	for check := range isListen {
		if check {
			evChan := hook.Start()
			go eventOutput(evChan)
		} else {
			fmt.Println("stop listen")
			//hook.End()
		}
	}

}

func eventOutput(evChan chan hook.Event) {
	fmt.Println("start listen")
	filename := filepath.Join(config.Settings.FilePath.Record, template.CurrentTemplate, "text.log")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var lastTime *time.Time
	for ev := range evChan {
		input := strings.ReplaceAll(ev.String(), ev.When.String(), "")
		input = fmt.Sprintf("[%v] %v \r\n", ev.When, input)
		if lastTime == nil {
			lastTime = TimeFormat(input)
		}
		currentTime := TimeFormat(input)
		duration := currentTime.Sub(*lastTime)
		// 计算鼠标操作是不是连续的
		if duration > time.Second {
			awaitInput := fmt.Sprintf(`[%v] - Event: {"Kind": "Await", "Sleep": "%v"}`, currentTime, duration)
			awaitInput = fmt.Sprintf("%v \r\n", awaitInput)
			_, err := f.Write([]byte(awaitInput))
			if err != nil {
				panic(err)
			}
			lastTime = currentTime
		} else {
			// 管道输出到文件中
			_, err := f.Write([]byte(input))
			if err != nil {
				panic(err)
			}
		}

	}
}

func TimeFormat(input string) *time.Time {
	re := regexp.MustCompile(`\[(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+)`)
	match := re.FindStringSubmatch(input)
	if len(match) >= 2 {
		timeString := match[1]
		ans, err := time.Parse("2006-01-02 15:04:05.9999999", timeString)
		if err != nil {
			fmt.Println("Error:", err)
		}
		return &ans
	} else {
		fmt.Println("No match found")
	}
	return nil
}
