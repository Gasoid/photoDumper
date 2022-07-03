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
			s := &SimpleStorage{}
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
		{
			name:    "error",
			fields:  fields{Dir: ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SimpleStorage{}
			got, err := s.Prepare(tt.fields.Dir)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSimpleStorage_FilePath(t *testing.T) {

	type args struct {
		dir      string
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no error",
			args: args{dir: "/tmp", filename: "photo.jpg"},
			want: "/tmp/photo.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SimpleStorage{}
			got := s.FilePath(tt.args.dir, tt.args.filename)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSimpleStorage_createAlbumDir(t *testing.T) {
	type args struct {
		albumName string
		rootDir   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "no error",
			args:    args{albumName: "album1", rootDir: "/tmp/photoD"},
			want:    "/tmp/photoD/album1",
			wantErr: false,
		},
		// {
		// 	name:    "error",
		// 	fields:  fields{Dir: ""},
		// 	args:    args{albumName: "brrr"},
		// 	want:    "",
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SimpleStorage{}
			got, err := s.CreateAlbumDir(tt.args.rootDir, tt.args.albumName)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSimpleStorage_DownloadPhoto(t *testing.T) {
	type args struct {
		url,
		albumName string
	}
	tests := []struct {
		name string
		args args
		//photoItem PhotoItem
	}{
		{
			name: "empty f",
			args: args{},
		},
		{
			name: "empty url",
			args: args{url: ""},
		},
		{
			name: "empty exif",
			args: args{url: "https://picsum.photos/200/300.jpg", albumName: "300"},
		},
		{
			name: "404",
			args: args{url: "https://github.com/Gasoid/photoDumper/sdf"},
		},
		{
			name: "dir is empty",
			args: args{url: "https://picsum.photos/200/300.jpg", albumName: "300"},
		},
		{
			name: "exif",
			args: args{
				url:       "https://picsum.photos/200/300.jpg",
				albumName: "300",
			},
		},
		{
			name: "gps is nil",
			args: args{url: "https://picsum.photos/200/300.jpg",
				albumName: "300",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SimpleStorage{}
			s.DownloadPhoto(tt.args.url, tt.args.albumName)
		})
	}
}

func TestSimpleStorage_SetExif(t *testing.T) {
	type args struct {
		filepath  string
		photoExif sources.ExifInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "gps is nil",
			args:    args{filepath: "/tmp/photoD/300.jpg", photoExif: &ExifInfo{}},
			wantErr: true,
		},
		{
			name:    "gps exists",
			args:    args{filepath: "/tmp/photoD/300.jpg", photoExif: &ExifInfo{gps: []float64{45.4545, 45.4545}}},
			wantErr: false,
		},
		{
			name:    "wrong path",
			args:    args{filepath: "/tmp/photoD/301.jpg"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SimpleStorage{}
			s.DownloadPhoto("https://picsum.photos/200/300.jpg", "/tmp/photoD/")
			err := s.SetExif(tt.args.filepath, tt.args.photoExif)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func Test_filename(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "no error",
			args:    args{path: "https://example.com/asd.jpg"},
			want:    "asd.jpg",
			wantErr: false,
		},
		{
			name:    "error",
			args:    args{path: "https://example.com/asd"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty",
			args:    args{path: ":/sdfsdf/sdfsdf"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filename(tt.args.path)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "not nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			assert.NotNil(t, got)
		})
	}
}

func Test_service_Kind(t *testing.T) {
	tests := []struct {
		name string
		s    *service
		want sources.Kind
	}{
		{
			name: "kind",
			s:    &service{},
			want: sources.KindStorage,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Kind()
			assert.Equal(t, tt.want, got)
			k := tt.s.Key()
			assert.Equal(t, "fs", k)
		})
	}
}

func Test_service_Constructor(t *testing.T) {
	tests := []struct {
		name string
		s    *service
		want func() sources.Storage
	}{
		{
			name: "new",
			s:    &service{},
			want: New,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{}
			got := s.Constructor()
			assert.NotEmpty(t, got)
		})
	}
}

func TestNewService(t *testing.T) {
	tests := []struct {
		name string
		want sources.ServiceStorage
	}{
		{
			name: "service",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = NewService()
		})
	}
}
