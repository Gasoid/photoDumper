package sources

import (
	"fmt"
	"log"
	"time"
)

type Kind int

const (
	KindSource Kind = iota
	KindStorage
)

var (
	registeredSources  = map[string]func(creds string) Source{}
	registeredStorages = map[string]func() Storage{}
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
	AlbumName() string
	ExifInfo() (ExifInfo, error)
}

type Storage interface {
	Prepare(dir string) (string, error)
	CreateAlbumDir(dir string) (string, error)
	DownloadPhoto(photoUrl, dir string) (string, error)
	SetExif(filepath string, info ExifInfo) error
}

type Social struct {
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

func (s *Social) DownloadAllAlbums(dir string) (string, error) {
	dir, err := s.storage.Prepare(dir)
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
			_, err := s.DownloadAlbum(albumID, dir)
			if err != nil {
				log.Println(err, "DownloadAllAlbums failed")
			}
		}(album["id"])
	}

	return dir, nil
}

// DownloadAlbum runs copying process to a particular directory
func (s *Social) DownloadAlbum(albumID, dir string) (string, error) {
	dir, err := s.storage.Prepare(dir)
	if err != nil {
		log.Println("DownloadAlbum(albumID, dir string)", err)
		return "", &StorageError{text: "dir can't be created", err: err}
	}
	s.source.AlbumPhotos(albumID, photoCh)
	return dir, nil
}

func (s *Social) savePhotos(photoCh chan Photo) {
	for file := range photoCh {
		f := file
		go func() {
			dir, err := s.storage.CreateAlbumDir(f.AlbumName())
			if err != nil {
				log.Println(err)
				return
			}
			filepath, err := s.storage.DownloadPhoto(f.Url(), dir)
			if err != nil {
				log.Println(err)
				return
			}
			exif, err := f.ExifInfo()
			if err != nil {
				log.Println(err)
				return
			}
			s.storage.SetExif(filepath, exif)
		}()
	}
	log.Println("channel closed")
}

// New creates a new instance of Social, you have to provide proper options
func New(sourceName, creds string) (*Social, error) {
	source, err := ProvideSource(sourceName, creds)
	if err != nil {
		return nil, err
	}
	storage, err := ProvideStorage()
	if err != nil {
		return nil, err
	}
	s := &Social{
		storage: storage,
		source:  source,
	}
	if photoCh == nil {
		photoCh = make(chan Photo, maxConcurrentFiles)
		go s.savePhotos(photoCh)
	}
	return s, nil
}

type Service interface {
	Key() string
}

type ServiceSource interface {
	Service
	Constructor() func(creds string) Source
}

type ServiceStorage interface {
	Service
	Constructor() func() Storage
}

func AddSource(s ServiceSource) {
	registeredSources[s.Key()] = s.Constructor()
}

func AddStorage(s ServiceStorage) {
	registeredStorages[s.Key()] = s.Constructor()
}

func ProvideSource(key string, creds string) (Source, error) {
	if newFunc, ok := registeredSources[key]; ok {
		return newFunc(creds), nil
	} else {
		return nil, &SourceError{text: "Source was not found"}
	}
}

func ProvideStorage() (Storage, error) {
	for _, v := range registeredStorages {
		return v(), nil
	}
	return nil, &StorageError{text: "no storages"}
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
