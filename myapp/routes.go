package main

import (
	// "log"
	"net/http"

	"github.com/djedjethai/celeritas"
	// "github.com/djedjethai/celeritas/filesystems/miniofilesystem"
	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes

	// add routes here
	a.get("/", a.Handlers.Home)

	a.App.Routes.Get("/users/login", a.Handlers.UserLogin)
	a.App.Routes.Post("/users/login", a.Handlers.PostUserLogin)

	a.App.Routes.Get("/auth/{provider}", a.Handlers.SocialLogin)
	a.App.Routes.Get("/auth/{provider}/callback", a.Handlers.SocialMediaCallback)

	a.get("/upload", a.Handlers.CeleritasUpload)
	a.post("/upload", a.Handlers.PostCeleritasUpload)

	a.get("/list-fs", a.Handlers.ListFS)
	a.get("/files/upload", a.Handlers.UploadToFS)
	a.post("/files/upload", a.Handlers.PostUploadToFS)

	a.get("/delete-from-fs", a.Handlers.DeleteFromFS)

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	// routes from celeritas
	a.App.Routes.Mount("/celeritas", celeritas.Routes())
	a.App.Routes.Mount("/api", a.ApiRoutes())

	return a.App.Routes
}
