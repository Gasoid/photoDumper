package localfs

import (
	"testing"
	"time"

	"github.com/Gasoid/photoDumper/sources"
	"github.com/stretchr/testify/assert"
)

type ExifInfo struct {
	description string
	created     time.Time
	gps         []float64
}

func (e *ExifInfo) Description() string {
	return e.description
}
func (e *ExifInfo) Created() time.Time {
	return e.created
}

func (e *ExifInfo) GPS() []float64 {
	return e.gps
}

type PhotoItem struct {
	url       string
	filename  string
	albumName string
	exifInfo  sources.ExifInfo
	err       error
}

func (p *PhotoItem) Url() string {
	return p.url
}
func (p *PhotoItem) Filename() string {
	return p.filename
}
func (p *PhotoItem) AlbumName() string {
	return p.albumName
}
func (p *PhotoItem) ExifInfo() (sources.ExifInfo, error) {
	return p.exifInfo, p.err
}

func TestSimpleStorage_dirPath(t *testing.T) {
	type fields struct {
		Dir string
	}
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "no error",
			fields:  fields{},
			args:    args{"/albumDir"},
			want:    "/albumDir",
			wantErr: false,
		},
		{
			name:    "error empty",
			fields:  fields{},
			args:    args{""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "home dir",
			fields:  fields{},
			args:    args{"~/photoD"},
			want:    "homeDir",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SimpleStorage{
				Dir: tt.fields.Dir,
			}
			got, err := s.dirPath(tt.args.dir)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.want == "homeDir" {
				assert.NotEmpty(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}

		})
	}
}

func TestSimpleStorage_Prepare(t *testing.T) {
	type fields struct {
		Dir string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "no error",
			fields:  fields{Dir: "/tmp/photoD"},
			want:    "/tmp/photoD",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SimpleStorage{
				Dir: tt.fields.Dir,
			}
			got, err := s.Prepare()
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSimpleStorage_FilePath(t *testing.T) {
	type fields struct {
		Dir string
	}
	type args struct {
		dir      string
		filename string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "no error",
			fields: fields{Dir: "/tmp/photoD"},
			args:   args{dir: "/tmp", filename: "photo.jpg"},
			want:   "/tmp/photo.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SimpleStorage{
				Dir: tt.fields.Dir,
			}
			got := s.FilePath(tt.args.dir, tt.args.filename)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSimpleStorage_createAlbumDir(t *testing.T) {
	type fields struct {
		Dir string
	}
	type args struct {
		albumName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "no error",
			fields:  fields{Dir: "/tmp/photoD"},
			args:    args{albumName: "album1"},
			want:    "/tmp/photoD/album1",
			wantErr: false,
		},
		{
			name:    "error",
			fields:  fields{Dir: ""},
			args:    args{albumName: "brrr"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SimpleStorage{
				Dir: tt.fields.Dir,
			}
			got, err := s.createAlbumDir(tt.args.albumName)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSimpleStorage_SavePhoto(t *testing.T) {
	type fields struct {
		Dir string
	}
	type args struct {
		f sources.Photo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//photoItem PhotoItem
	}{
		{
			name:   "empty f",
			fields: fields{Dir: "/tmp/photoD"},
			args:   args{},
		},
		{
			name:   "empty url",
			fields: fields{Dir: "/tmp/photoD"},
			args:   args{f: &PhotoItem{url: ""}},
		},
		{
			name:   "empty exif",
			fields: fields{Dir: "/tmp/photoD"},
			args:   args{f: &PhotoItem{url: "https://picsum.photos/200/300.jpg", filename: "300.jpg", albumName: "300"}},
		},
		{
			name:   "404",
			fields: fields{Dir: "/tmp/photoD"},
			args:   args{f: &PhotoItem{url: "https://github.com/Gasoid/photoDumper/sdf"}},
		},
		{
			name:   "dir is empty",
			fields: fields{Dir: ""},
			args:   args{f: &PhotoItem{url: "https://picsum.photos/200/300.jpg", filename: "300.jpg", albumName: "300"}},
		},
		{
			name:   "exif",
			fields: fields{Dir: "/tmp/photoD"},
			args: args{f: &PhotoItem{
				url:       "https://picsum.photos/200/300.jpg",
				filename:  "300.jpg",
				albumName: "300",
				exifInfo:  &ExifInfo{gps: []float64{0, 0}},
			}},
		},
		{
			name:   "gps is nil",
			fields: fields{Dir: "/tmp/photoD"},
			args: args{f: &PhotoItem{
				url:       "https://picsum.photos/200/300.jpg",
				filename:  "300.jpg",
				albumName: "300",
				exifInfo:  &ExifInfo{},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SimpleStorage{
				Dir: tt.fields.Dir,
			}
			s.SavePhoto(tt.args.f)
		})
	}
}

func TestSimpleStorage_setExifInfo(t *testing.T) {
	type fields struct {
		Dir string
	}
	type args struct {
		filepath  string
		photoExif sources.ExifInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "gps is nil",
			fields:  fields{Dir: "/tmp/photoD"},
			args:    args{filepath: "/tmp/photoD/1.jpg", photoExif: &ExifInfo{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SimpleStorage{
				Dir: tt.fields.Dir,
			}
			err := s.setExifInfo(tt.args.filepath, tt.args.photoExif)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
