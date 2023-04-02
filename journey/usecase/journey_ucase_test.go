// Package usecase_test tests all the application usecases
package usecase_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/coutcout/covoiturage-csvreader/configuration"
	"github.com/coutcout/covoiturage-csvreader/journey/service"
	"github.com/coutcout/covoiturage-csvreader/journey/usecase"
	"github.com/coutcout/covoiturage-csvreader/mocks"
	"github.com/gin-gonic/gin"

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
		nbAdded          int
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
			jRepo.On("Add", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("[]domain.Journey")).Return(test.nbAdded, nil)
			jCsvParser := service.NewJourneyCsvParser(&logger, config)

			journeyUsecase := usecase.NewJourneyUsecase(
				&logger,
				config,
				jRepo,
				jCsvParser,
			)
			nbJourneyImported, err := journeyUsecase.ImportFromCSVFile(&gin.Context{}, f)

			if test.shouldHaveErrors {
				assert.NotEmpty(t, err)
			} else {
				assert.Empty(t, err)
			}
			assert.Equal(t, test.nbAdded, int(nbJourneyImported))

			logger.Debugw("End of the test",
				"file", test.filename,
			)
		})
	}
}