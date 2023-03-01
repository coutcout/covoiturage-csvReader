package route_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"me/coutcout/covoiturage/journey/route"
	"me/coutcout/covoiturage/messaging"
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
	route.SetupRouter(&logger, r)

	type tmplTest struct {
		name     	string
		filename 	string
		statusCode  int
		message 	string
		errors 		[]string
	}

	tests := []tmplTest{
		{"nominal_case", "dataset_1.csv", http.StatusAccepted, route.MSG_FILE_ACCEPTED, []string{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, _ := os.Open(filepath.Join("testdata", test.filename))

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/import", f)
			r.ServeHTTP(w, req)

			assert.Equal(t, test.statusCode, w.Code)
			response := messaging.ResponseMessage{} 
			json.NewDecoder(w.Body).Decode(&response)
			assert.Equal(t, response.Message, test.message)
			assert.Equal(t, response.StatusCode, test.statusCode)
			assert.ElementsMatch(t, response.Errors, test.errors)
		})
	}
}