package sftpfilesystem

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/djedjethai/celeritas/filesystems"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTP struct {
	Host string
	User string
	Pass string
	Port string
}

func (s *SFTP) getCredentials() (*sftp.Client, error) {
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Pass),
		},
		// allow us to not get error key ..., make it easier
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}

	// set the current dir on the server
	cwd, err := client.Getwd()
	log.Println("current working directory: ", cwd)

	return client, nil
}

func (s *SFTP) Put(filename, folder string) error {
	client, err := s.getCredentials()
	if err != nil {
		return err
	}
	defer client.Close()

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Println("the path.Base(): ", path.Base(filename))
	f2, err := client.Create(path.Base(filename))
	if err != nil {
		return err
	}
	defer f2.Close()

	// destination id f2
	if _, err := io.Copy(f2, f); err != nil {
		return err
	}

	return nil
}

func (s *SFTP) List(prefix string) ([]filesystems.Listing, error) {
	var listing []filesystems.Listing
	client, err := s.getCredentials()
	if err != nil {
		return listing, err
	}
	defer client.Close()

	files, err := client.ReadDir(prefix)
	if err != nil {
		return listing, err
	}

	for _, x := range files {
		var item filesystems.Listing

		// we don't want any file startting with a .(like .env)
		// it s not call key in this package but Name()
		if !strings.HasPrefix(x.Name(), ".") {
			b := float64(x.Size())
			kb := b / 1024
			mb := kb / 1024
			item.Key = x.Name()
			item.Size = mb
			item.LastModified = x.ModTime()
			item.IsDir = x.IsDir()
			listing = append(listing, item)
		}
	}

	return listing, nil
}

func (s *SFTP) Delete(itemsToDelete []string) bool {
	client, err := s.getCredentials()
	if err != nil {
		return false
	}
	defer client.Close()

	for _, x := range itemsToDelete {
		deleteErr := client.Remove(x)
		if deleteErr != nil {
			return false
		}
	}

	return true
}

func (s *SFTP) Get(destination string, items ...string) error {
	client, err := s.getCredentials()
	if err != nil {
		return err
	}
	defer client.Close()

	for _, item := range items {
		// create a destination here on the app fs,
		// where we gonna save the file coming from sftp server
		dstFile, err := os.Create(fmt.Sprintf("%s/%s", destination, path.Base(item)))
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// open source file
		srcFile, err := client.Open(item)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// copy src file to dest
		bytes, err := io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
		log.Println(fmt.Sprintf("%d bytes copied: ", bytes))

		// flush the in-memory copy
		// that is specific to sftp
		err = dstFile.Sync()
		if err != nil {
			return err
		}

	}

	return nil
}
