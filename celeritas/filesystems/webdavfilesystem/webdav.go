package webdavfilesystem

import (
	"github.com/djedjethai/celeritas/filesystems"
)

type WebDAV struct {
	Host string
	User string
	Path string
}

func (s *WebDAV) Put(filename, folder string) error {
	return nil
}

func (s *WebDAV) List(prefix string) ([]filesystems.Listing, error) {
	var list []filesystems.Listing

	return list, nil
}

func (s *WebDAV) Delete(itemsToDelete []string) bool {
	return true
}

func (s *WebDAV) Get(destination string, items ...string) error {
	return nil
}
