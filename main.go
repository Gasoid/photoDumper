package main

import (
	"embed"
	"net/http"
	"time"

	_ "github.com/Gasoid/photoDumper/docs"
	"github.com/Gasoid/photoDumper/sources"
	"github.com/Gasoid/photoDumper/sources/instagram"
	"github.com/Gasoid/photoDumper/sources/vk"
	"github.com/pkg/browser"

	local "github.com/Gasoid/photoDumper/storage/localfs"
)

type engine interface {
	Run(addr ...string) error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

var (
	//go:embed build/*
	staticAssets    embed.FS
	setupRouterFunc func() engine = setupRouter
	openBrowserFunc func(string)  = openBrowser
)

// @title        PhotoDumper
// @version      1.2.0
// @description  app downloads photos from vk.

// @contact.name  Rinat Almakhov
// @contact.url   https://github.com/Gasoid/

// @license.name  MIT License
// @license.url   https://github.com/Gasoid/photoDumper/blob/main/LICENSE

// @host      localhost:8080
// @BasePath  /api/
// @securitydefinitions.apikey ApiKeyAuth
// @in query
// @name api_key
func main() {
	sources.AddSource(vk.NewService())
	sources.AddSource(instagram.NewService())
	sources.AddStorage(local.NewService())
	router := setupRouterFunc()
	if router != nil {
		go openBrowserFunc("http://localhost:8080")
		router.Run(":8080")
	}
}

func openBrowser(url string) {
	time.Sleep(time.Second * 5)
	browser.OpenURL(url)
}
