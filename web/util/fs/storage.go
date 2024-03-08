package fs

import "github.com/stellarisJAY/nesgo/web/config"

type FileStorage interface {
	Store(path string, data []byte) error
	Load(path string) ([]byte, error)
	Type() string
	Delete(path string) error
}

func NewFileStorage(storeType string) (FileStorage, error) {
	switch storeType {
	case "host":
		return NewHostFileSystemStorage(config.GetConfig().HostFileSystemStorageDir)
	default:
		panic("unimplemented store type")
	}
}
