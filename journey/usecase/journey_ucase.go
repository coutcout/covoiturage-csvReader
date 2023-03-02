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

func (ucase *journeyUsecase) ImportFromCSVFile(reader io.Reader){

}