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

	"me/coutcout/covoiturage/journey/router"
	"me/coutcout/covoiturage/messaging"
	"me/coutcout/covoiturage/mocks"
)

var logger zap.SugaredLogger

func init() {
	newLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Println("Error initializing logger!")
	}

	logger = *newLogger.Sugar()
}

func TestImportCSVFile(t *testing.T){
	r := gin.Default()
	mockJUsecase := new(mocks.JourneyUsecase)
	mockJUsecase.On("ImportFromCSVFile", mock.Anything).Return(int64(1), nil)
	router.NewJourneyRouter(&logger, r, mockJUsecase)

	type tmplTest struct {
		name     	string
		filename 	string
		statusCode  int
		message 	string
		errors 		[]string
	}

	tests := []tmplTest{
		{"nominal_case", "dataset_1.csv", http.StatusAccepted, router.MSG_FILE_ACCEPTED, []string{}},
		// {"good_format_wrong_extension", "dataset_1.csv.json", http.StatusAccepted, router.MSG_FILE_ACCEPTED, []string{}},
		// {"wrong_format", "dataset_1.json", http.StatusBadRequest, "", []string{"This is not a correct CSV formatted file"}},
		{"file_too_long", "dataset_too_long.csv", http.StatusBadRequest, "", []string{"file dataset_too_long.csv is too big"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			f, err := writer.CreateFormFile("files", test.filename)
			if(err != nil){
				logger.Error(err)
			}

			file, err := os.Open(filepath.Join("testdata", test.filename))
			if(err != nil){
				logger.Error(err)
			}

			_, err2 := io.Copy(f, file)
			if(err != nil){
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

			for _, fileMessage := range response.Files {
				if fileMessage.Filename == test.filename {
					assert.Equal(t, test.message, fileMessage.Message)
					assert.ElementsMatch(t, test.errors, fileMessage.Errors)
				}
			}
			
		})
	}
}