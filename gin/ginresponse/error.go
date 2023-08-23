package ginresponse

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type errorJSON struct {
	Error string `json:"error"`
}

func ErrorString(c *gin.Context, code int, err error, clientSideError string) {
	if err == nil {
		return
	}
	c.AbortWithStatusJSON(code, errorJSON{Error: err.Error()})
	c.Error(err)
}

func JSON(c *gin.Context, code int, json gin.H) {
	c.JSON(code, json)
	if err := c.Err(); err != nil {
		log.Errorf("error occured in gin context: %s\n", err)
	}
}

func MissingRequiredQuery(c *gin.Context, queryName string) error {
	err := fmt.Errorf("error %s query is required", queryName)
	ErrorString(c, http.StatusBadRequest, err, err.Error())
	return err
}

func MissingRequiredFormKey(c *gin.Context, queryName string) error {
	err := fmt.Errorf("error %s form key is required", queryName)
	ErrorString(c, http.StatusBadRequest, err, err.Error())
	return err
}
