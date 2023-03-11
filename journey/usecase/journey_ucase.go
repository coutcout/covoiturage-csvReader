package usecase

import (
	"io"
	"me/coutcout/covoiturage/domain"

	"go.uber.org/zap"
)

type journeyUsecase struct {
	logger				*zap.SugaredLogger
	journeyRepo 		domain.JourneyRepositoryInterface
	journeyCsvParser 	domain.JourneyParser
}

func NewJourneyUsecase(logger *zap.SugaredLogger, jRepo domain.JourneyRepositoryInterface, jCsvParser domain.JourneyParser) domain.JourneyUsecase {
	return &journeyUsecase{
		logger: logger,
		journeyRepo: jRepo,
		journeyCsvParser: jCsvParser,
	}
}

func (ucase *journeyUsecase) ImportFromCSVFile(reader io.Reader) (int64, error) {
	journeyChan := make(chan *domain.Journey)
	err := ucase.journeyCsvParser.Parse(reader, journeyChan)

	nbJourneyImported := 0
	for j := range journeyChan {
		if res, err := ucase.journeyRepo.Add(j); err == nil && res {
			nbJourneyImported += 1
		}
	}

	return int64(nbJourneyImported), err
}