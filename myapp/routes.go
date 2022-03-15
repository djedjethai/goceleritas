package main

import (
	"log"
	"net/http"

	"github.com/djedjethai/celeritas"
	"github.com/djedjethai/celeritas/filesystems/miniofilesystem"
	"github.com/djedjethai/celeritas/testfolder"
	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes

	// add routes here
	a.get("/", a.Handlers.Home)
	a.get("/test-route", testfolder.TestHandler)
	a.get("/test-minio", func(w http.ResponseWriter, r *http.Request) {
		// cast to miniofilesystem
		f := a.App.FileSystems["MINIO"].(miniofilesystem.Minio)
		files, err := f.List("")
		if err != nil {
			log.Println(err)
			return
		}

		for _, file := range files {
			log.Println(file)
		}

	})

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	// routes from celeritas
	a.App.Routes.Mount("/celeritas", celeritas.Routes())
	a.App.Routes.Mount("/api", a.ApiRoutes())

	return a.App.Routes
}
