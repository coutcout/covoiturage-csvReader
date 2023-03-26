// Package router_test tests all the routing configuration
package router_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/coutcout/covoiturage-csvReader/configuration"
	"github.com/coutcout/covoiturage-csvReader/journey/router"
	"github.com/coutcout/covoiturage-csvReader/messaging"
	"github.com/coutcout/covoiturage-csvReader/mocks"
)

var logger zap.SugaredLogger
var config *configuration.Config

func init() {
	newLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Println("Error initializing logger!")
	}

	logger = *newLogger.Sugar()

	config, err = configuration.NewConfig("../../resource/configurations/application-tu.yaml")
	if err != nil {
		logger.Error(err)
	}
}

func TestImportCSVFile(t *testing.T) {
	r := gin.Default()
	mockJUsecase := new(mocks.JourneyUsecase)
	router.NewJourneyRouter(&logger, config, r, mockJUsecase)

	type tmplTest struct {
		name                 string
		filename             string
		statusCode           int
		hasErrors            bool
		expectedImportedLine int
	}

	tests := []tmplTest{
		{"nominal_case", "dataset_1.csv", http.StatusAccepted, false, 3},
		{"good_format_wrong_extension", "dataset_1.csv.json", http.StatusAccepted, false, 3},
		{"wrong_format", "dataset_1.json", http.StatusBadRequest, true, 0},
		{"wrong_format_good_extension", "dataset_1.json.csv", http.StatusBadRequest, true, 0},
		{"file_too_long", "dataset_too_long.csv", http.StatusBadRequest, false, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var returnedErrors []string
			if test.hasErrors {
				returnedErrors = []string{"error"}
			}
			mock := mockJUsecase.On("ImportFromCSVFile", mock.Anything)
			mock.Return(int64(test.expectedImportedLine), returnedErrors)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			f, err := writer.CreateFormFile("files", test.filename)
			if err != nil {
				logger.Error(err)
			}

			file, err := os.Open(filepath.Join("testdata", test.filename))
			if err != nil {
				logger.Error(err)
			}

			_, err2 := io.Copy(f, file)
			if err != nil {
				logger.Error(err2)
			}

			writer.Close()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/import", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			r.ServeHTTP(w, req)

			assert.Equal(t, test.statusCode, w.Code)
			response := messaging.MultipleResponseMessage{}
			json.NewDecoder(w.Body).Decode(&response)

			assert.NotEmpty(t, response.Files)
			for _, fileMessage := range response.Files {
				if fileMessage.Filename == test.filename {
					assert.Equal(t, test.expectedImportedLine, fileMessage.NbLineImported)
					if test.hasErrors {
						assert.NotEmpty(t, fileMessage.Errors)
					}
				}
			}

			mock.Unset()
		})
	}
}

func TestImportCSVFile_wrongParameter(t *testing.T) {
	r := gin.Default()
	mockJUsecase := new(mocks.JourneyUsecase)
	router.NewJourneyRouter(&logger, config, r, mockJUsecase)

	t.Run("Wrong parameter name", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		f, err := writer.CreateFormFile("wrongFile", "dataset_1.csv")
		if err != nil {
			logger.Error(err)
		}

		file, err := os.Open(filepath.Join("testdata", "dataset_1.csv"))
		if err != nil {
			logger.Error(err)
		}

		_, err2 := io.Copy(f, file)
		if err != nil {
			logger.Error(err2)
		}

		writer.Close()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/import", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := messaging.SingleResponseMessage{}
		json.NewDecoder(w.Body).Decode(&response)

		assert.NotEmpty(t, response.Errors)
	})
}
