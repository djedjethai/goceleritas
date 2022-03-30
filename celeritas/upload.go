package celeritas

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/djedjethai/celeritas/fileSystems"
	"github.com/gabriel-vasile/mimetype"
)

func (c *Celeritas) UploadFile(r *http.Request, destination, field string, fs fileSystems.FS) error {
	fileName, err := c.getFileToUpload(r, field)
	if err != nil {
		c.ErrorLog.Println(err)
		return err
	}

	// fs fileSystems.FS is a pointer(but is an interface so i don't write it)
	// if we are using a remote fs we pass a pointer to it
	// or nil if we are uploading to a local fs
	if fs != nil {
		err = fs.Put(fileName, destination)
		if err != nil {
			c.ErrorLog.Println(err)
			return err
		}
	} else {
		err = os.Rename(fileName, fmt.Sprintf("%s/%s", destination, path.Base(fileName)))
		if err != nil {
			c.ErrorLog.Println(err)
			return err
		}

	}

	return nil
}

func (c *Celeritas) getFileToUpload(r *http.Request, fieldName string) (string, error) {
	_ = r.ParseMultipartForm(c.config.uploads.maxUploadSize)

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// mimetype.DetectReader() read the 512bytes of the file to detect the mimetype
	// so just after i need to reset the reader pointer to the beginning of the file
	mimeType, err := mimetype.DetectReader(file)
	if err != nil {
		return "", err
	}

	// reset the reader pointer to the beginning of the file
	// otherwise user won't be able to open the file(what ever it is, pic or doc)
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	// see any kind of mimetype i can use
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_type
	if !inSlice(c.config.uploads.allowedMimeTypes, mimeType.String()) {
		return "", errors.New("invalid file type uploaded")
	}

	dst, err := os.Create(fmt.Sprintf("./tmp/%s", header.Filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("./tmp/%s", header.Filename), nil
}

func inSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
