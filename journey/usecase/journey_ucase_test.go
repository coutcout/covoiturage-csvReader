// Package usecase_test tests all the application usecases
package usecase_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/coutcout/covoiturage-csvReader/configuration"
	"github.com/coutcout/covoiturage-csvReader/journey/service"
	"github.com/coutcout/covoiturage-csvReader/journey/usecase"
	"github.com/coutcout/covoiturage-csvReader/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
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

func TestImportFromCSVFile(t *testing.T) {
	type tmplTest struct {
		name             string
		filename         string
		nbAdded          int64
		shouldHaveErrors bool
	}

	tests := []tmplTest{
		{"nominal_case", "dataset_1.csv", 3, false},
		{"27fields_case", "dataset_27fields.csv", 5, false},
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
			jCsvParser := service.NewJourneyCsvParser(&logger, config)

			journeyUsecase := usecase.NewJourneyUsecase(
				&logger,
				config,
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
