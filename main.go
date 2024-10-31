package main

import (
	"fmt"
	"net/http"
	"syncpage/github"
	"syncpage/middleware"
	"syncpage/site"
)

const (
	PORT = "8080"
)

func main() {
	mux := http.NewServeMux()

	fnSite := site.Site{
		Name: "FancyNpcs",
		Repo: github.Repository{
			Owner: "fancymcplugins",
			Repo:  "fancynpcs",
		},
	}
	fnSite.Register(mux)

	fhSite := site.Site{
		Name: "FancyHolograms",
		Repo: github.Repository{
			Owner: "fancymcplugins",
			Repo:  "fancyholograms",
		},
	}
	fhSite.Register(mux)

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
