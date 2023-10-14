package videos_controller

import (
	"bytes"
	"gotube/app/helpers"
	"gotube/app/models"
	"gotube/config"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"gorm.io/gorm"
)

type MockDatabase struct{}

var raiseGormError = false

var mockVideo = models.Video{
	FilePath: "test.mp4",
	Title:    "Test Video",
}
var mockVideos = []models.Video{mockVideo}

func (db *MockDatabase) Find(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	if raiseGormError {
		return &gorm.DB{
			Error: gorm.ErrInvalidData,
		}
	}

	videos := dest.(*[]models.Video)

	*videos = mockVideos

	return &gorm.DB{
		Error: nil,
	}
}

func (db *MockDatabase) First(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	if raiseGormError {
		return &gorm.DB{
			Error: gorm.ErrRecordNotFound,
		}
	}

	video := dest.(*models.Video)
	stringId := conds[0].(string)

	id, _ := strconv.Atoi(stringId)
	mockVideo.ID = uint(id)

	*video = mockVideo

	return &gorm.DB{
		Error: nil,
	}
}

func (db *MockDatabase) Create(value interface{}) (tx *gorm.DB) {
	if raiseGormError {
		return &gorm.DB{
			Error: gorm.ErrModelValueRequired,
		}
	}

	return &gorm.DB{
		Error: nil,
	}
}

type testCase = struct {
	path                 string
	method               string
	raiseGormError       bool
	enableCSRFProtection bool
	expectedText         string
	expectedHttpStatus   int
}

var testCasesForGet = []testCase{
	{"/", "GET", false, true, mockVideo.Title, http.StatusOK},
	{"/", "GET", true, true, gorm.ErrInvalidData.Error(), http.StatusInternalServerError},
	{"/show/1", "GET", false, true, mockVideo.Title, http.StatusOK},
	{"/show/1", "GET", true, true, gorm.ErrRecordNotFound.Error(), http.StatusInternalServerError},
	{"/upload", "POST", false, false, "", http.StatusOK},
	{"/upload", "POST", false, true, "Bad Request", http.StatusBadRequest},
	{"/upload", "POST", true, false, gorm.ErrModelValueRequired.Error(), http.StatusInternalServerError},
	{"/videos/test.mp4", "GET", true, false, "", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	helpers.UploadPath = "../../../storage/uploads"
	helpers.TemplatesPath = "../../views"
	config.App.DB = &MockDatabase{}

	createTestVideoFile()

	for _, tc := range testCasesForGet {
		t.Logf("Running test for path: %s, method: %s, raiseGormError: %v", tc.path, tc.method, tc.raiseGormError)

		raiseGormError = tc.raiseGormError

		config.App.EnableCSRFProtection = tc.enableCSRFProtection
		ts := httptest.NewTLSServer(createTestRouter())
		defer ts.Close()

		var resp *http.Response
		var err error

		if tc.method == "GET" {
			resp, err = ts.Client().Get(ts.URL + tc.path)
		} else {
			var requestBody bytes.Buffer
			writer := multipart.NewWriter(&requestBody)
			writer.WriteField("title", "Test Title")
			writer.WriteField("description", "Test Description")
			part, _ := writer.CreateFormFile("videoFile", "test.mp4")

			part.Write([]byte("dummy file content"))
			writer.Close()

			resp, err = ts.Client().Post(ts.URL+tc.path, writer.FormDataContentType(), &requestBody)
		}

		if err != nil {
			t.Fatalf("HTTP request failed: %v", err)
		}
		defer resp.Body.Close()

		// Check the response status code
		if status := resp.StatusCode; status != tc.expectedHttpStatus {
			t.Errorf("Expected status %v, but got %v", tc.expectedHttpStatus, status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if !strings.Contains(string(body), tc.expectedText) {
			t.Errorf("Response body does not contain text: %s", tc.expectedText)
		}
	}
	deleteTestVideoFile()
}
