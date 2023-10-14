package videos_controller

import (
	"fmt"
	"gotube/app/helpers"
	"gotube/config"
	"os"

	"github.com/go-chi/chi/v5"
)

var testVideoFileName = "test.mp4"

func createTestRouter() *chi.Mux {
	router := config.BaseRouter()
	router.Get("/", IndexHandler)
	router.Post("/upload", UploadHandler)
	router.Get("/videos/*", ServeVideoHandler)
	router.Get("/show/{id}", ShowHandler)
	return router
}

func createTestVideoFile() {
	_, err := os.Create(helpers.UploadPath + "/" + testVideoFileName)

	if err != nil {
		panic(err)
	}

	fmt.Println("Created test video file")
}

func deleteTestVideoFile() {
	err := os.Remove(helpers.UploadPath + "/" + testVideoFileName)

	if err != nil {
		panic(err)
	}

	fmt.Println("Deleted test video file")
}
