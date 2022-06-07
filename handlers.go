package main

import (
	"net/http"

	"github.com/Gasoid/photoDumper/sources"
	"github.com/gin-gonic/gin"
)

// getSources godoc
// @Summary      Sources
// @Description  get sources list
// @Produce      json
// @Accept       json
// @Success      200  {array}  string  "getSources"
// @Router       /sources/ [get]
func getSources(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"sources": sources.Sources()})
}

// getAlbums godoc
// @Summary      Albums
// @Description  get albums list
// @Produce      json
// @Accept       json
// @Param        sourceName  path     string  true  "source name"
// @Success      200         {array}  string  "albums"
// @Security     ApiKeyAuth
// @Router       /albums/{sourceName}/ [get]
func getAlbums(c *gin.Context) {
	api_key := c.Query("api_key")
	source, err := sources.New(c.Param("sourceName"), api_key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	albums, err := source.GetAlbums()
	if err != nil {
		if source.IsAuthError(err) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"albums": albums})
}

// getAlbumPhotos godoc
// @Summary      Album Photos
// @Description  get photos of album
// @Produce      json
// @Accept       json
// @Param        sourceName  path     string  true  "source name"
// @Param        albumID     path     string  true  "album ID"
// @Success      200         {array}  string  "getAlbumPhotos"
// @Router       /album-photos/{albumID}/{sourceName}/ [get]
// @Security     ApiKeyAuth
func getAlbumPhotos(c *gin.Context) {
	api_key := c.Query("api_key")
	source, err := sources.New(c.Param("sourceName"), api_key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	photos, err := source.GetAlbumPhotos(c.Param("albumID"))
	if err != nil {
		if source.IsAuthError(err) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"photos": photos})
}

// downloadAlbum godoc
// @Summary      download photos of album
// @Description  download all photos of particular album
// @Produce      json
// @Accept       json
// @Param        sourceName  path     string  true  "source name"
// @Param        albumID     path     string  true  "album ID"
// @Param        dir         query    string  true  "directory where photos will be stored"
// @Success      200         {array}  string
// @Router       /download-album/{albumID}/{sourceName}/ [get]
// @Security     ApiKeyAuth
func downloadAlbum(c *gin.Context) {
	api_key := c.Query("api_key")
	source, err := sources.New(c.Param("sourceName"), api_key)
	if err != nil {
		if source.IsAuthError(err) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	go source.DownloadAlbum(c.Param("albumID"), c.Query("dir"))
	c.JSON(http.StatusOK, gin.H{})
}

// downloadAllAlbums godoc
// @Summary      download photos of albums
// @Description  download all photos of all albums
// @Produce      json
// @Accept       json
// @Param        sourceName  path     string  true  "source name"
// @Param        dir         query    string  true  "directory where photos will be stored"
// @Success      200         {array}  string
// @Router       /download-all-albums/{sourceName}/ [get]
// @Security     ApiKeyAuth
func downloadAllAlbums(c *gin.Context) {
	api_key := c.Query("api_key")
	source, err := sources.New(c.Param("sourceName"), api_key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = source.DownloadAllAlbums(c.Query("dir"))
	if err != nil {
		if source.IsAuthError(err) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
