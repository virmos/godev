package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/CloudyKit/jet/v6"
	"github.com/virmos/django/filesystems"
	"github.com/virmos/django/filesystems/miniofilesystem"
	"github.com/virmos/django/filesystems/s3filesystem"
	"github.com/virmos/django/filesystems/sftpfilesystem"
	"github.com/virmos/django/filesystems/webdavfilesystem"
)

func (h *Handlers) DjangoUpload(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "django-upload", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) PostDjangoUpload(w http.ResponseWriter, r *http.Request) {
	err := h.App.UploadFile(r, "", "formFile", &h.App.Minio)
	if err != nil {
		h.App.ErrorLog.Println(err)
		h.App.Session.Put(r.Context(), "error", err.Error())
	} else {
		h.App.Session.Put(r.Context(), "flash", "Uploaded!")
	}
	http.Redirect(w, r, "/upload", http.StatusSeeOther)
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

func (h *Handlers) DeleteFromFS(w http.ResponseWriter, r *http.Request) {
	var fs filesystems.FS
	fsType := r.URL.Query().Get("fs_type")
	item := r.URL.Query().Get("file")

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

	deleted := fs.Delete([]string{item})
	if deleted {
		h.App.Session.Put(r.Context(), "flash", fmt.Sprintf("%s was deleted", item))
		http.Redirect(w, r, "/list-fs?fs-type="+fsType, http.StatusSeeOther)
	}
}
