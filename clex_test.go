package clex

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFileSystem_Copy(t *testing.T) {
	type fields struct {
		FileSystem http.FileSystem
	}
	type args struct {
		dst string
		src string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "正常系",
			fields: fields{FileSystem: http.Dir("testdata")},
			args: args{
				dst: "test",
				src: "test",
			},
		},
		{
			name:   "notfound",
			fields: fields{FileSystem: http.Dir("testdata")},
			args: args{
				dst: "",
				src: "notfound",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := FileSystem{
				FileSystem: tt.fields.FileSystem,
			}
			if err := fs.Copy(tt.args.dst, tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("FileSystem.Copy() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil {
				return
			}
			defer os.Remove(tt.args.dst)

			var src bytes.Buffer
			f, err := fs.Open(tt.args.src)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			if _, err := io.Copy(&src, f); err != nil {
				t.Fatal(err)
			}

			var dst bytes.Buffer
			f, err = os.Open(tt.args.dst)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			if _, err := io.Copy(&dst, f); err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(dst.Bytes(), src.Bytes()) {
				t.Fatal("copied file contents are not same")
			}
		})
	}
}

func TestFileSystem_Readdir(t *testing.T) {
	type fields struct {
		FileSystem http.FileSystem
	}
	type args struct {
		path  string
		count int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantLen int
		wantErr bool
	}{
		{
			name:   "正常系",
			fields: fields{FileSystem: http.Dir("testdata")},
			args: args{
				path:  ".",
				count: 0,
			},
			wantLen: 1,
		},
		{
			name:   "notfound",
			fields: fields{FileSystem: http.Dir("testdata")},
			args: args{
				path:  "notfound",
				count: 0,
			},
			wantErr: true,
		},
		{
			name:   "not directory",
			fields: fields{FileSystem: http.Dir("testdata")},
			args: args{
				path:  "test",
				count: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := FileSystem{
				FileSystem: tt.fields.FileSystem,
			}
			got, err := fs.Readdir(tt.args.path, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileSystem.Readdir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(len(got), tt.wantLen) {
				t.Errorf("len(FileSystem.Readdir()) = %v, wantLen %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestFileSystem_Stats(t *testing.T) {
	type fields struct {
		FileSystem http.FileSystem
	}
	type args struct {
		path string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantName string
		wantErr  bool
	}{
		{
			name:     "正常系",
			fields:   fields{FileSystem: http.Dir("testdata")},
			args:     args{path: "test"},
			wantName: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := FileSystem{
				FileSystem: tt.fields.FileSystem,
			}
			got, err := fs.Stats(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileSystem.Stats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got.Name(), tt.wantName) {
				t.Errorf("FileSystem.Stats() = %v, want %v", got, tt.wantName)
			}
		})
	}
}
