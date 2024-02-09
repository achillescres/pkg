package ginmiddleware

import (
	"encoding/json"
	"github.com/achillescres/pkg/cache/redisCache"
	"github.com/achillescres/pkg/hash"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

func NewCacherMiddleware(log *log.Entry, client *redis.Client, expiration time.Duration) (cacheGet gin.HandlerFunc, cacheSet gin.HandlerFunc) {

	cache := redisCache.NewRedisCacher(client)

	hasher := hash.NewMD5H()

	cacheGet = func(c *gin.Context) {
		url := c.Request.URL.String()
		bodyR := c.Request.Body

		body, err := io.ReadAll(bodyR)
		if err != nil {
			c.Next()
			return
		}

		hashKey := hasher.Hash([]byte(url), body)

		val, found := cache.Get(c, hashKey)
		if !found {
			c.Next()
		}
		c.Data(http.StatusOK, "application/json", val)
		c.Abort()
	}

	cacheSet = func(c *gin.Context) {
		val, found := c.Get("cache")
		if !found {
			c.Next()
			return
		}

		bodyR := c.Request.Body
		body, err := io.ReadAll(bodyR)
		if err != nil {
			c.Next()
			return
		}

		url := c.Request.URL.String()

		hashKey := hasher.Hash([]byte(url), body)

		data, err2 := json.Marshal(val)
		if err2 != nil {
			return
		}
		err = cache.Set(c, hashKey, data, expiration)
		if err != nil {
			log.Errorln(err)
			return
		}
	}

	return
}
