package main

import (
	"errors"
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
func sourcesHandler(c *gin.Context) {
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
func albumsHandler(c *gin.Context) {
	api_key := c.Query("api_key")
	source, err := sources.New(c.Param("sourceName"), api_key, &SimpleStorage{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	albums, err := source.Albums()
	if err != nil {
		var e *sources.AuthError
		if errors.As(err, &e) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"albums": albums})
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
func downloadAlbumHandler(c *gin.Context) {
	api_key := c.Query("api_key")
	source, err := sources.New(c.Param("sourceName"), api_key, &SimpleStorage{c.Query("dir")})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dir, err := source.DownloadAlbum(c.Param("albumID"))
	if err != nil {
		var e *sources.AuthError
		if errors.As(err, &e) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": "", "dir": dir})
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
func downloadAllAlbumsHandler(c *gin.Context) {
	api_key := c.Query("api_key")
	source, err := sources.New(c.Param("sourceName"), api_key, &SimpleStorage{c.Query("dir")})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dir, err := source.DownloadAllAlbums()
	if err != nil {
		var e *sources.AuthError
		if errors.As(err, &e) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": "", "dir": dir})
}
