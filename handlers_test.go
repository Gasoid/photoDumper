package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gasoid/photoDumper/sources"
	"github.com/stretchr/testify/assert"
)

type StorageTest struct {
	dir               string
	err               error
	albumdir          string
	createalbumdirErr error
	downloadPhoto     string
	downloadPhotoErr  error
	setExifErr        error
}

func (s *StorageTest) Prepare(dir string) (string, error) {
	return s.dir, s.err
}

func (s *StorageTest) DownloadPhoto(photoUrl, dir string) (string, error) {
	return s.downloadPhoto, s.downloadPhotoErr
}

func (s *StorageTest) CreateAlbumDir(dir string) (string, error) {
	return s.albumdir, s.createalbumdirErr
}

func (s *StorageTest) SetExif(filepath string, data sources.ExifInfo) error {
	return s.setExifErr
}

type SourceTest struct {
	albums []map[string]string
	err    error
}

func (source *SourceTest) AllAlbums() ([]map[string]string, error) {
	return source.albums, source.err
}
func (source *SourceTest) AlbumPhotos(albumdID string) (sources.ItemFetcher, error) {
	return &testFetcher{}, source.err
}

type testFetcher struct{}

func (tf *testFetcher) Next() bool {
	return false
}

func (tf *testFetcher) Item() sources.Photo {
	return nil
}

type service struct {
	sourceError error
}

func (s *service) Kind() sources.Kind {
	return sources.KindSource
}

func (s *service) Key() string {
	return "test"
}

func (s *service) Constructor() func(creds string) sources.Source {
	return func(creds string) sources.Source {
		return &SourceTest{err: s.sourceError}
	}
}

type storage struct {
	err error
}

func (s *storage) Kind() sources.Kind {
	return sources.KindStorage
}

func (s *storage) Key() string {
	return "test"
}

func (s *storage) Constructor() func() sources.Storage {
	return func() sources.Storage {
		return &StorageTest{err: s.err}
	}
}

func Test_sources(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/sources/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_assets(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func Test_albums(t *testing.T) {
	// sourceTest := &SourceTest{}
	// storageTest := &StorageTest{}
	sources.AddSource(&service{})
	sources.AddStorage(&storage{})
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/albums/test1/?api_key=sdfsdf", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/api/albums/test/?api_key=sdfsdf", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func Test_albumsError(t *testing.T) {
	// sourceTest := &SourceTest{}
	// storageTest := &StorageTest{}
	sources.AddSource(&service{sourceError: errors.New("bad")})
	sources.AddStorage(&storage{})
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/albums/test/?api_key=sdfsdf", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func Test_albumsAccessError(t *testing.T) {
	// sourceTest := &SourceTest{}
	// storageTest := &StorageTest{}
	sources.AddSource(&service{sourceError: &sources.AccessError{}})
	sources.AddStorage(&storage{})
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/albums/test/?api_key=sdfsdf", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func Test_downloadAlbumStorageError(t *testing.T) {
	sources.AddSource(&service{sourceError: &sources.AccessError{}})
	sources.AddStorage(&storage{err: errors.New("bad")})
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/download-album/albumid/test/?api_key=sdfsdf", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func Test_downloadAlbumError(t *testing.T) {
	sources.AddSource(&service{sourceError: &sources.AccessError{}})
	sources.AddStorage(&storage{err: errors.New("bad")})
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/download-album/albumid/test1/?api_key=sdfsdf", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_downloadAlbumAccessError(t *testing.T) {
	sources.AddSource(&service{sourceError: &sources.AccessError{}})
	sources.AddStorage(&storage{err: &sources.AccessError{}})
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/download-album/albumid/test/?api_key=sdfsdf", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func Test_downloadAlbum(t *testing.T) {
	sources.AddSource(&service{})
	sources.AddStorage(&storage{})
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/download-album/albumid/test/?api_key=sdfsdf", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_downloadAllAlbums(t *testing.T) {
	sources.AddSource(&service{})
	sources.AddStorage(&storage{})
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/download-all-albums/test/?api_key=sdfsdf", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_downloadAllAlbumsAccessError(t *testing.T) {
	sources.AddSource(&service{})
	sources.AddStorage(&storage{err: &sources.AccessError{}})
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/download-all-albums/test/?api_key=sdfsdf", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func Test_downloadAllAlbumsError(t *testing.T) {
	sources.AddSource(&service{})
	sources.AddStorage(&storage{})
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/download-all-albums/test1/?api_key=sdfsdf", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_downloadAllAlbumsBadError(t *testing.T) {
	sources.AddSource(&service{sourceError: &sources.AccessError{}})
	sources.AddStorage(&storage{err: errors.New("bad")})
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/download-all-albums/test/?api_key=sdfsdf", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func Test_downloadAllAlbumsApiKeyError(t *testing.T) {
	sources.AddSource(&service{})
	sources.AddStorage(&storage{})
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/download-all-albums/test/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
