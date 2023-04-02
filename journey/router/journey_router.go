// Package router defines all the API path
package router

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/coutcout/covoiturage-csvreader/configuration"
	"github.com/coutcout/covoiturage-csvreader/domain"
	"github.com/coutcout/covoiturage-csvreader/messaging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type form struct {
	Files []*multipart.FileHeader `form:"files" binding:"required"`
}

type journeyRoute struct {
	logger         *zap.SugaredLogger
	journeyUsecase domain.JourneyUsecase
	cfg            *configuration.Config
}

// NewJourneyRouter creates a new router for journeys.
//
// @param logger - The logger to log to.
// @param cfg - The configuration of the application. Can be nil.
// @param mainRouter - The Gin Engine to add routes to.
// @param jUsecase - The domain.JourneyUsecase to use
func NewJourneyRouter(logger *zap.SugaredLogger, cfg *configuration.Config, mainRouter *gin.Engine, jUsecase domain.JourneyUsecase) {
	router := &journeyRoute{
		logger:         logger,
		cfg:            cfg,
		journeyUsecase: jUsecase,
	}
	logger.Debug("Creation of journey routes")
	mainURI := "/journey"
	mainRouter.POST(mainURI + "/import", func(c *gin.Context) {
		router.importJourney(c)
	})
	mainRouter.GET(mainURI, func(c *gin.Context){
		router.getJourneys(c)
	})
}

// importJourney imports files from a file upload
//
// @param j - route to respond to requests to import journeys
// @param c - gin. Context for request body to be passed
func (j *journeyRoute) importJourney(c *gin.Context) {
	var form form
	err := c.ShouldBind(&form)
	if err != nil {
		j.logger.Errorw("Error importing file",
			"error", err.Error(),
		)
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, messaging.SingleResponseMessage{
			Errors: []string{"'files' parameter is required"},
		})
		return
	}

	maxUploadFileSize := j.cfg.Journey.Import.MaxUploadFile * 1024
	response := messaging.MultipleResponseMessage{
		Files: []messaging.FileImportResponseMessage{},
		Data: messaging.FileImportData{
			TotalFilesImported: len(form.Files),
		},
	}

	for _, formFile := range form.Files {

		fileResponse := messaging.FileImportResponseMessage{
			Filename: formFile.Filename,
			Imported: true,
			Errors:   []string{},
		}

		if formFile.Size > maxUploadFileSize {
			err := fmt.Errorf("file %s is too big (current: %d - max: %d)", formFile.Filename, formFile.Size, maxUploadFileSize)
			c.Error(err)
			response.Data.NbFilesWithErrors++
			fileResponse.Imported = false
			fileResponse.Errors = append(fileResponse.Errors, err.Error())
			response.Files = append(response.Files, fileResponse)
			break
		}

		openedFile, err := formFile.Open()
		if err != nil {
			j.logger.Errorw("Error importing file",
				"error", err.Error(),
				"filename", formFile.Filename,
				"filesize", formFile.Size,
			)
			c.Error(err)
			response.Data.NbFilesWithErrors++
			fileResponse.Imported = false
			fileResponse.Errors = append(fileResponse.Errors, err.Error())
			break
		}

		nbLineImported, errors := j.journeyUsecase.ImportFromCSVFile(c, openedFile)
		response.Data.NbLineImported += int(nbLineImported)
		fileResponse.Errors = append(fileResponse.Errors, errors...)
		fileResponse.NbLineImported = int(nbLineImported)
		if len(errors) > 0 {
			response.Data.NbFilesWithErrors++
			if nbLineImported == 0 {
				fileResponse.Imported = false
			}
		} else {

			response.Data.NbFilesSucceded++
		}
		response.Files = append(response.Files, fileResponse)
	}

	responseStatus := http.StatusAccepted
	if response.Data.NbFilesWithErrors > 0 && response.Data.NbLineImported == 0 {
		responseStatus = http.StatusBadRequest
	}

	c.JSON(responseStatus, response)
}

func (r *journeyRoute) getJourneys(c *gin.Context){
	errorChan := make(chan error)
	journeyChan := r.journeyUsecase.GetJourneys(c, errorChan);
	defer close(errorChan)
	go func(c *gin.Context){
		err, ok := <- errorChan
		if !ok {
			r.logger.Debug("getJourneys ended without error")
			return
		}

		r.logger.Errorw("Error while retrieving all journeys with streaming",
			"error", err,
		)
		c.Abort()
	}(c)

	c.Stream(func(w io.Writer) bool {
		msg, ok := <- journeyChan
		if !ok {
			return false
		}
		encoder := json.NewEncoder(c.Writer)
		encoder.Encode(msg)
		return true
	})
	
}
