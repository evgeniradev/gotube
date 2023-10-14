package videos_controller

import (
	"gotube/app/helpers"
	"gotube/app/models"
	"gotube/config"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi/v5"
)

// Handles the HTTP GET request for the video index page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Get all videos from the database
	var videos []models.Video
	result := config.App.DB.Find(&videos)
	if result.Error != nil {
		helpers.LogAndDisplayError(result.Error, w)
		return
	}

	// Prepare data for rendering the template
	data := make(map[string]interface{})
	data["Videos"] = videos
	data["Title"] = "Videos"

	// Render the HTML template
	helpers.RenderTemplate(r, w, "videos", "index", data)
}

// Handles the HTTP GET request for viewing a specific video
func ShowHandler(w http.ResponseWriter, r *http.Request) {
	var video models.Video

	id := chi.URLParam(r, "id")

	// Retrieve the video details from the database based on the ID
	result := config.App.DB.First(&video, id)
	if result.Error != nil {
		helpers.LogAndDisplayError(result.Error, w)
		return
	}

	// Prepare data for rendering the template
	data := make(map[string]interface{})
	data["Video"] = video
	data["Title"] = video.Title

	// Render the HTML template
	helpers.RenderTemplate(r, w, "videos", "show", data)
}

// Handles the HTTP POST request for video uploads
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form with a 10MB limit
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		helpers.LogAndDisplayError(err, w)
		return
	}

	// Retrieve the uploaded file
	file, fileHeader, err := r.FormFile("videoFile")
	if err != nil {
		helpers.LogAndDisplayError(err, w)
		return
	}
	defer file.Close()

	// Extract the file extension
	fileExtension := path.Ext(fileHeader.Filename)

	// Generate a random file name with a 10-character hexadecimal value
	randomHexValue, err := helpers.GenerateRandomHexID(10)
	if err != nil {
		helpers.LogAndDisplayError(err, w)
		return
	}

	// Create a new file with the generated name in the "storage/uploads" directory
	newFileName := randomHexValue + fileExtension
	uploadedFile, err := os.Create(helpers.UploadPath + "/" + newFileName)
	if err != nil {
		helpers.LogAndDisplayError(err, w)
		return
	}
	defer uploadedFile.Close()

	// Copy the uploaded file data to the newly created file
	_, err = io.Copy(uploadedFile, file)
	if err != nil {
		helpers.LogAndDisplayError(err, w)
		return
	}

	// Create a new video record in the database
	video := &models.Video{
		FilePath:    newFileName,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}
	result := config.App.DB.Create(video)
	if result.Error != nil {
		helpers.LogAndDisplayError(result.Error, w)
		return
	}

	// Set a success flash message and redirect to the home page
	config.App.Session.Put(r.Context(), "flashSuccess", "Video uploaded successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handles the HTTP GET request for serving video files
func ServeVideoHandler(w http.ResponseWriter, r *http.Request) {
	// Create a file server handler to serve video files from the "storage/uploads" directory
	fs := http.StripPrefix("/videos", http.FileServer(http.Dir(helpers.UploadPath)))
	fs.ServeHTTP(w, r)
}
