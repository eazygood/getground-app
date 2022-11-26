package controller

import (
	"time"

	"github.com/eazygood/getground-app/internal/errors"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

var (
	requestTimeout = time.Second * 30
)

func logAndAbort(request *gin.Context, err errors.ApiError) {
	logger.Error(err.Message)
	request.AbortWithStatusJSON(err.Code, err)
}
