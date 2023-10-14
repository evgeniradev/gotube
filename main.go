// main.go
package main

import (
	"fmt"
	"gotube/app/controllers/videos_controller"
	"gotube/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Starting application...")

	config.App.DB = config.ConnectDB()

	http.ListenAndServe(":8080", router())
}

func router() *chi.Mux {
	router := config.BaseRouter()
	router.Get("/", videos_controller.IndexHandler)
	router.Post("/upload", videos_controller.UploadHandler)
	router.Get("/videos/*", videos_controller.ServeVideoHandler)
	router.Get("/show/{id}", videos_controller.ShowHandler)
	return router
}
