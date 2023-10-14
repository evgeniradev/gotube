package helpers

import (
	"bytes"
	"errors"
	"gotube/app/models"
	"gotube/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	TemplatesPath = "../views"
	request := http.Request{}
	writer := httptest.NewRecorder()
	data := make(map[string]interface{})
	data["Videos"] = []models.Video{}
	data["Title"] = "Page Title"

	config.SessionLoadMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RenderTemplate(r, w, "videos", "index", data)
	})).ServeHTTP(writer, &request)

	if writer.Code != http.StatusOK {
		t.Errorf("Did not return a %d status code", http.StatusOK)
	}

	expectedBody := data["Title"].(string)
	if body := writer.Body.String(); !strings.Contains(string(body), expectedBody) {
		t.Errorf("Did not incliude '%s' in the response body", expectedBody)
	}
}

func TestLogAndDisplayError(t *testing.T) {
	errorMessage := "error_test_message"
	raisedError := errors.New(errorMessage)

	writer := httptest.NewRecorder()

	logBuffer := bytes.Buffer{}
	config.App.ErrorLog.SetOutput(&logBuffer)

	LogAndDisplayError(raisedError, writer)

	if writer.Code != http.StatusInternalServerError {
		t.Errorf("Did not return a %d status code", http.StatusInternalServerError)
	}

	if body := writer.Body.String(); !strings.Contains(string(body), errorMessage) {
		t.Errorf("Did not incliude '%s' in the response body", errorMessage)
	}

	if logOutput := logBuffer.String(); !strings.Contains(logOutput, errorMessage) {
		t.Errorf("Did not log '%s' error message", logOutput)
	}
}

func TestGenerateRandomHexID(t *testing.T) {
	idLength := 10
	id1, err := GenerateRandomHexID(idLength)

	if err != nil {
		t.Errorf("Did not generate a random hexadecimal ID: %v", err)
	}

	if len(id1) != idLength {
		t.Errorf("Did not generate a random hexadecimal ID of length %d", idLength)
	}

	id2, err := GenerateRandomHexID(idLength)

	if err != nil {
		t.Errorf("Did not generate a random hexadecimal ID: %v", err)
	}

	if id1 == id2 {
		t.Error("Did not generate a random hexadecimal")
	}
}
