package sources

import (
	"fmt"
	"log"
	"time"
)

var (
	registeredSources  = map[string]func(string) Source{}
	photoCh            chan Photo
	maxConcurrentFiles = 5
)

type StorageError struct {
	text string
	err  error
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("Source error: %s", e.text)
}

func (e *StorageError) Unwrap() error {
	return e.err
}

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

type AccessError struct {
	Text string
	Err  error
}

func (e *AccessError) Error() string {
	return fmt.Sprintf("Auth error: %s", e.Text)
}

func (e *AccessError) Unwrap() error {
	return e.Err
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
		return nil, err
	}
	return albums, nil
}

func (s *Social) DownloadAllAlbums() (string, error) {
	dir, err := s.storage.Prepare()
	if err != nil {
		log.Println("DownloadAllAlbums(dir string)", err)
		return "", &StorageError{text: "dir can't be created", err: err}
	}

	albums, err := s.source.AllAlbums()
	if err != nil {
		return "", err
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
		return "", &StorageError{text: "dir can't be created", err: err}
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
	if sourceNew, ok := registeredSources[sourceName]; ok {
		s.source = sourceNew(creds)
	} else {
		return nil, &SourceError{text: "there is no such a source"}
	}
	if photoCh == nil {
		photoCh = make(chan Photo, maxConcurrentFiles)
		go s.storage.SavePhotos(photoCh)
	}
	return s, nil
}

func Sources() []string {
	listSources := make([]string, len(registeredSources))
	var i int
	for key := range registeredSources {
		listSources[i] = key
		i++
	}
	return listSources
}

func AddSource(sourceName string, newFunc func(string) Source) {
	registeredSources[sourceName] = newFunc
}
