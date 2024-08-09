package event

import (
	"encoding/json"
	"fmt"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"math"
	"regexp"
	"strings"
	"time"
)

type MouseMoveEvent struct {
	Kind      string        `json:"Kind"`
	Button    string        `json:"Button"`
	X         int           `json:"X,string"`
	Y         int           `json:"Y,string"`
	Amount    int           `json:"Amount,string"`
	Rotation  uint          `json:"Rotation,string"`
	Direction uint          `json:"Direction"`
	Clicks    int           `json:"Clicks,string"`
	Sleep     time.Duration `json:"Sleep,string"`
}

const (
	MOUSEMOVE  = "MouseMove"
	MOUSEDOWN  = "MouseDown"
	MOUSEWHEEL = "MouseWheel"
	MOUSEAWAIT = "Await"
)

func NewMouseMoveEvent(text string) *MouseMoveEvent {
	input := strings.Split(text, "- Event: ")[1]
	mObj := new(MouseMoveEvent)
	if strings.Contains(input, MOUSEAWAIT) {
		awaitMaps := make(map[string]string)
		err := json.Unmarshal([]byte(input), &awaitMaps)
		if err != nil {
			fmt.Println("Error:", err)
		}
		mObj.Kind = awaitMaps["Kind"]
		durationStr := awaitMaps["Sleep"]
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			panic(err)
		}
		mObj.Sleep = duration

	} else {
		re := regexp.MustCompile(`(\b\w+\b): (\w+)`)
		result := re.ReplaceAllString(input, `"$1": "$2"`)
		err := json.Unmarshal([]byte(result), mObj)
		if err != nil {
			panic(err)
		}
	}
	return mObj
}

func NewKeyboardEvent(text string) *KeyboardEvent {
	input := strings.Split(text, "- Event: ")[1]
	kObj := new(KeyboardEvent)
	re := regexp.MustCompile(`(\b\w+\b): (\w+)`)
	result := re.ReplaceAllString(input, `"$1": "$2"`)
	err := json.Unmarshal([]byte(result), kObj)
	if err != nil {
		panic(err)
	}
	return kObj
}

// MouseEventFormat 将监听到的日志 格式化写到一个文件中
func MouseEventFormat(mouseChan chan *MouseMoveEvent) {

	for mouse := range mouseChan {
		// todo 这里的if else 太丑了 ， 搞个设计模式优化这里的代码
		if mouse.Kind == MOUSEMOVE {
			robotgo.MouseSleep = 1 // 100 millisecond
			robotgo.Move(mouse.X, mouse.Y)
		}
		if mouse.Kind == MOUSEDOWN {
			robotgo.MouseSleep = 5
			robotgo.Move(mouse.X, mouse.Y)
			if mouse.Button == "1" {
				robotgo.Click()
			} else if mouse.Button == "2" {
				robotgo.Click("right")
			}
		}
		if mouse.Kind == MOUSEAWAIT {
			robotgo.MouseSleep = 5
			time.Sleep(mouse.Sleep)
		}
		if mouse.Kind == MOUSEWHEEL {
			robotgo.MouseSleep = 5
			robotgo.ScrollDir(int(math.Abs(float64(mouse.Rotation))), getDir(mouse.Rotation))
		}
	}

}

func (mouse *MouseMoveEvent) MouseEventFormat() {
	//
	//for mouse := range mouseChan {
	// todo 这里的if else 太丑了 ， 搞个设计模式优化这里的代码
	if mouse.Kind == MOUSEMOVE {
		robotgo.MouseSleep = 1 // 100 millisecond
		robotgo.Move(mouse.X, mouse.Y)
	}
	if mouse.Kind == MOUSEDOWN {
		robotgo.MouseSleep = 5
		robotgo.Move(mouse.X, mouse.Y)
		if mouse.Button == "1" {
			robotgo.Click()
		} else if mouse.Button == "2" {
			robotgo.Click("right")
		}
	}
	if mouse.Kind == MOUSEAWAIT {
		robotgo.MouseSleep = 5
		time.Sleep(mouse.Sleep)
	}
	if mouse.Kind == MOUSEWHEEL {
		robotgo.MouseSleep = 5
		robotgo.ScrollDir(int(math.Abs(float64(mouse.Rotation))), getDir(mouse.Rotation))
	}
	//}

}

func getDir(rotation uint) string {
	if rotation > 0 {
		return "up"
	} else {
		return "down"
	}
}

type KeyboardEvent struct {
	Kind    string `json:"Kind"`
	Keycode uint16 `json:"keycode"`
	RawCode uint16 `json:"Rawcode,string"`
	KeyChar rune   `json:"Keychar,string"`
}

const (
	KEYDOWN = "KeyDown"
	KEYUP   = "KeyUp"
)

func KeyboardEventFormat(keyboardChan chan *KeyboardEvent) {
	for keyboard := range keyboardChan {
		// FIXME 键盘按键操作还有bug 需要修复

		robotgo.KeySleep = 50
		err := robotgo.KeyDown(hook.RawcodetoKeychar(keyboard.RawCode))
		if err != nil {
			panic(err)
		}
	}
}
