package server

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sasthw/internal/archive"
	"strings"

	log "github.com/google/logger"
	"github.com/labstack/echo"
)

func (s *Server) treeHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		log.Error(err)
		return echo.ErrInternalServerError
	}

	src, err := file.Open()
	if err != nil {
		log.Error(err)
		return echo.ErrInternalServerError
	}
	defer func() { _ = src.Close() }()

	tmpDir, err := ioutil.TempDir("", "tree-*")
	if err != nil {
		log.Error(err)
		return echo.ErrInternalServerError
	}

	archiveDir := path.Join(tmpDir, "./data/")
	if err := archive.Unpack(src, path.Join(tmpDir, "./data/")); err != nil {
		log.Error(err)
		return echo.ErrInternalServerError
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	result := strings.Builder{}
	err = filepath.Walk(archiveDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			result.WriteString(fmt.Sprintf("%s %d B\n", html.EscapeString(path), info.Size()))
			return nil
		})
	if err != nil {
		log.Error(err)
		return echo.ErrInternalServerError
	}

	return c.String(http.StatusOK, result.String())
}
