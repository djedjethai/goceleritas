package s3filesystem

import (
	"github.com/djedjethai/celeritas/filesystems"
)

type S3 struct {
	Ket      string
	Secret   string
	Region   string
	Endpoint string
	Bucket   string
}

func (s *S3) Put(filename, folder string) error {
	return nil
}

func (s *S3) List(prefix string) ([]filesystems.Listing, error) {
	var list []filesystems.Listing

	return list, nil
}

func (s *S3) Delete(itemsToDelete []string) bool {
	return true
}

func (s *S3) Get(destination string, items ...string) error {
	return nil
}
