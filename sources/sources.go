package sources

import (
	"fmt"
	"log"
	"strings"
)

var (
	RegisteredSources = map[string]func(string) Source{}
	photoCh           chan interface{}
	concurrentFiles   = 5
)

type SourceError struct {
	text string
	err  error
}

func (e *SourceError) Error() string {
	return fmt.Sprintf("Source error: %s", e.text)
}

func (e *SourceError) Unwrap() error {
	return e.err
}

type AuthError struct {
	text string
	err  error
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("Auth error: %s", e.text)
}

func (e *AuthError) Unwrap() error {
	return e.err
}

type Source interface {
	GetAlbums() ([]map[string]string, error)
	AlbumPhotos(albumdID string, photo chan interface{}) error
}

type Photo interface {
	GetUrl() string
	GetFilename() string
	GetAlbumName() string
	GetExifInfo() (map[string]interface{}, error)
}

type Storage interface {
	Prepare() (string, error)
	SavePhotos(photo chan interface{})
}

type Social struct {
	name    string
	creds   string
	source  Source
	storage Storage
}

// GetAlbums returns albums
func (s *Social) GetAlbums() ([]map[string]string, error) {
	albums, err := s.source.GetAlbums()
	if err != nil {
		if strings.Contains(err.Error(), "Auth error") {
			return nil, &AuthError{"Albums are inaccessible", err}
		}
		return nil, &SourceError{"Albums are inaccessible", err}
	}
	return albums, nil
}

func (s *Social) DownloadAllAlbums() (string, error) {
	dir, err := s.storage.Prepare()
	if err != nil {
		log.Println("DownloadAllAlbums(dir string)", err)
		return "", &SourceError{text: "dir can't be created"}
	}

	albums, err := s.source.GetAlbums()
	if err != nil {
		if strings.Contains(err.Error(), "Auth error") {
			return "", &AuthError{"Albums are inaccessible", err}
		}
		return "", &SourceError{"Albums are inaccessible", err}
	}
	for _, album := range albums {
		go func(albumID string) {
			_, err := s.DownloadAlbum(albumID)
			if err != nil {
				log.Println(err, "DownloadAllAlbums failed")
			}
		}(album["id"])
	}

	return dir, nil
}

// DownloadAlbum runs copying process to a particular directory
func (s *Social) DownloadAlbum(albumID string) (string, error) {
	dir, err := s.storage.Prepare()
	if err != nil {
		log.Println("DownloadAlbum(albumID, dir string)", err)
		if strings.Contains(err.Error(), "Auth error") {
			return "", &AuthError{"Album is inaccessible", err}
		}
		return "", &SourceError{"Album is inaccessible", err}
	}
	s.source.AlbumPhotos(albumID, photoCh)
	return dir, nil
}

// New creates a new instance of Social, you have to provide proper options
func New(sourceName, creds string, storage interface{}) (*Social, error) {
	s := &Social{
		name:    sourceName,
		creds:   creds,
		storage: storage.(Storage),
	}
	if sourceNew, ok := RegisteredSources[sourceName]; ok {
		s.source = sourceNew(creds)
	} else {
		return nil, &SourceError{text: "there is no such a source"}
	}
	if photoCh == nil {
		photoCh = make(chan interface{}, concurrentFiles)
		go s.storage.SavePhotos(photoCh)
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
