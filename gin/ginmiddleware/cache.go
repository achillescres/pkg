package ginmiddleware

import (
	"encoding/json"
	"github.com/achillescres/pkg/cache/redisCache"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func NewCacherMiddleware(log *log.Entry, client *redis.Client, expiration time.Duration) (cacheGet gin.HandlerFunc, cacheSet gin.HandlerFunc) {

	cache := redisCache.NewRedisCacher(client)

	cacheGet = func(c *gin.Context) {
		url := c.Request.URL.String()
		val, found := cache.Get(c, url)
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
		url := c.Request.URL.String()
		data, err2 := json.Marshal(val)
		if err2 != nil {
			return
		}
		err := cache.Set(c, url, data, expiration)
		if err != nil {
			log.Errorln(err)
			return
		}
	}

	return
}
