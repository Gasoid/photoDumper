package vk

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"

	exif "github.com/Gasoid/simpleGoExif"
)

const (
	maxCount        = 1000
	concurrentFiles = 5
	VK              = "vk"
)

var (
	fileChannel chan DownloadFile
)

// `Vk` is a struct that contains a string, a pointer to a `PhotosGetAlbumsResponse` struct, an int,
// and a pointer to a `VK` struct.
// @property {string} token - VK API token
// @property Albums - a list of all albums in the user's account.
// @property {int} CurAlbum - the index of the current album in the Albums array
// @property vkAPI - the main API object, which is used to make requests to the VK API.
type Vk struct {
	token    string
	Albums   *api.PhotosGetAlbumsResponse
	CurAlbum int
	vkAPI    *api.VK
}

// DownloadFile is a struct that contains a directory, a URL, a creation time, an album name, and a
// longitude and latitude.
type DownloadFile struct {
	dir       string
	url       string
	created   time.Time
	albumName string
	longitude,
	latitude float64
}

// It takes a URL, parses it, and returns the base name of the path
func (f *DownloadFile) filePath() (string, error) {
	name, err := FileName(f.url)
	if err != nil {
		log.Println("filePath()", err)
		return "", err
	}
	return filepath.Join(f.dir, name), nil
}

// It's setting EXIF data for the downloaded file.
func (f *DownloadFile) setExifInfo() {
	filepath, err := f.filePath()
	if err != nil {
		log.Println("filePath()", err)
		return
	}
	image, err := exif.Open(filepath)
	if err != nil {
		log.Println("exif.Open", err)
	}
	defer image.Close()
	// Description
	description := fmt.Sprintf("Dumped by photoDumper. Source is vk. Album name: %s", f.albumName)
	image.SetDescription(description)
	image.SetTime(f.created)
	image.SetGPS(f.latitude, f.longitude)
}

type Albums interface {
	Add(name, cover string)
}

// It creates a new Vk object, which is a wrapper around the vkAPI object
func New(creds string) interface{} {
	if fileChannel == nil {
		fileChannel = make(chan DownloadFile, concurrentFiles)
		go downloadFile()
	}
	return &Vk{token: creds, vkAPI: api.NewVK(creds)}
}

// Getting albums from vk api
func (v *Vk) GetAlbums() ([]map[string]string, error) {
	resp, err := v.vkAPI.PhotosGetAlbums(api.Params{"need_covers": 1})
	if err != nil {
		return nil, fmt.Errorf("GetAlbums error: %w", err)
	}
	albums := make([]map[string]string, resp.Count)
	for i, album := range resp.Items {
		if album.ID < 0 {
			continue
		}
		created := time.Unix(int64(album.Created), 0)
		albums[i] = map[string]string{
			"thumb":   album.ThumbSrc,
			"title":   album.Title,
			"id":      fmt.Sprint(album.ID),
			"created": created.Format(time.RFC3339),
			"size":    fmt.Sprint(album.Size),
			// "count": album.,
		}
	}
	return albums, nil
}

// Getting photos from the album.
func (v *Vk) GetAlbumPhotos(albumId string) ([]map[string]string, error) {
	resp, err := v.vkAPI.PhotosGet(api.Params{"album_id": albumId, "count": maxCount, "photo_sizes": 1})
	if err != nil {
		return nil, fmt.Errorf("GetAlbumPhotos error: %w", err)
	}
	photos := make([]map[string]string, resp.Count)
	for i, photo := range resp.Items {
		photos[i] = map[string]string{"thumb": photo.Sizes[0].URL, "title": photo.Title, "id": fmt.Sprint(photo.ID)}
	}

	return photos, nil
}

// Downloading photos from a VK album.
func (v *Vk) DownloadAlbum(albumID, dir string) error {
	params := api.Params{"album_ids": albumID}
	if strings.Contains(albumID, "-") {
		params["need_system"] = 1
	}
	albumResp, err := v.vkAPI.PhotosGetAlbums(params)
	if err != nil {
		return fmt.Errorf("DownloadAlbum: %w", err)
	}
	// log.Println(albumID)
	resp, err := v.vkAPI.PhotosGet(api.Params{"album_id": albumID, "count": maxCount, "photo_sizes": 1})
	if err != nil {
		log.Println("DownloadAlbum:", err)
		return fmt.Errorf("DownloadAlbum: %w", err)
	}
	if albumResp.Count < 1 {
		return errors.New("no such an album")
	}
	if albumResp.Items[0].Title == "" {
		return errors.New("album title is empty")
	}
	albumDir := filepath.Join(dir, albumResp.Items[0].Title)
	_, err = os.Stat(albumDir)
	if err != nil {
		err = os.Mkdir(albumDir, 0750)
		if err != nil {
			log.Println("DownloadAlbum:", err)
			return fmt.Errorf("DownloadAlbum: %w", err)
		}
	}

	for _, photo := range resp.Items {
		var url string
		if photo.MaxSize().URL == "" {
			for _, s := range photo.Sizes {
				if s.Type == "x" {
					url = s.URL
				}
			}
		} else {
			url = photo.MaxSize().URL
		}

		created := time.Unix(int64(photo.Date), 0)
		fileChannel <- DownloadFile{
			dir:       albumDir,
			url:       url,
			created:   created,
			albumName: albumResp.Items[0].Title,
			latitude:  photo.Lat,
			longitude: photo.Long,
		}
	}

	return nil
}

// Downloading all albums from the user's account.
func (v *Vk) DownloadAllAlbums(dir string) error {
	resp, err := v.vkAPI.PhotosGetAlbums(api.Params{"need_covers": 1, "need_system": 1})
	if err != nil {
		return fmt.Errorf("GetAlbums error: %w", err)
	}
	for _, album := range resp.Items {
		go func(albumID string) {
			v.DownloadAlbum(albumID, dir)
		}(fmt.Sprint(album.ID))
	}

	return nil
}

// It takes a URL, parses it, and returns the base name of the path
func FileName(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return filepath.Base(u.Path), nil
}

// It downloads the file from the url, creates a file with the name of the file, and writes the body of
// the response to the file
func downloadFile() {
	for file := range fileChannel {
		f := file
		go func() {
			// Get the data
			resp, err := http.Get(f.url)
			if err != nil {
				log.Println(err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				log.Printf("%q is unavailable. code is 404", f.url)
				return
			}
			defer resp.Body.Close()
			filepath, err := f.filePath()
			if err != nil {
				log.Println(err)
				return
			}
			// Create the file
			out, err := os.Create(filepath)
			if err != nil {
				log.Println(err)
				return
			}

			// Write the body to file
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				log.Println(err)
				return
			}
			out.Close()
			f.setExifInfo()
		}()
	}
	log.Println("channel closed")
}
