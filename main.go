package main

import (
	"auto-record/app/auto"
	template2 "auto-record/app/template"
	"auto-record/config"
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var record *auto.AutoRecord
var template *template2.Template

func main() {
	// Create an instance of the app structure
	app := NewApp()
	appInit()

	defer close(record.IsListen)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "auto-record",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		//BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		LogLevel:  logger.DEBUG,
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
			record,
			template,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func appInit() {
	record = auto.NewAutoRecord()
	template = template2.NewTemplate()

	go record.Listen(record.IsListen)
	recordFile := config.Settings.FilePath.Record
	fmt.Println(recordFile)
	if err := os.Mkdir(recordFile, os.ModePerm); os.IsNotExist(err) {
		panic(fmt.Errorf("create directory error %v", err))
	}
}
