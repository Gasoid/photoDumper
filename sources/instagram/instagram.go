package instagram

import (
	"fmt"
	"time"

	"github.com/Gasoid/photoDumper/sources"
)

type PhotoItem struct {
	url       string
	albumName string
	created   time.Time
}

func (f *PhotoItem) Url() string {
	return f.url
}

func (f *PhotoItem) AlbumName() string {
	return f.albumName
}

// It's setting EXIF data for the downloaded file.
func (f *PhotoItem) ExifInfo() (sources.ExifInfo, error) {
	exif := &exifInfo{
		description: fmt.Sprintf("Dumped by photoDumper. Source is vk. Username: %s", f.albumName),
		created:     f.created,
	}
	return exif, nil
}

type exifInfo struct {
	description string
	created     time.Time
}

func (e *exifInfo) Description() string {
	return e.description
}

func (e *exifInfo) Created() time.Time {
	return e.created
}

func (e *exifInfo) GPS() []float64 {
	return nil
}

type service struct{}

func (s *service) Kind() sources.Kind {
	return sources.KindSource
}

func (s *service) Key() string {
	return "instagram"
}

func (s *service) Constructor() func(creds string) sources.Source {
	return New
}

func NewService() sources.ServiceSource {
	return &service{}
}

func New(creds string) sources.Source {
	api := &InstagramApi{access_token: creds}
	return &Instagram{api: api}
}

type Instagram struct {
	api *InstagramApi
}

func (ig *Instagram) AllAlbums() ([]map[string]string, error) {
	resp := ig.api.Me("id", "username", "media_count")
	media, err := ig.api.MeMedia("id", "media_url", "timestamp", "caption")
	if err != nil {
		return nil, &sources.AccessError{Err: err, Text: "token is invalid?"}
	}
	album := map[string]string{}
	for media.Next() {
		if media.Item().MediaType != IMAGE_TYPE {
			continue
		}
		album = map[string]string{
			"thumb":   media.Item().MediaUrl,
			"title":   "All Instagram photos and videos",
			"id":      "all_photos_and_videos",
			"created": media.Item().Timestamp,
			"size":    fmt.Sprint(resp.MediaCount),
		}
		break
	}

	albums := make([]map[string]string, 1)
	albums[0] = album
	return albums, nil
}

type fetcher struct {
	media *PagingResponse
}

func (f *fetcher) Next() bool {
	return f.media.Next()
}

func (f *fetcher) Item() sources.Photo {
	photo := f.media.Item()
	date, err := time.Parse("2006-01-02T15:04:05-0700", photo.Timestamp)
	if err != nil {
		date = time.Now()
	}
	return &PhotoItem{
		url:       photo.MediaUrl,
		albumName: photo.Username,
		created:   date,
		// latitude:  photo.Lat,
		// longitude: photo.Long,
	}
}

func (ig *Instagram) AlbumPhotos(albumID string) (sources.ItemFetcher, error) {
	media, err := ig.api.MeMedia("id", "media_url", "timestamp", "caption", "username")
	if err != nil {
		return nil, &sources.AccessError{Err: err, Text: "token is invalid?"}
	}
	return &fetcher{media: media}, nil
}
