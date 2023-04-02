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

	"github.com/coutcout/covoiturage-csvreader/configuration"
	"github.com/coutcout/covoiturage-csvreader/domain"
	"github.com/coutcout/covoiturage-csvreader/journey/router"
	"github.com/coutcout/covoiturage-csvreader/journey/usecase"
	"github.com/coutcout/covoiturage-csvreader/messaging"
	"github.com/coutcout/covoiturage-csvreader/mocks"
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
			mock := mockJUsecase.On("ImportFromCSVFile", mock.Anything, mock.Anything)
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
			req, _ := http.NewRequest("POST", "/journey/import", body)
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
		req, _ := http.NewRequest("POST", "/journey/import", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := messaging.SingleResponseMessage{}
		json.NewDecoder(w.Body).Decode(&response)

		assert.NotEmpty(t, response.Errors)
	})
}


func TestGetJourneys(t *testing.T) {
	r := gin.Default()
	jRepoMock := mocks.NewJourneyRepositoryInterface(t)
	jParserMock := mocks.NewJourneyParser(t)
	jUCase := usecase.NewJourneyUsecase(&logger, config, jRepoMock, jParserMock)
	router.NewJourneyRouter(&logger, config, r, jUCase)

	type test struct {
		name string
	}

	tests := []test{
		{
			name: "Nominal",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jRepoMock.On("FindAll", mock.Anything, int64(0), config.Journey.Get.Stream.BufferSize).Return([]domain.Journey{
				domain.Journey{JourneyId: 1},
				domain.Journey{JourneyId: 2},
			}, nil).Once()

			jRepoMock.On("FindAll", mock.Anything, config.Journey.Get.Stream.BufferSize, config.Journey.Get.Stream.BufferSize).Return([]domain.Journey{
				domain.Journey{JourneyId: 3},
				domain.Journey{JourneyId: 4},
			}, nil).Once()

			jRepoMock.On("FindAll", mock.Anything, 2*config.Journey.Get.Stream.BufferSize, config.Journey.Get.Stream.BufferSize).Return([]domain.Journey{}, nil).Once()

			w := CreateTestResponseRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/journey", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var receivedJourneys []domain.Journey
			var journeysBuffer []domain.Journey
			decoder := json.NewDecoder(w.Body)
			nbStreams := 0
			for {
				err := decoder.Decode(&journeysBuffer)
				logger.Debugw("Receiving bytes of response",
					"data", journeysBuffer,
				)

				if err == io.EOF {
					break
				} else if err != nil {
					logger.Errorw("Error receiving stream",
						"error", err,
					)
				}
				receivedJourneys = append(receivedJourneys, journeysBuffer...)
				nbStreams++
			}
			
			assert.Equal(t, 4, len(receivedJourneys))
			assert.Equal(t, 2, nbStreams)
		})
	}
}

type TestResponseRecorder struct {
	*httptest.ResponseRecorder
	closeChannel chan bool
}

func (r *TestResponseRecorder) CloseNotify() <-chan bool {
	return r.closeChannel
}

func (r *TestResponseRecorder) closeClient() {
	r.closeChannel <- true
}

func CreateTestResponseRecorder() *TestResponseRecorder {
	return &TestResponseRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}