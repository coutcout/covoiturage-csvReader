package usecase_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	ucase "me/coutcout/covoiturage/journey/usecase"
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

func TestParse(t *testing.T) {
	mockJourneyRepo := new(mocks.JourneyRepositoryInterface)
	mockJourneyRepo.On("Add", mock.Anything).Return(true, nil)

	type tmplTest struct {
		name     string
		filename string
		nbAdded  int64
	}

	tests := []tmplTest{
		{"nominal_case", "dataset_1.csv", 3},
		{"nominal_case", "dataset_empty.csv", 0},
		{"nominal_case", "dataset_headersOnly.csv", 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, _ := os.Open(filepath.Join("testdata", test.filename))
			defer f.Close()

			res, err := ucase.Parse(&logger, f)
			assert.NoError(t, err)
			assert.Equal(t, test.nbAdded, res)
		})
	}

}
