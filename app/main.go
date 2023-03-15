package main

import (
	"me/coutcout/covoiturage/journey/repository"
	"me/coutcout/covoiturage/journey/router"
	"me/coutcout/covoiturage/journey/service"
	"me/coutcout/covoiturage/journey/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger zap.SugaredLogger

func main() {
	newLogger, _ := zap.NewDevelopment()
	logger = *newLogger.Sugar()
	r := gin.Default()

	// Repositories
	journeyRepo := repository.NewDbJourneyRepository(&logger)

	// Services
	journeyParser := service.NewJourneyCsvParser(&logger)

	// Usecases
	journeyUC := usecase.NewJourneyUsecase(
		&logger,
		journeyRepo,
		journeyParser,
	)

	router.NewJourneyRouter(&logger, r, journeyUC)

	r.Run()
}