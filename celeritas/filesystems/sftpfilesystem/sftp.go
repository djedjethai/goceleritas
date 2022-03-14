package sftpfilesystem

import (
	"github.com/djedjethai/celeritas/filesystems"
)

type SFTP struct {
	Host string
	User string
	Path string
	Port string
}

func (s *SFTP) Put(filename, folder string) error {
	return nil
}

func (s *SFTP) List(prefix string) ([]filesystems.Listing, error) {
	var list []filesystems.Listing

	return list, nil
}

func (s *SFTP) Delete(itemsToDelete []string) bool {
	return true
}

func (s *SFTP) Get(destination string, items ...string) error {
	return nil
}
