package hgin

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

// StaticZipFS
// src should be []byte or filename
func StaticZipFS(engine *gin.Engine, relativePath string, src any, hName func(name string) string) {
	var (
		r   *zip.Reader
		err error
	)

	switch v := src.(type) {
	case string:
		if !strings.HasSuffix(v, ".zip") {
			panic("string src only support filename")
		}

		rc, err := zip.OpenReader(v)
		if err != nil {
			panic(err)
		}
		r = &rc.Reader
	case []byte:
		r, err = zip.NewReader(bytes.NewReader(v), int64(len(v)))
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("unsupported src type: %v", src))
	}

	zfs := &zipFileSystem{
		reader: r,
		hName:  hName,
	}
	engine.StaticFS(relativePath, http.FS(zfs))
}

type zipFileSystem struct {
	reader *zip.Reader
	hName  func(name string) string
}

func (z *zipFileSystem) Open(name string) (fs.File, error) {
	slog.Info("zipfs", slog.Any("raw_path", name))
	if z.hName != nil {
		name = z.hName(name)
	}

	f, err := z.reader.Open(strings.TrimPrefix(name, "/"))
	if err != nil {
		slog.Error("zipfs", slog.Any("path", name), slog.Any("err", err.Error()))
		return nil, os.ErrNotExist
	} else {
		slog.Info("zipfs", slog.Any("path", name))
	}
	return f, nil
}
