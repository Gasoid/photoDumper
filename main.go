package main

import (
	_ "github.com/Gasoid/photoDumper/docs"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title        PhotoDumper
// @version      1.0
// @description  app downloads photos from vk.

// @contact.name  Rinat Almakhov
// @contact.url   https://github.com/Gasoid/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/
// @securitydefinitions.apikey ApiKeyAuth
// @in query
// @name api_key
func main() {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/sources/", getSources)
		auth := api.Group("/", Auth())
		{
			auth.GET("/albums/:sourceName/", getAlbums)
			auth.GET("/album-photos/:sourceName/:albumID/", getAlbumPhotos)
			auth.GET("/download-all-albums/:sourceName/", downloadAllAlbums)
			auth.GET("/download-album/:sourceName/:albumID/", downloadAlbum)
		}

	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}
