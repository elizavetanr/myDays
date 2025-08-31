package storage

import (
	"archive/zip"
	"errors"
	"io"
	"os"
)

var (
	ErrStorageEmpty = errors.New("архив пуст")
)

type ZipStorage struct {
	*Storage
}

func NewZipStorage(filename string) *ZipStorage {
	return &ZipStorage{
		&Storage{filename: filename},
	}
}

func (s *ZipStorage) Save(data []byte) error {
	f, err := os.Create(s.GetFileName())
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	w, err := zw.Create("data")
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (s *ZipStorage) Load() ([]byte, error) {
	r, err := zip.OpenReader(s.GetFileName())
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if len(r.File) == 0 {
		return nil, ErrStorageEmpty
	}

	file := r.File[0]
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return io.ReadAll(rc)
}
