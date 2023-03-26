// Package repo manage data
package repo

import (
	"me/coutcout/covoiturage/domain"
	"me/coutcout/covoiturage/configuration"

	"go.uber.org/zap"
)

type dbJourneyRepository struct {
	logger *zap.SugaredLogger
	cfg *configuration.Config
}

// NewDbJourneyRepository make an instance of a dbJourneyRepository
func NewDbJourneyRepository(logger *zap.SugaredLogger, cfg *configuration.Config) domain.JourneyRepositoryInterface {
	return &dbJourneyRepository{
		logger: logger,
		cfg: cfg,
	}
}

// Add adds journey to the repository. Not Implemented yet.
// 
// @param r - jouney the journey to be
// @param journey
func (r *dbJourneyRepository) Add(journey *domain.Journey) (bool, error) {
	r.logger.Debugw("Not Implemented yet",
		"object", journey,
	)
	return true, nil
}
