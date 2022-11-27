package controller

import (
	"fmt"
	"time"

	"github.com/eazygood/getground-app/internal/errors"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

var (
	requestTimeout       = time.Second * 30
	customTimeLayoutList = []string{
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"02 Jan 2006 15:04:05 +0200",
		"2006-01-02 00:00:00 +0000",
		"02.01.2006 15:04:05",
		"02/01/2006 15:04:05",
		"02.01.2006 15:04:05",
		"2006/01/02 15:04:05",
	}
)

func logAndAbort(request *gin.Context, err errors.ApiError) {
	logger.Error(err.Message)
	request.AbortWithStatusJSON(err.Code, err)
}

func toTimePtr(dateValue string) (*time.Time, error) {
	for _, layout := range customTimeLayoutList {
		if t, err := time.Parse(layout, dateValue); err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("unable to find proper date")
}
