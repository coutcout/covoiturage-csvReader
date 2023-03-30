// Package usecase implements all the application usecases
package usecase

import (
	"io"
	"sync"

	"github.com/coutcout/covoiturage-csvreader/configuration"
	"github.com/coutcout/covoiturage-csvreader/domain"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

type journeyUsecase struct {
	logger           *zap.SugaredLogger
	cfg              *configuration.Config
	journeyRepo      domain.JourneyRepositoryInterface
	journeyCsvParser domain.JourneyParser
}

// NewJourneyUsecase creates a new journey usecase.
//
// @param logger - Logger to log to. Must not be nil.
// @param cfg - Configuration for the journey repository. Must not be nil.
// @param jRepo - Journey repository to use. Must not be nil.
// @param jCsvParser - Journey parser to use. Must not be nil
func NewJourneyUsecase(logger *zap.SugaredLogger, cfg *configuration.Config, jRepo domain.JourneyRepositoryInterface, jCsvParser domain.JourneyParser) domain.JourneyUsecase {
	return &journeyUsecase{
		logger:           logger,
		cfg:              cfg,
		journeyRepo:      jRepo,
		journeyCsvParser: jCsvParser,
	}
}

// ImportFromCSVFile imports journeys from a CSV file.
//
// @param reader - the reader to read the csv file
func (ucase *journeyUsecase) ImportFromCSVFile(c *gin.Context, reader io.Reader) (int64, []string) {
	journeyChan := make(chan *domain.Journey)
	errorChan := make(chan string)
	ucase.journeyCsvParser.Parse(reader, journeyChan, errorChan)
	errors := []string{}

	nbJourneyImported := 0
	var workerGroup sync.WaitGroup

	worker := func(repo domain.JourneyRepositoryInterface, journeyChan <-chan *domain.Journey, bufferSize int){
		var journeyBuffer []domain.Journey
		for {
			select {
			case journey, ok := <-journeyChan:
				if !ok {
					ucase.logger.Debug("Worker ended, flushing buffer")
					repo.Add(c, journeyBuffer)
					return
				}

				journeyBuffer = append(journeyBuffer, *journey)
				if len(journeyBuffer) == bufferSize {
					ucase.logger.Debugw("Buffer is full, flushing it",
						"bufferSize", bufferSize,
					)
					repo.Add(c, journeyBuffer)
					journeyBuffer = nil
				}
			}
		}
	}

	ucase.logger.Debug("Workers Initializing")
	for w := 0; w < ucase.cfg.Journey.Insertion.WorkerPoolSize; w++ {
		workerGroup.Add(1)
		go func() {
			defer workerGroup.Done()
			worker(ucase.journeyRepo, journeyChan, ucase.cfg.Journey.Insertion.BulkInsertSize)
		}()
	}


	workerGroup.Add(1)
	go func() {
		defer workerGroup.Done()
		for e := range errorChan {
			errors = append(errors, e)
		}
	}()

	workerGroup.Wait()

	return int64(nbJourneyImported), errors
}
