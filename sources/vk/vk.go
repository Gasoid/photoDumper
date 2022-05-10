package vk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"io/ioutil"

	"github.com/Gasoid/go-dms/dms"
	"github.com/SevereCloud/vksdk/v2/api"

	exif "github.com/dsoprea/go-exif/v2"
	exifcommon "github.com/dsoprea/go-exif/v2/common"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure"
)

const (
	maxCount        = 1000
	concurrentFiles = 5
)

var (
	fileChannel chan DownloadFile
)

type Vk struct {
	token    string
	Albums   *api.PhotosGetAlbumsResponse
	CurAlbum int
	vkAPI    *api.VK
}

type DownloadFile struct {
	dir       string
	url       string
	created   time.Time
	albumName string
	longitude,
	latitude float64
}

func (f *DownloadFile) filePath() (string, error) {
	name, err := FileName(f.url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return path.Join(f.dir, name), nil
}

// EXIF HELL
func (f *DownloadFile) setExifInfo() {
	filepath, err := f.filePath()
	if err != nil {
		log.Println("filePath()", err)
		return
	}
	jmp := jpegstructure.NewJpegMediaParser()
	intfc, err := jmp.ParseFile(filepath)
	if err != nil {
		log.Println("ParseFile(filepath)", err)
		return
	}
	sl := intfc.(*jpegstructure.SegmentList)
	rootIb, err := sl.ConstructExifBuilder()
	if err != nil {
		im := exif.NewIfdMappingWithStandard()
		ti := exif.NewTagIndex()
		err := exif.LoadStandardTags(ti)
		if err != nil {
			log.Println("ConstructExifBuilder()", err)
			return
		}

		rootIb = exif.NewIfdBuilder(im, ti, exifcommon.IfdPathStandard, exifcommon.EncodeDefaultByteOrder)
	}

	ifd0Ib, err := exif.GetOrCreateIbFromRootIb(rootIb, "IFD0")
	if err != nil {
		log.Println("GetOrCreateIbFromRootIb(rootIb, ifd0Path)", err)
		return
	}

	// Description
	description := fmt.Sprintf("Dumped by photoDumper. Source is vk. Album name: %s", f.albumName)
	err = ifd0Ib.SetStandardWithName("ImageDescription", description)
	if err != nil {
		log.Println("SetStandardWithName(ImageDescription)", err)
		return
	}

	dateTime := exif.ExifFullTimestampString(f.created)
	err = ifd0Ib.SetStandardWithName("DateTime", dateTime)
	if err != nil {
		log.Println("SetStandardWithName(DateTime)", err)
		return
	}

	if f.latitude != 0 && f.longitude != 0 {
		//log.Println("There are GPS coordinates:", f.latitude, f.longitude, filepath)

		childIb, err := exif.GetOrCreateIbFromRootIb(rootIb, "IFD/GPSInfo")
		if err != nil {
			log.Println("GetOrCreateIbFromRootIb(rootIbf.latitude, GPSInfo)", err)
			return
		}
		lat, lon, err := dms.NewDMS(f.latitude, f.longitude)
		if err != nil {
			log.Println("dms.NewDMS(f.latitude, f.longitude)", err)
			return
		}
		updatedGiLat := exif.GpsDegrees{
			Degrees: float64(lat.Degrees),
			Minutes: float64(lat.Minutes),
			Seconds: lat.Seconds,
		}

		err = childIb.SetStandardWithName("GPSLatitude", updatedGiLat.Raw())
		if err != nil {
			log.Println("SetStandardWithName(GPS)", err)
			return
		}
		updatedGiLong := exif.GpsDegrees{
			Degrees: float64(lon.Degrees),
			Minutes: float64(lon.Minutes),
			Seconds: lon.Seconds,
		}

		err = childIb.SetStandardWithName("GPSLongitude", updatedGiLong.Raw())
		if err != nil {
			log.Println("SetStandardWithName(GPS)", err)
			return
		}
	}

	err = sl.SetExif(rootIb)
	if err != nil {
		log.Println("SetExif()", err)
		return
	}
	b := bytes.NewBufferString("")
	err = sl.Write(b)
	if err != nil {
		log.Println("Write(b)", err)
		return
	}
	ioutil.WriteFile(filepath, b.Bytes(), 0666)
}

type Albums interface {
	Add(name, cover string)
}

func New(creds string) *Vk {
	if fileChannel == nil {
		fileChannel = make(chan DownloadFile, concurrentFiles)
		go downloadFile()
	}
	return &Vk{token: creds, vkAPI: api.NewVK(creds)}
}

func (v *Vk) GetAlbums() ([]map[string]string, error) {
	resp, err := v.vkAPI.PhotosGetAlbums(api.Params{"need_covers": 1, "need_system": 1})
	if err != nil {
		return nil, fmt.Errorf("GetAlbums error: %w", err)
	}
	albums := make([]map[string]string, resp.Count)
	for i, album := range resp.Items {
		created := time.Unix(int64(album.Created), 0)
		albums[i] = map[string]string{
			"thumb":   album.ThumbSrc,
			"title":   album.Title,
			"id":      fmt.Sprint(album.ID),
			"created": created.Format(time.RFC3339),
			// "count": album.,
		}
	}
	return albums, nil
}

func (v *Vk) GetAlbumPhotos(albumId string) ([]map[string]string, error) {
	resp, err := v.vkAPI.PhotosGet(api.Params{"album_id": albumId, "count": maxCount, "photo_sizes": 1})
	if err != nil {
		return nil, fmt.Errorf("GetAlbumPhotos error: %w", err)
	}
	photos := make([]map[string]string, resp.Count)
	for i, photo := range resp.Items {
		photos[i] = map[string]string{"thumb": photo.Sizes[0].URL, "title": photo.Title, "id": fmt.Sprint(photo.ID)}
	}

	return photos, nil
}

func (v *Vk) DownloadAlbum(albumID, dir string) error {
	params := api.Params{"album_ids": albumID}
	if strings.Contains(albumID, "-") {
		params["need_system"] = 1
	}
	albumResp, err := v.vkAPI.PhotosGetAlbums(params)
	if err != nil {
		return fmt.Errorf("DownloadAlbum: %w", err)
	}
	// log.Println(albumID)
	resp, err := v.vkAPI.PhotosGet(api.Params{"album_id": albumID, "count": maxCount, "photo_sizes": 1})
	if err != nil {
		log.Println("DownloadAlbum:", err)
		return fmt.Errorf("DownloadAlbum: %w", err)
	}
	if albumResp.Count < 1 {
		return errors.New("no such an album")
	}
	if albumResp.Items[0].Title == "" {
		return errors.New("album title is empty")
	}
	albumDir := path.Join(dir, albumResp.Items[0].Title)
	_, err = os.Stat(albumDir)
	if err != nil {
		err = os.Mkdir(albumDir, 0750)
		if err != nil {
			log.Println("DownloadAlbum:", err)
			return fmt.Errorf("DownloadAlbum: %w", err)
		}
	}

	for _, photo := range resp.Items {
		var url string
		if photo.MaxSize().URL == "" {
			for _, s := range photo.Sizes {
				if s.Type == "x" {
					url = s.URL
				}
			}
		} else {
			url = photo.MaxSize().URL
		}

		created := time.Unix(int64(photo.Date), 0)
		fileChannel <- DownloadFile{
			dir:       albumDir,
			url:       url,
			created:   created,
			albumName: albumResp.Items[0].Title,
			latitude:  photo.Lat,
			longitude: photo.Long,
		}
	}

	return nil
}

func (v *Vk) DownloadAllAlbums(dir string) error {
	resp, err := v.vkAPI.PhotosGetAlbums(api.Params{"need_covers": 1, "need_system": 1})
	if err != nil {
		return fmt.Errorf("GetAlbums error: %w", err)
	}
	for _, album := range resp.Items {
		go func(albumID string) {
			v.DownloadAlbum(albumID, dir)
		}(fmt.Sprint(album.ID))
	}

	return nil
}

func FileName(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return filepath.Base(u.Path), nil
}

func downloadFile() {
	for file := range fileChannel {
		f := file
		go func() {
			// Get the data
			resp, err := http.Get(f.url)
			if err != nil {
				log.Println(err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				log.Printf("%q is unavailable. code is 404", f.url)
				return
			}
			defer resp.Body.Close()
			filepath, err := f.filePath()
			if err != nil {
				log.Println(err)
				return
			}
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
			f.setExifInfo()
		}()
	}
	log.Println("channel closed")
}
