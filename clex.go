package clex

import (
	"io"
	"net/http"
	"os"

	"golang.org/x/xerrors"
)

type FileSystem struct {
	http.FileSystem
}

func (fs FileSystem) Copy(dst, src string) error {
	srcf, err := fs.Open(src)
	if err != nil {
		return xerrors.Errorf("failed to open src file: %w", err)
	}
	defer srcf.Close()

	dstf, err := os.Create(dst)
	if err != nil {
		return xerrors.Errorf("failed to create dst file: %w", err)
	}
	defer dstf.Close()

	if _, err := io.Copy(dstf, srcf); err != nil {
		return xerrors.Errorf("failed to copy file contents: %w", err)
	}
	return nil
}

func (fs FileSystem) Readdir(path string, count int) ([]os.FileInfo, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, xerrors.Errorf("failed to open path: %w", err)
	}
	info, err := f.Readdir(count)
	if err != nil {
		return nil, xerrors.Errorf("failed to get directory info: %w", err)
	}
	return info, nil
}

func (fs FileSystem) Stats(path string) (os.FileInfo, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, xerrors.Errorf("failed to open path: %w", err)
	}
	info, err := f.Stat()
	if err != nil {
		return nil, xerrors.Errorf("failed to get directory info: %w", err)
	}
	return info, nil
}
