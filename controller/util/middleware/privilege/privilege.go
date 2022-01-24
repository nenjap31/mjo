package privilege

import (
	"mjo/controller/util/response"
	authServiceDefined "mjo/service/auth/defined"
	"mjo/config"
	"mjo/repository/util/redis"
	"errors"
	"net/http"
	"reflect"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

type ParamCheckPrivilege struct {
	Submenu string
	Access  []int
}

const (
	invalidJWT = "Invalid Token"
	LOG_PREFIX = "[PRIVILEGE MIDDLEWARE]"
)

var mapping = map[string]int{
	"GET":    1,
	"POST":   2,
	"PUT":    2,
	"DELETE": 3,
}

type Privilege struct {
	redis             redis.IRepository
}

func New(redis redis.IRepository) Privilege {
	return Privilege{
		redis:             redis,
	}
}
func (priv *Privilege) CheckPrivilege() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessToken := c.Request().Header.Get("Authorization")
			authScheme := "Bearer"
			l := len(authScheme)
			if len(accessToken) > l+1 && accessToken[:l] == authScheme {
				accessToken = accessToken[l+1:]
			}
			t := reflect.ValueOf(&authServiceDefined.AccessTokenClaims{}).Type().Elem()
			claims := reflect.New(t).Interface().(jwt.Claims)
			token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
				if _, success := token.Method.(*jwt.SigningMethodHMAC); !success {
					return nil, errors.New("internal server error")
				}
				return []byte(config.GetConfig().SecretKey), nil
			})
			if err != nil || !token.Valid {
				log.Info(LOG_PREFIX, " Token invalid")
				return c.JSON(http.StatusBadRequest, response.NewResponse("", response.Map["badRequest"], invalidJWT))
			}
			claimsData := token.Claims.(*authServiceDefined.AccessTokenClaims)

			jwtTokenRedis := claimsData.Username + "|" + accessToken
			key, _ := priv.redis.Get(jwtTokenRedis)
			if key == nil {
				return c.JSON(http.StatusUnauthorized, response.NewResponse("", response.Map["unauthorized"], invalidJWT))
			}
			c.Set("user", token)
			return next(c)
		}
	}
}
