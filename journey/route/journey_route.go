package route

import (
	"me/coutcout/covoiturage/messaging"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	MSG_FILE_ACCEPTED = "File accepted"
)

func SetupRouter(logger *zap.SugaredLogger, router *gin.Engine) {
	logger.Debug("Creation of journey routes")
	router.POST("/import", func(c *gin.Context) {
		importJourney(c)
	})
}

func importJourney(c *gin.Context) {

	response := messaging.ResponseMessage{
		StatusCode: http.StatusAccepted,
		Message:    "File accepted",
		Errors:     []string{},
	}
	c.JSON(response.StatusCode, response)
}
