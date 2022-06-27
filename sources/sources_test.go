package sources

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type PhotoItem struct {
	url       string
	albumName string
	exifInfo  ExifInfo
	err       error
}

func (p *PhotoItem) Url() string {
	return p.url
}
func (p *PhotoItem) AlbumName() string {
	return p.albumName
}
func (p *PhotoItem) ExifInfo() (ExifInfo, error) {
	return p.exifInfo, p.err
}

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

func (s *StorageTest) SetExif(filepath string, data ExifInfo) error {
	return s.setExifErr
}

type SourceTest struct {
	albums []map[string]string
	err    error
}

func (source *SourceTest) AllAlbums() ([]map[string]string, error) {
	return source.albums, source.err
}
func (source *SourceTest) AlbumPhotos(albumdID string, photo chan Photo) error {
	return source.err
}

type service struct{}

func (s *service) Kind() Kind {
	return KindSource
}

func (s *service) Key() string {
	return "test"
}

func (s *service) Constructor() func(creds string) Source {
	return func(creds string) Source {
		return &SourceTest{}
	}
}

type storage struct{}

func (s *storage) Kind() Kind {
	return KindStorage
}

func (s *storage) Key() string {
	return "test"
}

func (s *storage) Constructor() func() Storage {
	return func() Storage {
		return &StorageTest{}
	}
}

