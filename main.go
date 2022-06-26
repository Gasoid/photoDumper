package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	_ "github.com/Gasoid/photoDumper/docs"
	"github.com/Gasoid/photoDumper/sources"
	"github.com/Gasoid/photoDumper/sources/vk"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	local "github.com/Gasoid/photoDumper/storage/localfs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	sources.AddStorage(local.NewService())
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	assets, err := fs.Sub(staticAssets, "build")
	if err != nil {
		fmt.Println("build folder is not readable")
		return
	}
	assetsFS := http.FS(assets)
	router := gin.Default()
	router.Use(cors.New(config))
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/assets/index.html")
	})
	router.StaticFS("/assets/", assetsFS)
	api := router.Group("/api")
	{
		api.GET("/sources/", sourcesHandler)
		auth := api.Group("/", Auth())
		{
			auth.GET("/albums/:sourceName/", albumsHandler)
			auth.GET("/download-all-albums/:sourceName/", downloadAllAlbumsHandler)
			auth.GET("/download-album/:albumID/:sourceName/", downloadAlbumHandler)
		}

	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}
