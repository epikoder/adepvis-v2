package main

import (
	"context"
	"embed"

	"github.com/epikoder/adepvis/src/service"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := service.NewApp()
	auth := service.NewAuth()
	scanner := service.NewScanner()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Adepvis",
		Width:  824, //1024,
		Height: 568, //768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.Startup(ctx)
			auth.Init(ctx)
			scanner.Init(ctx)
		},
		OnDomReady: func(ctx context.Context) {
			runtime.LogInfo(ctx, "DOM READY")
		},
		Bind: []interface{}{
			app,
			auth,
			scanner,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
