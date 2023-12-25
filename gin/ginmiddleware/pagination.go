package ginmiddleware

import (
	"fmt"
	"github.com/achillescres/pkg/gin/ginresponse"
	"github.com/achillescres/pkg/vars"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

const (
	PaginationKey = "paginationKey"
)

var ErrPaginationQueriesIsRequired = fmt.Errorf("error paginations queries is required")

func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		var quoteInt int
		quoteS, ok := c.GetQuery("dayQuote")
		if !ok || quoteS == "" {
			ginresponse.ErrorString(c, http.StatusBadRequest, ErrPaginationQueriesIsRequired, "dayQuote is required")
			log.Errorln(ErrPaginationQueriesIsRequired)
			return
		}
		quoteInt, err = strconv.Atoi(quoteS)
		if err != nil || quoteInt <= 0 {
			ginresponse.ErrorString(c, http.StatusBadRequest, err, "invalid dayQuote, it is natural number")
			log.Errorln(err)
			return
		}

		dayQuote := uint64(quoteInt)

		lastSeenFltDate, okFltDate := c.GetQuery("lastSeenFltDate")
		if !okFltDate || lastSeenFltDate == "" {
			lastSeenFltDate = "00010101"
		} else if _, err := time.Parse(vars.FltDateLayout, lastSeenFltDate); err != nil {
			err := fmt.Errorf("invalid lastSeenFltDate query, must be YYYYMMDD format: %w", err)
			ginresponse.ErrorString(c, http.StatusBadRequest, err, "invalid lastSeenFltDate, format is YYYYMMDD")
			log.Errorln(err)
			return
		}
		c.Set(PaginationKey, PaginationQueries{
			DayQuote:        dayQuote,
			LastSeenFltDate: lastSeenFltDate,
		})
		c.Next()
	}
}

type PaginationQueries struct {
	DayQuote        uint64
	LastSeenFltDate string
}
