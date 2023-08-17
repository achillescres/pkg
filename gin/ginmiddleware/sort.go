package ginmiddleware

import (
	"fmt"
	"github.com/achillescres/pkg/gin/ginresponse"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	ASC     = "ASC"
	DESC    = "DESC"
	SortKey = "sortKey"
)

var ErrSortQueriesIsRequired = fmt.Errorf("error sort queries are required")

func SortMiddleware(defaultOrder, defaultBy string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("SortMiddleware")
		defaultOrder = strings.ToUpper(defaultOrder)

		by, ok := c.GetQuery("by")
		by = strings.ToUpper(by)
		if !ok || by == "" {
			by = defaultBy
		}

		order, ok := c.GetQuery("order")
		order = strings.ToUpper(order)
		if !ok || order == "" {
			order = defaultOrder
		} else if order != ASC && order != DESC {
			err := fmt.Errorf("error incorrect order: %s", order)
			ginresponse.ErrorString(c, http.StatusBadRequest, err, "invalid order")
			log.Errorln(err)
		}

		c.Set(SortKey, SortQueries{
			By:    by,
			Order: order,
		})
		c.Next()
	}
}

type SortQueries struct {
	By, Order string
}
