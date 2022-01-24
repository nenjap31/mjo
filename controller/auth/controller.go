package auth

import (
	controllerAuthDefined "mjo/controller/auth/defined"
	"mjo/controller/util/response"
	"mjo/service/auth"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

const LOGPREFIX = "[AUTH CONTROLLER]"

type Controller struct {
	service auth.IService
}

func NewController(service auth.IService) *Controller {
	return &Controller{service}
}

func (controller *Controller) Login(c echo.Context) error {
	bodyRequest := new(controllerAuthDefined.LoginRequest)
	if err := c.Bind(bodyRequest); err != nil {
		log.Error(LOGPREFIX, " Login - Bind error - ", err.Error())
		return c.JSON(http.StatusBadRequest, response.NewResponse("", response.Map["badRequest"], err.Error()))
	}
	if err := c.Validate(bodyRequest); err != nil {
		log.Error(LOGPREFIX, " Login - Validation error - ", err.Error())
		return c.JSON(http.StatusBadRequest, response.NewResponse("", response.Map["badRequest"], err.Error()))
	}
	data, err := controller.service.GenerateToken(bodyRequest.Username, bodyRequest.Password)
	if err != nil {
		
		return c.JSON(http.StatusBadRequest, response.NewResponse("", response.Map["badRequest"], err.Error()))
	}
	log.Info(LOGPREFIX, " Login - Login successfull")
	return c.JSON(http.StatusCreated, response.NewResponse("", response.Map["ok"], data))
}

func (controller *Controller) Logout(c echo.Context) error {
	reqToken := c.Get("user").(*jwt.Token)
	err := controller.service.RevokeToken(reqToken)
	if err != nil {
		log.Error(LOGPREFIX, " Logout - Bad request - ", err.Error())
		return c.JSON(http.StatusBadRequest, response.NewResponse("", response.Map["badRequest"], err.Error()))
	}
	log.Info(LOGPREFIX, " Logout - Logout successfull")
	return c.JSON(http.StatusCreated, response.NewResponse("", response.Map["created"], nil))
}

