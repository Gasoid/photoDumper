package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
func (s *SimpleStorage) SavePhotos(photoCh chan sources.Photo) {
	for file := range photoCh {
		f := file
		go func() {
			// Get the data
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
			s.setExifInfo(filepath, f)
		}()
	}
	log.Println("channel closed")
}

// It's setting EXIF data for the downloaded file.
func (s *SimpleStorage) setExifInfo(filepath string, photo sources.Photo) error {
	image, err := exif.Open(filepath)
	if err != nil {
		log.Println("exif.Open", err)
		return err
	}
	defer image.Close()
	// Description
	exifInfo, err := photo.ExifInfo()
	if err != nil {
		log.Println("photo.ExifInfo()", err)
		return err
	}
	if description, ok := exifInfo["description"].(string); ok {
		err = image.SetDescription(description)
		if err != nil {
			log.Println("image.SetDescription", err)
			return err
		}
	}
	if created, ok := exifInfo["created"].(time.Time); ok {
		err = image.SetTime(created)
		if err != nil {
			log.Println("image.SetTime", err)
			return err
		}
	}
	if gps, ok := exifInfo["gps"].([]float64); ok {
		err = image.SetGPS(gps[0], gps[1])
		if err != nil {
			log.Println("image.SetGPS", err)
			return err
		}
	}

	return nil
}
