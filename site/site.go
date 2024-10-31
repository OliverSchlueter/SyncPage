package site

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"syncpage/github"
	"time"
)

const (
	UpdateInterval = time.Second * 30
)

var (
	NoAssetError = errors.New("no asset found for release")
)

type Site struct {
	Name string
	Repo github.Repository
}

func (s *Site) Register(mux *http.ServeMux) {
	go s.startUpdateLoop()

	mux.HandleFunc("/"+s.Name+"/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len(s.Name)+1:]
		if path == "" {
			path = "index.html"
		}

		content, err := os.ReadFile("sites/" + s.Name + "/" + path)
		if err != nil {
			if os.IsNotExist(err) {
				s.TryToReturnIndex(w, r)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(content)
	})
}

func (s *Site) TryToReturnIndex(w http.ResponseWriter, r *http.Request) {
	path := "sites/" + s.Name + "/index.html"

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(content)
}

func (s *Site) startUpdateLoop() {
	ticker := time.NewTicker(UpdateInterval)
	defer ticker.Stop()

	for range ticker.C {
		err := s.updateFiles()
		if err != nil {
			fmt.Printf("Error updating site %s: %v\n", s.Name, err)
		}

		fmt.Printf("Updated site %s\n", s.Name)
	}
}

func (s *Site) updateFiles() error {
	release, err := s.Repo.GetLatestRelease()
	if err != nil {
		return err
	}

	if release.Assets == nil {
		return NoAssetError
	}
	asset := release.Assets[0]

	content, err := asset.Download()
	if err != nil {
		return err
	}

	err = unpackZip(content, fmt.Sprintf("%s/%s", "sites", s.Name))
	if err != nil {
		return err
	}

	return nil
}

func unpackZip(zipContent []byte, destinationDir string) error {
	reader, err := zip.NewReader(bytes.NewReader(zipContent), int64(len(zipContent)))
	if err != nil {
		return err
	}

	// Clear destination directory
	err = os.RemoveAll(destinationDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(destinationDir, 0755)
	if err != nil {
		return err
	}

	for _, f := range reader.File {
		path := filepath.Join(destinationDir, f.Name)

		if f.FileInfo().IsDir() {
			// Create directory
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				return err
			}
			continue
		}

		// Ensure the directory for the file exists
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}

		// Open destination file for writing
		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// Open source file in the ZIP archive
		srcFile, err := f.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// Copy file contents
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}
