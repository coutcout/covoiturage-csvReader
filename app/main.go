// Package main contains the main file
package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/coutcout/covoiturage-csvreader/configuration"
	"github.com/coutcout/covoiturage-csvreader/journey/repo"
	"github.com/coutcout/covoiturage-csvreader/journey/router"
	"github.com/coutcout/covoiturage-csvreader/journey/service"
	"github.com/coutcout/covoiturage-csvreader/journey/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger zap.SugaredLogger

func main() {
	params, _, err := configuration.ParseFlag(os.Args[0], os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := configuration.NewConfig(params.ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}

	newLogger, _ := zap.NewProduction()
	logger = *newLogger.Sugar()
	r := gin.Default()

	// Repositories
	// ** DB Connections **
	configMongo := cfg.Database.Mongo
	dbOpts := options.Client().ApplyURI("mongodb://" + configMongo.Username + ":" + configMongo.Password + "@" + configMongo.Hostname + ":" + configMongo.Port + "/" + configMongo.Options)
	mongoClient, err := mongo.Connect(context.TODO(), dbOpts)
	if err != nil {
		log.Fatal(err)
	}

	mongoDB := mongoClient.Database(configMongo.DbName)
	journeyRepo := repo.NewDbJourneyMongoRepository(
		&logger,
		cfg,
		mongoDB,
	)

	// Services
	journeyParser := service.NewJourneyCsvParser(
		&logger,
		cfg,
	)

	// Usecases
	journeyUC := usecase.NewJourneyUsecase(
		&logger,
		cfg,
		journeyRepo,
		journeyParser,
	)

	router.NewJourneyRouter(
		&logger,
		cfg,
		r,
		journeyUC,
	)

	log.Fatal(r.Run(cfg.Server.Host + ":" + cfg.Server.Port))
}
