// Define all app usecases
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

// Constructor
func NewJourneyUsecase(logger *zap.SugaredLogger, cfg *configuration.Config, jRepo domain.JourneyRepositoryInterface, jCsvParser domain.JourneyParser) domain.JourneyUsecase {
	return &journeyUsecase{
		logger: logger,
		cfg: cfg,
		journeyRepo: jRepo,
		journeyCsvParser: jCsvParser,
	}
}

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