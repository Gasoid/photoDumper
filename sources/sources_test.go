package sources

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StorageTest struct {
	dir string
	err error
}

func (s *StorageTest) Prepare() (string, error) {
	return s.dir, s.err
}

func (s *StorageTest) SavePhotos(photo chan Photo) {
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

func TestNew(t *testing.T) {
	sourceTest := &SourceTest{}
	storageTest := &StorageTest{}
	AddSource("test", func(creds string) Source { return sourceTest })
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
				name:    "test",
				creds:   "secrets",
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
			got, err := New(tt.args.sourceName, tt.args.creds, tt.args.storage)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSources(t *testing.T) {
	sourceTest := &SourceTest{}
	registeredSources = map[string]func(string) Source{}
	AddSource("test", func(creds string) Source { return sourceTest })
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
				storage: &StorageTest{"dir", nil},
			},
			args: args{
				albumID: "123",
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
				storage: &StorageTest{"dir", errors.New("error")},
			},
			args: args{
				albumID: "123",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Social{
				name:    tt.fields.name,
				creds:   tt.fields.creds,
				source:  tt.fields.source,
				storage: tt.fields.storage,
			}
			got, err := s.DownloadAlbum(tt.args.albumID)
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
		name    string
		creds   string
		source  Source
		storage Storage
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "no error",
			fields: fields{
				name:    "test",
				creds:   "secret",
				source:  sourceTest,
				storage: &StorageTest{"dir", nil},
			},
			want:    "dir",
			wantErr: false,
		},
		{
			name: "source error",
			fields: fields{
				name:    "test",
				creds:   "secret",
				source:  &SourceTest{err: errors.New("error")},
				storage: &StorageTest{"dir", nil},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "storage error",
			fields: fields{
				name:    "test",
				creds:   "secret",
				source:  &SourceTest{},
				storage: &StorageTest{err: errors.New("error")},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "download no error",
			fields: fields{
				name:    "test",
				creds:   "secret",
				source:  &SourceTest{albums: albums},
				storage: &StorageTest{},
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Social{
				name:    tt.fields.name,
				creds:   tt.fields.creds,
				source:  tt.fields.source,
				storage: tt.fields.storage,
			}
			got, err := s.DownloadAllAlbums()
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
				name:    tt.fields.name,
				creds:   tt.fields.creds,
				source:  tt.fields.source,
				storage: tt.fields.storage,
			}
			got, err := s.Albums()
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
