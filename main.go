package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"syncpage/github"
	"syncpage/middleware"
	"syncpage/site"
)

const (
	DataDirPath         = "data/"
	SitesCollectionPath = "data/sites.json"
	Port                = "8080"
)

func main() {
	mux := http.NewServeMux()

	sites, err := loadSites()
	if err != nil {
		panic(err)
	}

	for _, s := range sites {
		s.Register(mux)
	}

	go func() {
		err := http.ListenAndServe(":"+Port, middleware.Middleware(mux))
		if err != nil {
			return
		}
	}()

	fmt.Printf("Listening on port %s\n", Port)

	c := make(chan bool)
	<-c
}

func loadSites() ([]site.Site, error) {
	err := os.MkdirAll(DataDirPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	file, err := os.ReadFile(SitesCollectionPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = saveDefaultSites()
			if err != nil {
				return nil, fmt.Errorf("failed to save default sites: %w", err)
			}
			fmt.Println("Default sites saved, please edit data/sites.json and restart the program")
			os.Exit(0)
			return nil, nil
		}

		return nil, fmt.Errorf("failed to read sites collection: %w", err)
	}

	var sites []site.Site
	err = json.Unmarshal(file, &sites)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal sites collection: %w", err)
	}

	return sites, nil
}

func saveDefaultSites() error {
	s := []site.Site{
		{
			Name: "docs",
			Repo: github.Repository{
				Owner:     "OWNER",
				Repo:      "REPO",
				AuthToken: "TOKEN",
			},
			WorkflowName: "WORKFLOW_NAME",
			ArtifactName: "ARTIFACT_NAME",
			FileName:     "\\.zip$",
		},
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default site: %w", err)
	}

	err = os.WriteFile(SitesCollectionPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write default site: %w", err)
	}

	return nil
}
