package router

import (
	"me/coutcout/covoiturage/domain"
	"me/coutcout/covoiturage/messaging"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	MSG_FILE_ACCEPTED = "File accepted"
)

type Form struct {
	Files []*multipart.FileHeader `form:"files" binding:"required"`
}

type journeyRoute struct {
	logger         *zap.SugaredLogger
	journeyUsecase domain.JourneyUsecase
}

func NewJourneyRouter(logger *zap.SugaredLogger, mainRouter *gin.Engine, jUsecase domain.JourneyUsecase) {
	router := &journeyRoute{
		logger:         logger,
		journeyUsecase: jUsecase,
	}

	logger.Debug("Creation of journey routes")
	mainRouter.POST("/import", func(c *gin.Context) {
		router.importJourney(c)
	})
}

func (j *journeyRoute) importJourney(c *gin.Context) {
	var form Form
	err := c.ShouldBind(&form)
	if err != nil {
		j.logger.Errorw("Error importing file",
			"error", err.Error(),
		)
		c.Error(err)
	}

	for _, formFile := range form.Files {
		openedFile, err := formFile.Open()
		if err != nil {
			j.logger.Errorw("Error importing file",
				"error", err.Error(),
				"filename", formFile.Filename,
				"filesize", formFile.Size,
			)
			c.Error(err)
			break
		}

		j.journeyUsecase.ImportFromCSVFile(openedFile)
	}

	response := messaging.ResponseMessage{
		StatusCode: http.StatusAccepted,
		Message:    "File accepted",
		Errors:     c.Errors.Errors(),
	}
	c.JSON(response.StatusCode, response)
}
