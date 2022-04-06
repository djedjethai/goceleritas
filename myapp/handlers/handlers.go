package handlers

import (
	"fmt"
	"io"
	"myapp/data"
	"net/http"
	"net/url"
	"os"

	"github.com/CloudyKit/jet/v6"
	"github.com/djedjethai/celeritas"
	"github.com/djedjethai/celeritas/filesystems"
	"github.com/djedjethai/celeritas/filesystems/miniofilesystem"
	"github.com/djedjethai/celeritas/filesystems/s3filesystem"
	"github.com/djedjethai/celeritas/filesystems/sftpfilesystem"
	"github.com/djedjethai/celeritas/filesystems/webdavfilesystem"
)

// Handlers is the type for handlers, and gives access to Celeritas and models
type Handlers struct {
	App    *celeritas.Celeritas
	Models data.Models
}

// Home is the handler to render the home page
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "home", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

// render the form form page for the user to upload a file
func (h *Handlers) CeleritasUpload(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "celeritas-upload", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

// upload a file
func (h *Handlers) PostCeleritasUpload(w http.ResponseWriter, r *http.Request) {

	// like it it have export to the file-system
	// err := h.App.UploadFile(r, "", "formFile", "")

	// like this it will upload to minio
	err := h.App.UploadFile(r, "", "formFile", &h.App.Minio)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
		h.App.Session.Put(r.Context(), "error", err.Error())
	} else {
		h.App.Session.Put(r.Context(), "flash", "uploaded!")
	}

	http.Redirect(w, r, "/upload", http.StatusSeeOther)

}

func inSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func (h *Handlers) ListFS(w http.ResponseWriter, r *http.Request) {
	var fs filesystems.FS
	var list []filesystems.Listing

	fsType := ""
	if r.URL.Query().Get("fs-type") != "" {
		fsType = r.URL.Query().Get("fs-type")
	}

	curPath := "/"
	if r.URL.Query().Get("curPath") != "" {
		curPath = r.URL.Query().Get("curPath")
		// unescape our query data(take out special char)
		curPath, _ = url.QueryUnescape(curPath)
	}

	if fsType != "" {
		switch fsType {
		case "MINIO":
			f := h.App.FileSystems["MINIO"].(miniofilesystem.Minio)
			fs = &f
			fsType = "MINIO"
		case "SFTP":
			f := h.App.FileSystems["SFTP"].(sftpfilesystem.SFTP)
			fs = &f
			fsType = "SFTP"
		case "WEBDAV":
			f := h.App.FileSystems["WEBDAV"].(webdavfilesystem.WebDAV)
			fs = &f
			fsType = "WEBDAV"
		case "S3":
			f := h.App.FileSystems["S3"].(s3filesystem.S3)
			fs = &f
			fsType = "S3"

		}

		l, err := fs.List(curPath)
		if err != nil {
			h.App.ErrorLog.Println(err)
			return
		}

		list = l
	}

	vars := make(jet.VarMap)
	vars.Set("list", list)
	vars.Set("fs_type", fsType)
	vars.Set("curPath", curPath)
	err := h.render(w, r, "list-fs", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) UploadToFS(w http.ResponseWriter, r *http.Request) {
	fsType := r.URL.Query().Get("type")

	vars := make(jet.VarMap)
	vars.Set("fs_type", fsType)

	err := h.render(w, r, "upload", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) PostUploadToFS(w http.ResponseWriter, r *http.Request) {

	fileName, err := getFileToUpload(r, "formFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uploadType := r.Form.Get("upload-type")

	switch uploadType {
	case "MINIO":
		fs := h.App.FileSystems["MINIO"].(miniofilesystem.Minio)
		err = fs.Put(fileName, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "SFTP":
		fs := h.App.FileSystems["SFTP"].(sftpfilesystem.SFTP)
		err = fs.Put(fileName, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "WEBDAV":
		fs := h.App.FileSystems["WEBDAV"].(webdavfilesystem.WebDAV)
		err = fs.Put(fileName, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "S3":
		fs := h.App.FileSystems["S3"].(s3filesystem.S3)
		err = fs.Put(fileName, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	h.App.Session.Put(r.Context(), "flash", "File uploaded!")
	http.Redirect(w, r, "/files/upload?type="+uploadType, http.StatusSeeOther)

}

func getFileToUpload(r *http.Request, fieldName string) (string, error) {
	// 10 << 20 ends to max 10megs??? of data
	_ = r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

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

func (h *Handlers) DeleteFromFS(w http.ResponseWriter, r *http.Request) {
	var fs filesystems.FS
	fsType := r.URL.Query().Get("fs_type")
	file := r.URL.Query().Get("file")

	switch fsType {
	case "MINIO":
		f := h.App.FileSystems["MINIO"].(miniofilesystem.Minio)
		fs = &f
	case "SFTP":
		f := h.App.FileSystems["SFTP"].(sftpfilesystem.SFTP)
		fs = &f
	case "WEBDAV":
		f := h.App.FileSystems["WEBDAV"].(webdavfilesystem.WebDAV)
		fs = &f
	case "S3":
		f := h.App.FileSystems["S3"].(s3filesystem.S3)
		fs = &f

	}

	deleted := fs.Delete([]string{file})
	if deleted {
		h.App.Session.Put(r.Context(), "flash", fmt.Sprintf("%s was deleted", file))
		http.Redirect(w, r, "/list-fs?fs-type="+fsType, http.StatusSeeOther)
	}
}
