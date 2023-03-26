// Manage data
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

// Constructor
func NewDbJourneyRepository(logger *zap.SugaredLogger, cfg *configuration.Config) domain.JourneyRepositoryInterface {
	return &dbJourneyRepository{
		logger: logger,
		cfg: cfg,
	}
}

func (r *dbJourneyRepository) Add(journey *domain.Journey) (bool, error) {
	r.logger.Debugw("Not Implemented yet",
		"object", journey,
	)
	return true, nil
}
