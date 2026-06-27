package main

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/akke/ease-ui/internal/app"
)

//go:embed all:frontend/dist
var assets embed.FS

func runApp() error {
	a, err := app.New(app.Options{})
	if err != nil {
		return err
	}

	err = wails.Run(&options.App{
		Title:     "Ease",
		Width:     1100,
		Height:    720,
		MinWidth:  900,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0x0A, G: 0x0A, B: 0x0A, A: 1},
		OnStartup: func(ctx context.Context) {
			wailsruntime.LogInfo(ctx, "ease-ui starting, version "+version)
		},
		Bind: []interface{}{ a },
	})
	return err
}
