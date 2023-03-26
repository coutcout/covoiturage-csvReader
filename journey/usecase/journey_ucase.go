// package usecase implements all the application usecases
package usecase

import (
	"io"
	"me/coutcout/covoiturage/configuration"
	"me/coutcout/covoiturage/domain"
	"sync"

	"go.uber.org/zap"
)

type journeyUsecase struct {
	logger				*zap.SugaredLogger
	cfg 				*configuration.Config
	journeyRepo 		domain.JourneyRepositoryInterface
	journeyCsvParser 	domain.JourneyParser

}

// NewJourneyUsecase creates a new journey usecase.
// 
// @param logger - Logger to log to. Must not be nil.
// @param cfg - Configuration for the journey repository. Must not be nil.
// @param jRepo - Journey repository to use. Must not be nil.
// @param jCsvParser - Journey parser to use. Must not be nil
func NewJourneyUsecase(logger *zap.SugaredLogger, cfg *configuration.Config, jRepo domain.JourneyRepositoryInterface, jCsvParser domain.JourneyParser) domain.JourneyUsecase {
	return &journeyUsecase{
		logger: logger,
		cfg: cfg,
		journeyRepo: jRepo,
		journeyCsvParser: jCsvParser,
	}
}

// ImportFromCSVFile imports journeys from a CSV file.
// 
// @param reader - the reader to read the csv file
func (ucase *journeyUsecase) ImportFromCSVFile(reader io.Reader) (int64, []string) {
	journeyChan := make(chan *domain.Journey)
	errorChan := make(chan string)
	ucase.journeyCsvParser.Parse(reader, journeyChan, errorChan)
	errors := []string{}

	nbJourneyImported := 0
	var workerGroup sync.WaitGroup
	
	workerGroup.Add(1)
	go func(){
		defer workerGroup.Done()
		for j := range journeyChan {
			if res, err := ucase.journeyRepo.Add(j); err == nil && res {
				nbJourneyImported ++
			}
		}
	}()

	workerGroup.Add(1)
	go func(){
		defer workerGroup.Done()
		for e := range errorChan {
			errors = append(errors, e)
		}
	}()

	workerGroup.Wait()

	return int64(nbJourneyImported), errors
}