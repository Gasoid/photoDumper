package sources

import (
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	RegisteredSources = map[string]func(string) Source{}
	photoCh           chan Photo
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
	AllAlbums() ([]map[string]string, error)
	AlbumPhotos(albumdID string, photo chan Photo) error
}

type ExifInfo interface {
	Description() string
	Created() time.Time
	GPS() []float64
}

type Photo interface {
	Url() string
	Filename() string
	AlbumName() string
	ExifInfo() (ExifInfo, error)
}

type Storage interface {
	Prepare() (string, error)
	SavePhotos(photo chan Photo)
}

type Social struct {
	name    string
	creds   string
	source  Source
	storage Storage
}

// Albums returns albums
func (s *Social) Albums() ([]map[string]string, error) {
	albums, err := s.source.AllAlbums()
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

	albums, err := s.source.AllAlbums()
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
func New(sourceName, creds string, storage Storage) (*Social, error) {
	s := &Social{
		name:    sourceName,
		creds:   creds,
		storage: storage,
	}
	if sourceNew, ok := RegisteredSources[sourceName]; ok {
		s.source = sourceNew(creds)
	} else {
		return nil, &SourceError{text: "there is no such a source"}
	}
	if photoCh == nil {
		photoCh = make(chan Photo, concurrentFiles)
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

func AddSource(sourceName string, newFunc func(string) Source) {
	RegisteredSources[sourceName] = newFunc
}
