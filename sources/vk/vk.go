package vk

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
)

const (
	maxCount = 1000
)

type Vk struct {
	token    string
	Albums   *api.PhotosGetAlbumsResponse
	CurAlbum int
	vkAPI    *api.VK
}

type Albums interface {
	Add(name, cover string)
}

func New(creds string) *Vk {
	return &Vk{token: creds, vkAPI: api.NewVK(creds)}
}

func (v *Vk) GetAlbums() ([]map[string]string, error) {
	resp, err := v.vkAPI.PhotosGetAlbums(api.Params{"need_covers": 1, "need_system": 1})
	if err != nil {
		return nil, fmt.Errorf("GetAlbums error: %w", err)
	}
	albums := make([]map[string]string, resp.Count)
	for i, album := range resp.Items {
		albums[i] = map[string]string{"thumb": album.ThumbSrc, "title": album.Title, "id": fmt.Sprint(album.ID)}
	}
	return albums, nil
}

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
	albumDir := path.Join(dir, albumResp.Items[0].Title)
	_, err = os.Stat(albumDir)
	if err != nil {
		err = os.Mkdir(albumDir, 0750)
		if err != nil {
			log.Println("DownloadAlbum:", err)
			return fmt.Errorf("DownloadAlbum: %w", err)
		}
	}

	for _, photo := range resp.Items {
		if photo.MaxSize().URL == "" {
			log.Printf("DownloadFile: photo.MaxSize().URL (%q, album is %q) is empty", photo.ID, albumResp.Items[0].Title)
			continue
		}
		err := DownloadFile(albumDir, photo.MaxSize().URL)
		if err != nil {
			log.Println("DownloadFile:", err)
		}
	}

	return nil
}

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

func FileName(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return filepath.Base(u.Path), nil
}

func DownloadFile(dir string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	name, err := FileName(url)
	if err != nil {
		return err
	}
	filepath := path.Join(dir, name)
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
