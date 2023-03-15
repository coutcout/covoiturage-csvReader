package repository

import (
	"me/coutcout/covoiturage/domain"

	"go.uber.org/zap"
)

type dbJourneyRepository struct {
	logger *zap.SugaredLogger
}

func NewDbJourneyRepository(logger *zap.SugaredLogger) domain.JourneyRepositoryInterface {
	return &dbJourneyRepository{
		logger: logger,
	}
}

func (r *dbJourneyRepository) Add(journey *domain.Journey) (bool, error) {
	r.logger.Debugw("Not Implemented yet",
		"object", journey,
	)
	return true, nil
}