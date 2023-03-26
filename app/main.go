// Package main contains the main file
package main

import (
	"log"
	"os"

	"github.com/coutcout/covoiturage-csvReader/configuration"
	"github.com/coutcout/covoiturage-csvReader/journey/repo"
	"github.com/coutcout/covoiturage-csvReader/journey/router"
	"github.com/coutcout/covoiturage-csvReader/journey/service"
	"github.com/coutcout/covoiturage-csvReader/journey/usecase"

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
	journeyRepo := repo.NewDbJourneyRepository(
		&logger,
		cfg,
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
