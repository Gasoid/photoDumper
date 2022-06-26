package localfs

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Gasoid/photoDumper/sources"
	exif "github.com/Gasoid/simpleGoExif"
)

type SimpleStorage struct {
	Dir string
}

// It's a method of Social struct. It's checking if the path is absolute or relative.
func (s *SimpleStorage) dirPath(dir string) (string, error) {
	if len(dir) < 1 {
		return "", fmt.Errorf("len of dir is less 1")
	}
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

func (s *SimpleStorage) Prepare() (string, error) {
	dir, err := s.dirPath(s.Dir)
	if err != nil {
		log.Println("prepareDir", err)
		return "", err
	}
	err = os.MkdirAll(dir, 0750)
	if err != nil {
		log.Println("prepareDir", err)
	}
	return dir, err
}

// It takes a URL, parses it, and returns the base name of the path
func (s *SimpleStorage) FilePath(dir, filename string) string {
	return filepath.Join(dir, filename)
}

func (s *SimpleStorage) createAlbumDir(albumName string) (string, error) {
	dir, err := s.dirPath(s.Dir)
	if err != nil {
		log.Println("createAlbumDir", err)
		return "", err
	}
	albumDir := filepath.Join(dir, albumName)
	err = os.MkdirAll(albumDir, 0750)
	if err != nil {
		log.Println("createAlbumDir:", err)
		return "", fmt.Errorf("createAlbumDir: %w", err)
	}
	return albumDir, nil
}

// It downloads the file from the url, creates a file with the name of the file, and writes the body of
// the response to the file
func (s *SimpleStorage) SavePhoto(f sources.Photo) {
	if f == nil {
		return
	}
	resp, err := http.Get(f.Url())
	if err != nil {
		log.Println(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("%q is unavailable. code is %d", f.Url(), resp.StatusCode)
		return
	}
	defer resp.Body.Close()
	dir, err := s.createAlbumDir(f.AlbumName())
	if err != nil {
		log.Println(err)
		return
	}
	filepath := s.FilePath(dir, f.Filename())
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
		return
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	out.Close()
	photoExif, err := f.ExifInfo()
	if err != nil {
		log.Println(err)
		return
	}
	s.setExifInfo(filepath, photoExif)
}

// It's setting EXIF data for the downloaded file.
func (s *SimpleStorage) setExifInfo(filepath string, photoExif sources.ExifInfo) error {
	image, err := exif.Open(filepath)
	if err != nil {
		log.Println("exif.Open", err)
		return err
	}
	defer image.Close()
	if photoExif == nil {
		return errors.New("exif is empty")
	}
	err = image.SetDescription(photoExif.Description())
	if err != nil {
		log.Println("image.SetDescription", err)
		return err
	}
	err = image.SetTime(photoExif.Created())
	if err != nil {
		log.Println("image.SetTime", err)
		return err
	}
	gps := photoExif.GPS()
	if gps == nil {
		return errors.New("gps is empty")
	}
	err = image.SetGPS(gps[0], gps[1])
	if err != nil {
		log.Println("image.SetGPS", err)
		return err
	}

	return nil
}

func New() sources.Storage {
	return &SimpleStorage{}
}

type service struct{}

func (s *service) Kind() sources.Kind {
	return sources.KindStorage
}

func (s *service) Key() string {
	return "fs"
}

func (s *service) Constructor() func() sources.Storage {
	return New
}

func NewService() sources.ServiceStorage {
	return &service{}
}
