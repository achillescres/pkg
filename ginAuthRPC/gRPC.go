package ginAuthRPC

import (
	"context"
	"fmt"
	authProto2 "github.com/achillescres/pkg/gin/ginmiddleware/authProto"
	"github.com/achillescres/pkg/gin/ginresponse"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

type mode struct{}

func (mode) Local() string {
	return "local"
}

func (mode) Dev() string {
	return "dev"
}

func (mode) Prod() string {
	return "prod"
}

var Mode = mode{}

type AuthServer interface {
	AuthMiddleware(c *gin.Context)
}

type authServerRPC struct {
	externalAuthClient authProto2.ExternalAuthClient
}

const (
	envLocalFilename = ".env.local"
	envDevFilename   = ".env.dev"
	envProdFilename  = ".env.prod"
)

func ConfigRPCAuth(mode string) error {
	envFilename := ""
	switch mode {
	case Mode.Local():
		envFilename = envLocalFilename
	case Mode.Dev():
		envFilename = envDevFilename
	case Mode.Prod():
		envFilename = envProdFilename
	default:
		panic("invalid mode")
	}

	err := godotenv.Load(envFilename)
	if err != nil {
		return err
	}
	err = viper.BindEnv("gRPCAuthAddr", "GRPC_AUTH_ADDR")
	if err != nil {
		return err
	}
	return nil
}

func NewRPCServer(ctx context.Context, opt ...interface{}) (AuthServer, error) {
	addr := viper.GetString("gRPCAuthAddr")
	fmt.Println(addr)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := authProto2.NewExternalAuthClient(conn)

	return &authServerRPC{
		externalAuthClient: c,
	}, nil
}

const (
	userKey         = "userKey"
	userLoginKey    = "userLoginKey"
	restrictionsKey = "restrictions"
)

func (a *authServerRPC) AuthMiddleware(c *gin.Context) {
	cookie, err := c.Cookie("Access")
	if err != nil {
		log.Errorln(err)
		ginresponse.ErrorString(c, http.StatusUnauthorized, err, "cookies are missing")
		return
	}
	userInfo, err := a.externalAuthClient.Permissions(c, &authProto2.CookieAccess{Access: cookie})
	if err != nil {
		log.Errorln(err)
		ginresponse.ErrorString(c, http.StatusUnauthorized, err, err.Error())
		return
	}

	if userInfo == nil {
		log.Errorln(err)
		ginresponse.ErrorString(c, http.StatusUnauthorized, err, err.Error())
		return
	}

	c.Set(userKey, userInfo.User)
	c.Set(userLoginKey, userInfo.UserLogin)
	c.Set(restrictionsKey, newRestrictions(userInfo.AirlCode, userRole(userInfo.UserRole)))

	c.Next()
}

type restrictions struct {
	AirlCode string
	UserRole userRole
}

func newRestrictions(airline string, userRole userRole) *restrictions {
	return &restrictions{AirlCode: airline, UserRole: userRole}
}

type userRole int

var userRoles = map[string]userRole{
	"tester":    -1,
	"manager":   0,
	"admin":     1,
	"developer": 2,
}

const (
	userRoleTester = -1 + iota
	userRoleManager
	userRoleAdmin
	userRoleDeveloper
)

func (r userRole) HigherOrEqualThan(or userRole) bool {
	return r >= or
}

func (r userRole) LowerThan(or userRole) bool {
	return r < or
}

func (r userRole) Equal(or userRole) bool {
	return r == or
}
