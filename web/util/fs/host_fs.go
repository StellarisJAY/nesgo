package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

type HostFileSystemStorage struct {
	dir string
}

func NewHostFileSystemStorage(dir string) (*HostFileSystemStorage, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to make host fs dir: %v", err)
	}
	return &HostFileSystemStorage{dir: dir}, nil
}

func (h *HostFileSystemStorage) Store(path string, data []byte) error {
	path = filepath.Join(h.dir, path)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create save dir: %v", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

func (h *HostFileSystemStorage) Load(path string) ([]byte, error) {
	data, err := os.ReadFile(filepath.Join(h.dir, path))
	if err != nil {
		return nil, fmt.Errorf("read file error: %v", err)
	}
	return data, nil
}

func (h *HostFileSystemStorage) Delete(path string) error {
	return os.Remove(path)
}

func (h *HostFileSystemStorage) Type() string {
	return "host"
}
