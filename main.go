package main

import (
	"fmt"
	"net/http"
	"regexp"
	"syncpage/github"
	"syncpage/middleware"
	"syncpage/site"
)

const (
	PORT = "8080"
)

func main() {
	mux := http.NewServeMux()

	fpDocs := site.Site{
		Name: "docs",
		Repo: github.Repository{
			Owner:     "fancymcplugins",
			Repo:      "docs",
			AuthToken: "",
		},
		WorkflowName: "Build documentation",
		ArtifactName: "docs",
		FileName:     regexp.MustCompile("\\.zip$"),
	}
	fpDocs.Register(mux)

	go func() {
		err := http.ListenAndServe(":"+PORT, middleware.Middleware(mux))
		if err != nil {
			return
		}
	}()

	fmt.Printf("Listening on port %s\n", PORT)

	c := make(chan bool)
	<-c
}
