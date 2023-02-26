package usecase_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"me/coutcout/covoiturage/domain"
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
		{"empty_file_case", "dataset_empty.csv", 0},
		{"headers_only_case", "dataset_headersOnly.csv", 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, _ := os.Open(filepath.Join("testdata", test.filename))
			defer f.Close()
			resChan := make(chan *domain.Journey)
			var nbReadedLines int64 = 0

			err := ucase.Parse(&logger, f, resChan)

			for range resChan {
				nbReadedLines += 1
			}
			assert.NoError(t, err)
			assert.Equal(t, test.nbAdded, nbReadedLines)
			logger.Debugw("End of the test",
				"file", test.filename,
			)
		})
	}

}
