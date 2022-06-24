package localfs

import (
	"testing"

	// . "github.com/Gasoid/photoDumper/sources"
	"github.com/stretchr/testify/assert"
)

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SimpleStorage{
				Dir: tt.fields.Dir,
			}
			got, err := s.dirPath(tt.args.dir)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
