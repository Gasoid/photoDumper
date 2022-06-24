package instagram

// import (
// 	"fmt"
// 	"net/url"
// 	"path/filepath"
// 	"time"

// 	"github.com/Gasoid/photoDumper/sources"
// )

// const (
// 	maxCount = 1000
// 	ID       = "vk"
// )

// type IG struct {
// 	token string
// 	api   *ig.Client
// }

// // PhotoItem is a struct that contains a directory, a URL, a creation time, an album name, and a
// // longitude and latitude.
// type PhotoItem struct {
// 	url       string
// 	created   time.Time
// 	albumName string
// 	longitude,
// 	latitude float64
// }

// func (f *PhotoItem) Url() string {
// 	return f.url
// }

// func (f *PhotoItem) AlbumName() string {
// 	return f.albumName
// }

// func (f *PhotoItem) Filename() string {
// 	u, err := url.Parse(f.url)
// 	if err != nil {
// 		return ""
// 	}
// 	return filepath.Base(u.Path)
// }

// // It's setting EXIF data for the downloaded file.
// func (f *PhotoItem) ExifInfo() (sources.ExifInfo, error) {
// 	exif := &exifInfo{
// 		description: fmt.Sprintf("Dumped by photoDumper. Source is vk. Album name: %s", f.albumName),
// 		created:     f.created,
// 		gps:         []float64{f.latitude, f.longitude},
// 	}
// 	return exif, nil
// }

// type exifInfo struct {
// 	description string
// 	created     time.Time
// 	gps         []float64
// }

// func (e *exifInfo) Description() string {
// 	return e.description
// }

// func (e *exifInfo) Created() time.Time {
// 	return e.created
// }

// func (e *exifInfo) GPS() []float64 {
// 	return e.gps
// }

// // It creates a new Vk object, which is a wrapper around the vkAPI object
// func New(creds string) sources.Source {
// 	return &IG{token: creds, api: ig.NewClient(creds)}
// }

// // Getting albums from vk api
// func (ic *IG) AllAlbums() ([]map[string]string, error) {

// 	return nil, nil
// }

// // Downloading photos from Instagram.
// func (ic *IG) AlbumPhotos(albumID string, photoCh chan sources.Photo) error {
// 	return nil
// }

// // func makeError(err error, text string) error {
// // 	if errors.Is(err, api.ErrSignature) || errors.Is(err, api.ErrAccess) || errors.Is(err, api.ErrAuth) {
// // 		return &sources.AccessError{Text: text, Err: err}
// // 	}
// // 	return fmt.Errorf("%s: %w", text, err)
// // }
