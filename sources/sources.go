package sources

import (
	"errors"

	"github.com/Gasoid/photoDumper/sources/vk"
)

const (
	//Instagram string = "instagram"
	VK string = "vk"
)

var (
	sources = map[string]bool{
		//Instagram: true,
		VK: true,
	}
	ErrSourceNotFound = errors.New("there is no such a source")
)

type Source interface {
	GetAlbums() ([]map[string]string, error)
	GetAlbumPhotos(albumId string) ([]map[string]string, error)
	DownloadAllAlbums(dir string) error
	DownloadAlbum(albumdID, dir string) error
}

type Social struct {
	name   string
	creds  string
	source Source
}

func (s *Social) GetAlbums() ([]map[string]string, error) {
	return s.source.GetAlbums()
}

func (s *Social) DownloadAllAlbums(dir string) error {
	return s.source.DownloadAllAlbums(dir)
}

// DownloadAlbum runs copying process to a particular directory
func (s *Social) DownloadAlbum(albumID, dir string) error {
	return s.source.DownloadAlbum(albumID, dir)
}

type Photo struct {
	id  string
	url string
}

func (s *Social) GetAlbumPhotos(albumID string) ([]map[string]string, error) {
	return s.source.GetAlbumPhotos(albumID)
}

type Options struct {
	Source string
}

// New creates a new instance of Social, you have to provide proper options
func New(sourceName, creds string) (*Social, error) {
	if _, ok := sources[sourceName]; !ok {
		return nil, ErrSourceNotFound
	}
	s := &Social{name: sourceName, creds: creds}
	switch sourceName {
	case "vk":
		s.source = vk.New(creds)
	}
	return s, nil
}

func Sources() []string {
	return []string{VK}
}
