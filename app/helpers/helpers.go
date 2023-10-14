package helpers

import (
	"encoding/hex"
	"fmt"
	"gotube/config"
	"math/rand"
	"net/http"
	"text/template"
	"time"

	"github.com/justinas/nosurf"
)

// Represents data needed for rendering HTML templates
type templateData struct {
	Data         interface{} // Dynamic data to be rendered in templates
	CSRFToken    string      // Cross-Site Request Forgery (CSRF) token
	FlashSuccess string      // Success flash message
	FlashError   string      // Error flash message
	FlashWarning string      // Warning flash message
}

var TemplatesPath = "./app/views"
var UploadPath = "./storage/uploads"

// Renders an HTML template with the given data
func RenderTemplate(r *http.Request, w http.ResponseWriter, modelName string, templateName string, data interface{}) {
	// Parse HTML templates for the specified model and layout
	t, err :=
		template.ParseFiles(
			fmt.Sprintf("%s/%s/%s.page.tmpl", TemplatesPath, modelName, templateName),
			fmt.Sprintf("%s/layouts/base.layout.tmpl", TemplatesPath),
		)

	if err != nil {
		LogAndDisplayError(err, w)
		return
	}

	// Prepare template data
	temlplateData := templateData{
		Data:         data,
		CSRFToken:    nosurf.Token(r),
		FlashSuccess: config.App.Session.PopString(r.Context(), "flashSuccess"),
		FlashError:   config.App.Session.PopString(r.Context(), "flashError"),
		FlashWarning: config.App.Session.PopString(r.Context(), "flashWarning"),
	}

	// Execute the template with the provided data
	t.ExecuteTemplate(w, "base", temlplateData)
}

// Logs an error and sends a 500 Internal Server Error response
func LogAndDisplayError(err error, w http.ResponseWriter) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	config.App.ErrorLog.Println(err)
}

// Generates a random hexadecimal ID of the specified length
func GenerateRandomHexID(length int) (string, error) {
	bytes := make([]byte, length)

	// Seed the random number generator with the current Unix timestamp
	rand.Seed(time.Now().UnixNano())

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Convert the random bytes to a hexadecimal string of the specified length
	hexID := hex.EncodeToString(bytes)[0:length]

	return hexID, nil
}
