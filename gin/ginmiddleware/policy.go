package ginmiddleware

import (
	"context"
	"errors"
	"github.com/achillescres/pkg/gin/ginresponse"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	policyInfoKey = "policyInfoKey_97asd8fguidso"
)

type TokenChecker[PolicyData any] func(ctx context.Context, token string) (PolicyData, error)

// UserPolicy returns a middleware that checks and validates auth cookie with TokenChecker
// PolicyData must be value type NOT POINTER!
func UserPolicy[PolicyData any](log logrus.Entry, cookieName string, check TokenChecker[PolicyData]) func(c *gin.Context) {
	return func(c *gin.Context) {
		access, err := c.Cookie(cookieName)
		if err != nil {
			log.Errorf("UserPolicy - (middleware): get access cookie: %s\n", err)
			ginresponse.ErrorString(c, http.StatusBadRequest, err, "access cookie is empty")
			return
		}

		policyInfo, err := check(c, access)
		if err != nil {
			log.Errorf("UserPolicy - (middleware): check token: %s\n", err)
			ginresponse.ErrorString(c, http.StatusUnauthorized, err, "token's permission check failed")
			return
		}

		c.Set(policyInfoKey, policyInfo)
		c.Next()
	}
}

// GetPolicyData returns PolicyData that was injected by UserPolicy middleware
func GetPolicyData[PolicyData any](ctx context.Context) (pd PolicyData, err error) {
	pd, ok := ctx.Value(policyInfoKey).(PolicyData)
	if !ok {
		return pd, errors.New("error no policy data in context")
	}
	return pd, nil
}
