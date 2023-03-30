// Package repo manage data
package repo

import (
	"github.com/coutcout/covoiturage-csvreader/configuration"
	"github.com/coutcout/covoiturage-csvreader/domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"go.uber.org/zap"
)

type dbJourneyRepository struct {
	logger       *zap.SugaredLogger
	cfg          *configuration.Config
	dbConnection *mongo.Database
	journeyCollection *mongo.Collection
}

const journeyCollectionName = "journey"

// NewDbJourneyRepository make an instance of a dbJourneyRepository
func NewDbJourneyMongoRepository(logger *zap.SugaredLogger, cfg *configuration.Config, mongoDb *mongo.Database) domain.JourneyRepositoryInterface {
	dbJourneyRepository := &dbJourneyRepository{
		logger:      logger,
		cfg:         cfg,
		dbConnection: mongoDb,
	}
	dbJourneyRepository.journeyCollection = mongoDb.Collection(journeyCollectionName)

	return dbJourneyRepository
}


// Add adds journey to the repository. Not Implemented yet.
// 
// @param r - jouey the journey to be
// @param journey
func (r *dbJourneyRepository) Add(c *gin.Context, journey *domain.Journey) (bool, error) {

	result, err := r.journeyCollection.InsertOne(c, *journey)
	return result.InsertedID != nil, err
}
