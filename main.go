package main

import (
	"auto-record/app/auto"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	record := auto.NewAutoRecord()
	defer close(record.IsListen)
	go record.Listen(record.IsListen)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "auto-record",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			record,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
