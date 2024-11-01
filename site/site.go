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
	"regexp"
	"strings"
	"syncpage/github"
	"time"
)

const (
	SitesDir       = "data/sites"
	UpdateInterval = 5 * time.Minute
)

var (
	AdminAuthToken = os.Getenv("ADMIN_AUTH_TOKEN")

	ErrNoArtifacts = errors.New("no artifacts found in workflow run")
)

type Site struct {
	Name         string            `json:"name"`
	Repo         github.Repository `json:"repo"`
	WorkflowName string            `json:"workflow_name"`
	ArtifactName string            `json:"artifact_name"`
	FileName     string            `json:"file_name"`
}

func (s *Site) Register(mux *http.ServeMux) {
	go s.startUpdateLoop()
	if err := s.updateFiles(); err != nil {
		fmt.Printf("error while updating site %s: %v\n", s.Name, err)
	}

	mux.HandleFunc("/api/v1/update/"+s.Name, s.HandleForceUpdate)

	pattern := "/" + s.Name + "/"
	if s.Name == "home" {
		pattern = "/"
	}
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/") {
			path += "index.html"
		}

		content, err := os.ReadFile(SitesDir + "/" + s.Name + "/" + path)
		if err != nil {
			if os.IsNotExist(err) {
				s.TryToReturnIndex(w, r)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if strings.HasSuffix(path, ".html") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		} else if strings.HasSuffix(path, ".css") {
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		} else if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		}

		w.Write(content)
	})

	fmt.Printf("Registered site %s\n", s.Name)
}

func (s *Site) TryToReturnIndex(w http.ResponseWriter, r *http.Request) {
	path := SitesDir + "/" + s.Name + "/index.html"

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

func (s *Site) HandleForceUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != mustGetAdminAuthToken() {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err := s.updateFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Site) startUpdateLoop() {
	ticker := time.NewTicker(UpdateInterval)
	defer ticker.Stop()

	for range ticker.C {
		err := s.updateFiles()
		if err != nil {
			fmt.Printf("error while updating site %s: %v\n", s.Name, err)
		}
	}
}

func (s *Site) updateFiles() error {
	run, err := s.Repo.GetLatestWorkflowRun(s.WorkflowName)
	if err != nil {
		return fmt.Errorf("error while getting latest workflow run: %w", err)
	}

	artifacts, err := s.Repo.GetArtifacts(run)
	if err != nil {
		return fmt.Errorf("error while getting artifacts: %w", err)
	}

	if len(artifacts) == 0 {
		return ErrNoArtifacts
	}
	var artifact github.Artifact
	for _, a := range artifacts {
		if a.Name == s.ArtifactName {
			artifact = a
			break
		}
	}

	artifactRawContent, err := artifact.Download(&s.Repo)
	if err != nil {
		return fmt.Errorf("error while downloading artifact: %w", err)
	}

	artifactContent, err := readFileFromZip(artifactRawContent, regexp.MustCompile(s.FileName))
	if err != nil {
		return fmt.Errorf("error while unpacking outer artifact: %w", err)
	}

	err = unpackZip(artifactContent, SitesDir+"/"+s.Name)
	if err != nil {
		return fmt.Errorf("error while unpacking inner artifact: %w", err)
	}

	fmt.Printf("Updated site %s\n", s.Name)
	return nil
}

func readFileFromZip(zipContent []byte, filename *regexp.Regexp) ([]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipContent), int64(len(zipContent)))
	if err != nil {
		return nil, err
	}

	var file *zip.File
	for _, f := range reader.File {
		if filename.MatchString(f.Name) {
			file = f
		}
	}
	if file == nil {
		return nil, fmt.Errorf("file %s not found in zip", filename)
	}

	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	content, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, err
	}

	return content, nil
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

func mustGetAdminAuthToken() string {
	if AdminAuthToken != "" {
		return AdminAuthToken
	}

	AdminAuthToken = os.Getenv("ADMIN_AUTH_TOKEN")
	if AdminAuthToken == "" {
		AdminAuthToken = "admin"
	}

	return AdminAuthToken
}
