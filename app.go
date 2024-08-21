package main

import (
	capture2 "auto-record/app/capture"
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	eventInit(a)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func eventInit(app *App) {
	// TODO 新增一个events事件，但是这里的写法有点太没结构了 要改
	capture := capture2.NewCapture(app.ctx)

	go capture.CaptureEvent(capture2.ImageChan)
}
