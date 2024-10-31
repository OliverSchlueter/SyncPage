package main

import (
	"net/http"
	"syncpage/github"
	"syncpage/site"
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
		err := http.ListenAndServe(":8080", mux)
		if err != nil {
			return
		}
	}()

	c := make(chan bool)
	<-c
}
