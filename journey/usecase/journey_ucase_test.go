package usecase_test

import (
	"log"
	"me/coutcout/covoiturage/journey/usecase"
	"me/coutcout/covoiturage/mocks"
	"os"
	"path/filepath"
	"testing"

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

			jRepo := new(mocks.JourneyRepositoryInterface)
			jCsvParser := new(mocks.JourneyParser)

			journeyUsecase := usecase.NewJourneyUsecase(
				&logger,
				jRepo,
				jCsvParser,
			)
			journeyUsecase.ImportFromCSVFile(f)
		})
	}

}