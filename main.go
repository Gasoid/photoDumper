package main

import (
	"embed"

	_ "github.com/Gasoid/photoDumper/docs"
	"github.com/Gasoid/photoDumper/sources"
	"github.com/Gasoid/photoDumper/sources/instagram"
	"github.com/Gasoid/photoDumper/sources/vk"

	local "github.com/Gasoid/photoDumper/storage/localfs"
)

//go:embed build/*
var staticAssets embed.FS

// @title        PhotoDumper
// @version      1.1.2
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
	router := setupRouter()
	if router != nil {
		router.Run(":8080")
	}
}
