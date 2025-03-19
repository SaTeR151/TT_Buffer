package handlers

import (
	"context"
	"net/http"

	"github.com/SaTeR151/TT_Buffer/internal/models"
	"github.com/SaTeR151/TT_Buffer/internal/service"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

func PostFactsToBuffer(s service.ServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var fact models.Fact
		if err := c.ShouldBindJSON(&fact); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			logger.Error(err)
			return
		}
		err := s.Insert(context.Background(), fact)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			logger.Error(err)
			return
		}
		c.Status(http.StatusOK)
	}
}
