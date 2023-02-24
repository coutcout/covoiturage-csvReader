package repositories

import (
	"me/coutcout/covoiturage/domain"
)

type JourneyRepositoryInterface interface {
	Add(journey domain.Journey) (bool, error)
}
