package main

import (
	"fmt"
	"os"
	api "sport-space-api/api"
	"sport-space-api/config"
	"sport-space-api/docs"
	"sport-space-api/logger"
	"sport-space-api/model"
	"sport-space-api/tools/email"
	"sport-space-api/tools/jwt"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	log *logger.Logger
)

func init() {
	log = logger.New("main")
	log.INFO("init app")
	config.Init()
	api.Init()
	model.Init(config.DBCfg{})
	email.Init(config.MailCfg{})

	jwt.Secret = []byte(config.App.JWTSecret)
	jwt.AccessTokenLongTime = time.Duration(config.App.JWTLongTime * int(time.Minute))
}

func initRoute() {
	r := gin.Default()
	r.Use(api.LoggingMiddleware())
	r.Use(api.CORSMiddleware())
	store := cookie.NewStore([]byte(config.App.CookieSecret))
	r.Use(sessions.Sessions("sport-space", store))

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/auth")
		eg.Use(func(c *gin.Context) {})
		{
			eg.POST("/otp", api.GetAuthCode)
			eg.POST("/login", api.Authorize)
			eg.POST("/refresh", api.Refresh)
		}
		authorized := v1.Group("/user")
		authorized.Use(api.AuthRequiredMiddleware())
		{
			authorized.GET("/", api.GetUser)
		}

	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}

func main() {
	log.INFO("Start app")
	initRoute()
}