func TestNew(t *testing.T) {
	sourceTest := &SourceTest{}
	storageTest := &StorageTest{}
	AddSource(&service{})
	AddStorage(&storage{})
	type args struct {
		sourceName string
		creds      string
		storage    Storage
	}
	tests := []struct {
		name    string
		args    args
		want    *Social
		wantErr bool
	}{
		{
			name: "no error",
			args: args{
				sourceName: "test",
				creds:      "secrets",
				storage:    storageTest,
			},
			want: &Social{
				source:  sourceTest,
				storage: storageTest,
			},
			wantErr: false,
		},
		{
			name: "has error",
			args: args{
				sourceName: "nonExistent",
				creds:      "secrets",
				storage:    storageTest,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.sourceName, tt.args.creds)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNew_NoStorage(t *testing.T) {
	storageTest := &StorageTest{}
	registeredStorages = map[string]func() Storage{}
	AddSource(&service{})
	type args struct {
		sourceName string
		creds      string
		storage    Storage
	}
	tests := []struct {
		name    string
		args    args
		want    *Social
		wantErr bool
	}{
		{
			name: "error",
			args: args{
				sourceName: "test",
				creds:      "secrets",
				storage:    storageTest,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.sourceName, tt.args.creds)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSources(t *testing.T) {
	registeredSources = map[string]func(string) Source{}
	AddSource(&service{})
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "1 item",
			want: []string{"test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sources()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSocial_DownloadAlbum(t *testing.T) {
	sourceTest := &SourceTest{}
	type fields struct {
		name    string
		creds   string
		source  Source
		storage Storage
	}
	type args struct {
		albumID string
		dest    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "no error",
			fields: fields{
				name:    "test",
				creds:   "secret",
				source:  sourceTest,
				storage: &StorageTest{dir: "dir", err: nil},
			},
			args: args{
				albumID: "123",
				dest:    "/tmp/photoD/album1",
			},
			want:    "dir",
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				name:    "test",
				creds:   "secret",
				source:  sourceTest,
				storage: &StorageTest{dir: "dir", err: errors.New("error")},
			},
			args: args{
				albumID: "123",
				dest:    "/tmp/photoD/album1",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Social{
				source:  tt.fields.source,
				storage: tt.fields.storage,
			}
			got, err := s.DownloadAlbum(tt.args.albumID, tt.args.dest)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSocial_DownloadAllAlbums(t *testing.T) {
	albums := []map[string]string{
		{
			"id": "1",
		},
	}
	sourceTest := &SourceTest{}
	type fields struct {
		source  Source
		storage Storage
	}
	type args struct {
		dest string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "no error",
			fields: fields{
				source:  sourceTest,
				storage: &StorageTest{dir: "dir", err: nil},
			},
			args:    args{dest: "/tmp/photoD"},
			want:    "dir",
			wantErr: false,
		},
		{
			name: "source error",
			fields: fields{
				source:  &SourceTest{err: errors.New("error")},
				storage: &StorageTest{dir: "dir", err: nil},
			},
			args:    args{dest: "/tmp/photoD"},
			want:    "",
			wantErr: true,
		},
		{
			name: "storage error",
			fields: fields{
				source:  &SourceTest{},
				storage: &StorageTest{err: errors.New("error")},
			},
			args:    args{dest: "/tmp/photoD"},
			want:    "",
			wantErr: true,
		},
		{
			name: "download no error",
			fields: fields{
				source:  &SourceTest{albums: albums},
				storage: &StorageTest{},
			},
			args:    args{dest: "/tmp/photoD"},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Social{
				source:  tt.fields.source,
				storage: tt.fields.storage,
			}
			got, err := s.DownloadAllAlbums(tt.args.dest)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSocial_Albums(t *testing.T) {
	albums := []map[string]string{
		{
			"id": "1",
		},
	}
	// sourceTest := &SourceTest{}
	type fields struct {
		name    string
		creds   string
		source  Source
		storage Storage
	}
	tests := []struct {
		name    string
		fields  fields
		want    []map[string]string
		wantErr bool
	}{
		{
			name: "no error",
			fields: fields{
				name:    "test",
				creds:   "secret",
				source:  &SourceTest{albums: albums},
				storage: &StorageTest{},
			},
			want:    albums,
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				name:    "test",
				creds:   "secret",
				source:  &SourceTest{err: errors.New("error")},
				storage: &StorageTest{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Social{
				source:  tt.fields.source,
				storage: tt.fields.storage,
			}
			got, err := s.Albums()
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSocial_savePhotos(t *testing.T) {
	type fields struct {
		source  Source
		storage Storage
	}
	type args struct {
		photoCh chan Photo
		exifErr error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "no error",
			fields: fields{
				source:  &SourceTest{},
				storage: &StorageTest{albumdir: "asd", downloadPhoto: "/tmp/photoD/asd.jpg"},
			},
		},
		{
			name: "album error",
			fields: fields{
				source:  &SourceTest{},
				storage: &StorageTest{createalbumdirErr: errors.New("something goes wrong")},
			},
		},
		{
			name: "download error",
			fields: fields{
				source:  &SourceTest{},
				storage: &StorageTest{downloadPhotoErr: errors.New("something goes wrong")},
			},
		},
		{
			name: "exif error",
			fields: fields{
				source:  &SourceTest{},
				storage: &StorageTest{albumdir: "asd", downloadPhoto: "/tmp/photoD/asd.jpg"},
			},
			args: args{exifErr: errors.New("something goes wrong")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.photoCh = make(chan Photo, maxConcurrentFiles)
			s := &Social{
				source:  tt.fields.source,
				storage: tt.fields.storage,
			}
			tt.args.photoCh <- &PhotoItem{albumName: "album1", url: "https://example.com/asd.jpg", err: tt.args.exifErr}
			go func() {
				time.Sleep(1 * time.Second)
				close(tt.args.photoCh)
			}()
			s.savePhotos(tt.args.photoCh)
		})
	}
}

func TestStorageError_Error(t *testing.T) {
	type fields struct {
		text string
		err  error
	}
	tests := []struct {
		name string
		err  error
		// want    string
		// wantErr string
	}{
		{
			name: "Storage error",
			err: &StorageError{
				text: "test",
				err:  errors.New("original error"),
			},
		},
		{
			name: "Source error",
			err: &SourceError{
				text: "test",
				err:  errors.New("original error"),
			},
		},
		{
			name: "Acess error",
			err: &AccessError{
				Text: "test",
				Err:  errors.New("original error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := tt.err.Error()
			errors.Unwrap(tt.err)
			assert.NotEmpty(t, text)
			//tt.err.Unwrap()
		})
	}
}
