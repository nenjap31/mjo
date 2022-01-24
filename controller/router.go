package controller

import (
	"mjo/controller/user"
	"mjo/controller/merchant"
	"mjo/controller/outlet"
	"mjo/controller/transaction"
	"mjo/controller/auth"
	"mjo/controller/util/middleware/privilege"
	"net/http"

	"github.com/labstack/echo/v4"
)

const LOG_IDENTIFIER = "API_ROUTER"

type Controllers struct {
	User  user.Controller
	Merchant merchant.Controller
	Outlet outlet.Controller
	Transaction transaction.Controller
	AuthController        auth.Controller
	MiddlewarePrivilege   privilege.Privilege
}

func RegisterPath(
	e *echo.Echo,
	controllers Controllers,
) {
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	authV1 := e.Group("api/v1/auth")
	authV1.POST("/login", controllers.AuthController.Login)
	authV1.DELETE("/logout", controllers.AuthController.Logout, controllers.MiddlewarePrivilege.CheckPrivilege())

	userV1 := e.Group("api/v1/user")
	userV1.GET("", controllers.User.List, controllers.MiddlewarePrivilege.CheckPrivilege())
	userV1.POST("", controllers.User.Create, controllers.MiddlewarePrivilege.CheckPrivilege())
	userV1.GET("/:id", controllers.User.FindById, controllers.MiddlewarePrivilege.CheckPrivilege())
	userV1.PUT("/:id", controllers.User.UpdateById, controllers.MiddlewarePrivilege.CheckPrivilege())
	userV1.DELETE("/:id", controllers.User.DeleteById, controllers.MiddlewarePrivilege.CheckPrivilege())

	merchantV1 := e.Group("api/v1/merchant")
	merchantV1.GET("", controllers.Merchant.List, controllers.MiddlewarePrivilege.CheckPrivilege())
	merchantV1.POST("", controllers.Merchant.Create, controllers.MiddlewarePrivilege.CheckPrivilege())
	merchantV1.GET("/:id", controllers.Merchant.FindById, controllers.MiddlewarePrivilege.CheckPrivilege())
	merchantV1.PUT("/:id", controllers.Merchant.UpdateById, controllers.MiddlewarePrivilege.CheckPrivilege())
	merchantV1.DELETE("/:id", controllers.Merchant.DeleteById, controllers.MiddlewarePrivilege.CheckPrivilege())

	outletV1 := e.Group("api/v1/outlet")
	outletV1.GET("", controllers.Outlet.List, controllers.MiddlewarePrivilege.CheckPrivilege())
	outletV1.POST("", controllers.Outlet.Create, controllers.MiddlewarePrivilege.CheckPrivilege())
	outletV1.GET("/:id", controllers.Outlet.FindById, controllers.MiddlewarePrivilege.CheckPrivilege())
	outletV1.PUT("/:id", controllers.Outlet.UpdateById, controllers.MiddlewarePrivilege.CheckPrivilege())
	outletV1.DELETE("/:id", controllers.Outlet.DeleteById, controllers.MiddlewarePrivilege.CheckPrivilege())

	transactionV1 := e.Group("api/v1/transaction")
	transactionV1.GET("", controllers.Transaction.List, controllers.MiddlewarePrivilege.CheckPrivilege())
	transactionV1.POST("", controllers.Transaction.Create, controllers.MiddlewarePrivilege.CheckPrivilege())
	transactionV1.GET("/:id", controllers.Transaction.FindById, controllers.MiddlewarePrivilege.CheckPrivilege())
	transactionV1.PUT("/:id", controllers.Transaction.UpdateById, controllers.MiddlewarePrivilege.CheckPrivilege())
	transactionV1.DELETE("/:id", controllers.Transaction.DeleteById, controllers.MiddlewarePrivilege.CheckPrivilege())
	transactionV1.POST("/reportmerchant", controllers.Transaction.MonthlyReport, controllers.MiddlewarePrivilege.CheckPrivilege())
	transactionV1.POST("/reportoutlet", controllers.Transaction.MonthlyOutletReport, controllers.MiddlewarePrivilege.CheckPrivilege())
}
