package vk

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
)

const (
	maxCount = 1000
	VK       = "vk"
)

type AccessError struct {
	text string
	err  error
}

func (e *AccessError) Error() string {
	return fmt.Sprintf("Auth error %s", e.text)
}

func (e *AccessError) Unwrap() error {
	return e.err
}

// `Vk` is a struct that contains a string, a pointer to a `PhotosGetAlbumsResponse` struct, an int,
// and a pointer to a `VK` struct.
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

func (f *DownloadFile) GetUrl() string {
	return f.url
}

func (f *DownloadFile) GetAlbumName() string {
	return f.albumName
}

func (f *DownloadFile) GetFilename() string {
	u, err := url.Parse(f.url)
	if err != nil {
		return ""
	}
	return filepath.Base(u.Path)
}

// It's setting EXIF data for the downloaded file.
func (f *DownloadFile) GetExifInfo() (map[string]interface{}, error) {
	exifInfo := map[string]interface{}{
		"description": fmt.Sprintf("Dumped by photoDumper. Source is vk. Album name: %s", f.albumName),
		"created":     f.created,
		"gps":         []float64{f.latitude, f.longitude},
	}

	return exifInfo, nil
}

// It creates a new Vk object, which is a wrapper around the vkAPI object
func New(creds string) interface{} {
	return &Vk{token: creds, vkAPI: api.NewVK(creds)}
}

// Getting albums from vk api
func (v *Vk) GetAlbums() ([]map[string]string, error) {
	resp, err := v.vkAPI.PhotosGetAlbums(api.Params{"need_covers": 1})
	if err != nil {
		return nil, makeError(err, "GetAlbums failed")
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

// Downloading photos from a VK album.
func (v *Vk) AlbumPhotos(albumID string, photoCh chan interface{}) error {
	params := api.Params{"album_ids": albumID}
	if strings.Contains(albumID, "-") {
		params["need_system"] = 1
	}
	albumResp, err := v.vkAPI.PhotosGetAlbums(params)
	if err != nil {
		return makeError(err, "DownloadAlbum failed")
	}
	// log.Println(albumID)
	resp, err := v.vkAPI.PhotosGet(api.Params{"album_id": albumID, "count": maxCount, "photo_sizes": 1})
	if err != nil {
		log.Println("DownloadAlbum:", err)
		return makeError(err, "DownloadAlbum failed")
	}
	if albumResp.Count < 1 {
		return errors.New("no such an album")
	}
	if albumResp.Items[0].Title == "" {
		return errors.New("album title is empty")
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
		photoCh <- &DownloadFile{
			url:       url,
			created:   created,
			albumName: albumResp.Items[0].Title,
			latitude:  photo.Lat,
			longitude: photo.Long,
		}
	}

	return nil
}

func makeError(err error, text string) error {
	if errors.Is(err, api.ErrSignature) || errors.Is(err, api.ErrAccess) || errors.Is(err, api.ErrAuth) {
		return &AccessError{text: text, err: err}
	}
	return fmt.Errorf("%s: %w", text, err)
}
