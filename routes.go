package main

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter() *gin.Engine {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	assets, err := fs.Sub(staticAssets, "build")
	if err != nil {
		fmt.Println("build folder is not readable")
		return nil
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
	return router
}
