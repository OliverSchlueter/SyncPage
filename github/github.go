package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Repository struct {
	Owner     string
	Repo      string
	AuthToken string
}

type Release struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func (r *Repository) GetLatestRelease() (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", r.Owner, r.Repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if r.AuthToken != "" {
		req.Header.Set("Authorization", "token "+r.AuthToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func (a *Asset) Download() ([]byte, error) {
	resp, err := http.Get(a.BrowserDownloadURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
