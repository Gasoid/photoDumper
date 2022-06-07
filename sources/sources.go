package sources

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

var (
	RegisteredSources = map[string]func(string) Source{}
	ErrSourceNotFound = errors.New("there is no such a source")
)

type Source interface {
	GetAlbums() ([]map[string]string, error)
	GetAlbumPhotos(albumId string) ([]map[string]string, error)
	DownloadAllAlbums(dir string) error
	DownloadAlbum(albumdID, dir string) error
	IsAuthError(err error) bool
}

type Social struct {
	name   string
	creds  string
	source Source
}

// GetAlbums returns albums
func (s *Social) GetAlbums() ([]map[string]string, error) {
	return s.source.GetAlbums()
}

// GetAlbums returns albums
func (s *Social) IsAuthError(err error) bool {
	return s.source.IsAuthError(err)
}

// It's a method of Social struct. It's checking if the path is absolute or relative.
func (s *Social) dirPath(dir string) (string, error) {
	if dir[:1] == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Println("filePath()", err)
			return "", err
		}
		dir = filepath.Join(home, filepath.FromSlash(dir[1:]))
	}
	return dir, nil
}

// Creating a directory if it doesn't exist.
func (s *Social) prepareDir(dir string) (string, error) {
	dir, err := s.dirPath(dir)
	if err != nil {
		log.Println("DownloadAllAlbums(dir string)", err)
		return "", err
	}
	err = os.MkdirAll(dir, 0750)
	if err != nil {
		log.Println("DownloadAllAlbums(dir string)", err)
	}
	return dir, err
}

func (s *Social) DownloadAllAlbums(dir string) error {
	dir, err := s.prepareDir(dir)
	if err != nil {
		log.Println("DownloadAllAlbums(dir string)", err)
		return err
	}
	return s.source.DownloadAllAlbums(dir)
}

// DownloadAlbum runs copying process to a particular directory
func (s *Social) DownloadAlbum(albumID, dir string) error {
	dir, err := s.prepareDir(dir)
	if err != nil {
		log.Println("DownloadAlbum(albumID, dir string)", err)
		return err
	}
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
	s := &Social{name: sourceName, creds: creds}
	if sourceNew, ok := RegisteredSources[sourceName]; ok {
		s.source = sourceNew(creds)
	} else {
		return nil, ErrSourceNotFound
	}
	return s, nil
}

func Sources() []string {
	listSources := make([]string, len(RegisteredSources))
	var i int
	for key := range RegisteredSources {
		listSources[i] = key
		i++
	}
	return listSources
}

func AddSource(sourceName string, newFunc func(string) interface{}) {
	RegisteredSources[sourceName] = func(creds string) Source {
		result := newFunc(creds)
		return result.(Source)
	}
}
