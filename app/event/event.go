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

// MouseEventFormat 将监听到的日志 格式化写到一个文件中
func MouseEventFormat(text string) {
	input := strings.Split(text, "- Event: ")[1]
	mObj := new(MouseMoveEvent)
	if strings.Contains(input, MOUSEAWAIT) {
		awaitMaps := make(map[string]string)
		fmt.Println(input)
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
	if mObj.Kind == MOUSEMOVE {
		robotgo.MouseSleep = 50 // 100 millisecond
		robotgo.Move(mObj.X, mObj.Y)
	}
	if mObj.Kind == MOUSEDOWN {
		robotgo.MouseSleep = 100
		robotgo.Move(mObj.X, mObj.Y)
		if mObj.Button == "1" {
			robotgo.Click()
		} else if mObj.Button == "2" {
			robotgo.Click("right")
		}
	}
	if mObj.Kind == MOUSEAWAIT {
		robotgo.MouseSleep = 50
		time.Sleep(mObj.Sleep)
	}
	if mObj.Kind == MOUSEWHEEL {
		robotgo.MouseSleep = 50
		robotgo.ScrollDir(int(math.Abs(float64(mObj.Rotation))), getDir(mObj.Rotation))
	}
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

func KeyboardEventFormat(text string) {
	input := strings.Split(text, "- Event: ")[1]
	kObj := new(KeyboardEvent)
	re := regexp.MustCompile(`(\b\w+\b): (\w+)`)
	result := re.ReplaceAllString(input, `"$1": "$2"`)
	err := json.Unmarshal([]byte(result), kObj)
	if err != nil {
		panic(err)
	}
	robotgo.KeySleep = 50
	err = robotgo.KeyDown(hook.RawcodetoKeychar(kObj.RawCode))
	if err != nil {
		panic(err)
	}
}
