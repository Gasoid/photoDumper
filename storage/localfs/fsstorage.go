package localfs

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Gasoid/photoDumper/sources"
	exif "github.com/Gasoid/simpleGoExif"
)

type SimpleStorage struct {
	dir string
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

func (s *SimpleStorage) Prepare(dir string) (string, error) {
	dir, err := s.dirPath(dir)
	if err != nil {
		log.Println("prepareDir", err)
		return "", err
	}
	err = os.MkdirAll(dir, 0750)
	if err != nil {
		log.Println("prepareDir", err)
	}
	s.dir = dir
	return dir, err
}

// It takes a URL, parses it, and returns the base name of the path
func (s *SimpleStorage) FilePath(dir, filename string) string {
	return filepath.Join(dir, filename)
}

func filename(path string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	name := filepath.Base(u.Path)
	if filepath.Ext(name) == "" {
		return "", errors.New("no ext")
	}
	return name, nil
}

func (s *SimpleStorage) CreateAlbumDir(albumName string) (string, error) {
	albumDir := filepath.Join(s.dir, albumName)
	err := os.MkdirAll(albumDir, 0750)
	if err != nil {
		log.Println("createAlbumDir:", err)
		return "", fmt.Errorf("createAlbumDir: %w", err)
	}
	return albumDir, nil
}

// It downloads the file from the url, creates a file with the name of the file, and writes the body of
// the response to the file
func (s *SimpleStorage) DownloadPhoto(url, dir string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("%q is unavailable. code is %d", url, resp.StatusCode)
		return "", err
	}
	defer resp.Body.Close()
	name, err := filename(url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	filepath := s.FilePath(dir, name)
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	out.Close()
	return filepath, nil
}

// It's setting EXIF data for the downloaded file.
func (s *SimpleStorage) SetExif(filepath string, photoExif sources.ExifInfo) error {
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
