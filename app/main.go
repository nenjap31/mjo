package main

import (
	"context"
	//"fmt"
	"mjo/config"
	router "mjo/controller"
	authController "mjo/controller/auth"
	userController "mjo/controller/user"
	merchantController "mjo/controller/merchant"
	outletController "mjo/controller/outlet"
	transactionController "mjo/controller/transaction"
	customValidator "mjo/controller/util/validator"
	"mjo/repository/mysql"
	userRepository "mjo/repository/user"
	merchantRepository "mjo/repository/merchant"
	outletRepository "mjo/repository/outlet"
	transactionRepository "mjo/repository/transaction"
	authService "mjo/service/auth"
	userService "mjo/service/user"
	merchantService "mjo/service/merchant"
	outletService "mjo/service/outlet"
	transactionService "mjo/service/transaction"
	cache "mjo/repository/util/redis"
	middlewarePrivilege "mjo/controller/util/middleware/privilege"
	"mjo/util/gorm"
	"mjo/util/logger"
	"os"
	"os/signal"
	"time"
	//"github.com/labstack/echo/v4/middleware"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	logRotator "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const LOG_IDENTIFIER = "APP_MAIN"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func CustomValidator() *validator.Validate {
	customValidator := validator.New()
	return customValidator
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})
	logrus.SetLevel(logrus.InfoLevel)

	logFile, err := logRotator.New(
		"resource/log/log.%Y%m%d%H%M",
		logRotator.WithLinkName("resource/log/log"),
		logRotator.WithMaxAge(time.Duration(3600)*time.Second),
		logRotator.WithRotationTime(time.Duration(3600)*time.Second),
	)
	if err == nil {
		logrus.SetOutput(logFile)
	} else {
		log.Error(err)
		panic(err)
	}
	defer logFile.Close()

	dbCon, err := mysql.Connect(config.GetConfig().MysqlDb)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	gormDb, err := gorm.InitGorm(dbCon)
	if err != nil {
		dbCon.Close()
		log.Fatal(err)
		panic(err)
	}

	redisUrl := config.GetConfig().RedisHost + ":" + config.GetConfig().RedisPort
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     redisUrl,
			Password: config.GetConfig().RedisPassword,
			DB:       config.GetConfig().RedisDb,
		})

	user, err := userRepository.NewRepository(gormDb)
	if err != nil {
		log.Error(err)
	}
	merchant, err := merchantRepository.NewRepository(gormDb)
	if err != nil {
		log.Error(err)
	}
	outlet, err := outletRepository.NewRepository(gormDb)
	if err != nil {
		log.Error(err)
	}
	transaction, err := transactionRepository.NewRepository(gormDb)
	if err != nil {
		log.Error(err)
	}

	cacheServer := cache.NewRedisRepository(redisClient)

	privilegeMiddleware := middlewarePrivilege.New(cacheServer)

	authService := authService.NewService(user, cacheServer)
	authController := authController.NewController(authService)

	userService := userService.NewService(user)
	userController := userController.NewController(userService)
	merchantService := merchantService.NewService(merchant)
	merchantController := merchantController.NewController(merchantService)	
	outletService := outletService.NewService(outlet)
	outletController := outletController.NewController(outletService)
	transactionService := transactionService.NewService(transaction)
	transactionController := transactionController.NewController(transactionService)	

	controllers := router.Controllers{
		MiddlewarePrivilege:   privilegeMiddleware,
		User:  *userController,
		Merchant:  *merchantController,
		Outlet:  *outletController,
		Transaction:  *transactionController,
		AuthController:  *authController,
	}

	e := echo.New()
	e.Validator = &customValidator.BodyRequestValidator{Validator: CustomValidator()}
	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"*"},
	//}))
	router.RegisterPath(e, controllers)

	go func() {
		if err := e.Start(":2022"); err != nil {
			dbCon.Close()
			log.Fatal(err)
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		dbCon.Close()
		log.Fatal(err)
		panic(err)
	}
}