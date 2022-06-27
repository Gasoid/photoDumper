package main

import (
	"errors"
	"net/http"

	"github.com/Gasoid/photoDumper/sources"
	"github.com/gin-gonic/gin"
)

// sourcesHandler godoc
// @Summary      Sources
// @Description  returns sources
// @Produce      json
// @Accept       json
// @Success      200  {array}  string  "sources"
// @Router       /sources/ [get]
func sourcesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"sources": sources.Sources()})
}

// albumsHandler godoc
// @Summary      Albums
// @Description  returns albums
// @Produce      json
// @Accept       json
// @Param        sourceName  path     string  true  "source name"
// @Success      200         {array}  string
// @Failure      400         {string}  string    "error"
// @Failure      401         {string}  string    "error"
// @Failure      403         {string}  string    "error"
// @Failure      500         {string}  string    "error"
// @Security     ApiKeyAuth
// @Router       /albums/{sourceName}/ [get]
func albumsHandler(c *gin.Context) {
	source, err := sources.New(c.Param("sourceName"), c.Query("api_key"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	albums, err := source.Albums()
	if err != nil {
		var e *sources.AccessError
		if errors.As(err, &e) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"albums": albums})
}

// downloadAlbumHandler godoc
// @Summary      download photos of album
// @Description  download all photos of particular album, returns destination of your photos
// @Produce      json
// @Accept       json
// @Param        sourceName  path     string  true  "source name"
// @Param        albumID     path     string  true  "album ID"
// @Param        dir         query    string  true  "directory where photos will be stored"
// @Success      200         {array}  string
// @Failure      400         {string}  string    "error"
// @Failure      401         {string}  string    "error"
// @Failure      403         {string}  string    "error"
// @Failure      500         {string}  string    "error"
// @Router       /download-album/{albumID}/{sourceName}/ [get]
// @Security     ApiKeyAuth
func downloadAlbumHandler(c *gin.Context) {
	source, err := sources.New(c.Param("sourceName"), c.Query("api_key"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dir, err := source.DownloadAlbum(c.Param("albumID"), c.Query("dir"))
	if err != nil {
		var e *sources.AccessError
		if errors.As(err, &e) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"dir": dir, "error": ""})
}

// downloadAllAlbumsHandler godoc
// @Summary      download photos of albums
// @Description  download all photos of all albums, returns destination of your photos
// @Produce      json
// @Accept       json
// @Param        sourceName  path     string  true  "source name"
// @Param        dir         query    string  true  "directory where photos will be stored"
// @Success      200         {array}  string
// @Failure      400         {string}  string    "error"
// @Failure      401         {string}  string    "error"
// @Failure      403         {string}  string    "error"
// @Failure      500         {string}  string    "error"
// @Router       /download-all-albums/{sourceName}/ [get]
// @Security     ApiKeyAuth
func downloadAllAlbumsHandler(c *gin.Context) {
	source, err := sources.New(c.Param("sourceName"), c.Query("api_key"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dir, err := source.DownloadAllAlbums(c.Query("dir"))
	if err != nil {
		var e *sources.AccessError
		if errors.As(err, &e) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"dir": dir, "error": ""})
}
