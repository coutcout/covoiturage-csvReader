package router

import (
	"fmt"
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

	const MAX_UPLOAD_FILE = 1024 * 1024
	response := messaging.MultipleResponseMessage{
		Files: []messaging.FileImportResponseMessage{},
		Data: messaging.FileImportData{
			TotalImported: len(form.Files),
		},
	}

	for _, formFile := range form.Files {

		fileResponse := messaging.FileImportResponseMessage{
			Filename: formFile.Filename,
			Imported: true,
			Errors: []string{},
		}

		if formFile.Size > MAX_UPLOAD_FILE {
			err := fmt.Errorf("file %s is too big", formFile.Filename) 
			c.Error(err)
			response.Data.NbErrors += 1
			fileResponse.Imported = false
			fileResponse.Errors = append(fileResponse.Errors, err.Error())	
			response.Files = append(response.Files, fileResponse)
			break;
		}

		openedFile, err := formFile.Open()
		if err != nil {
			j.logger.Errorw("Error importing file",
				"error", err.Error(),
				"filename", formFile.Filename,
				"filesize", formFile.Size,
			)
			c.Error(err)
			response.Data.NbErrors += 1
			fileResponse.Imported = false
			fileResponse.Errors = append(fileResponse.Errors, err.Error())	
			break
		}

		j.journeyUsecase.ImportFromCSVFile(openedFile)
		response.Data.NbSucceded += 1
	}

	responseStatus := http.StatusAccepted
	if response.Data.NbErrors > 0 {
		responseStatus = http.StatusBadRequest
	}

	c.JSON(responseStatus, response)
}

