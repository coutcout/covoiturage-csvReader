// Test usecases
package usecase_test

import (
	"log"
	"me/coutcout/covoiturage/journey/service"
	"me/coutcout/covoiturage/journey/usecase"
	"me/coutcout/covoiturage/mocks"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

var logger zap.SugaredLogger

func init() {
	newLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Println("Error initializing logger!")
	}

	logger = *newLogger.Sugar()
}

func TestImportFromCSVFile(t *testing.T) {
	type tmplTest struct {
		name     string
		filename string
		nbAdded  int64
		shouldHaveErrors 	 bool
	}

	tests := []tmplTest{
		{"nominal_case", "dataset_1.csv", 3, false},
		{"empty_file_case", "dataset_empty.csv", 0, false},
		{"headers_only_case", "dataset_headersOnly.csv", 0, false},
		{"json", "dataset_1.json", 0, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, _ := os.Open(filepath.Join("testdata", test.filename))
			defer f.Close()

			jRepo := new(mocks.JourneyRepositoryInterface)
			jRepo.On("Add", mock.AnythingOfType("*domain.Journey")).Return(true, nil)
			jCsvParser := service.NewJourneyCsvParser(&logger)

			journeyUsecase := usecase.NewJourneyUsecase(
				&logger,
				jRepo,
			
				jCsvParser,
			)
			nbJourneyImported, err := journeyUsecase.ImportFromCSVFile(f)

			if test.shouldHaveErrors {
				assert.NotEmpty(t, err)
			} else {
				assert.Empty(t, err)
			}
			assert.Equal(t, test.nbAdded, nbJourneyImported)

			logger.Debugw("End of the test",
				"file", test.filename,
			)
		})
	}

}