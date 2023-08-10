package main

import (
	"fmt"
	"os"
	api "sport-space-api/api"
	"sport-space-api/config"
	"sport-space-api/docs"
	"sport-space-api/model"
	"sport-space-api/tools/email"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	config.Init()
	model.Init(config.DBCfg{})
	email.Init(config.MailCfg{})
}

func initRoute() {
	r := gin.Default()
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
		authorized.Use(api.AuthRequired())
		{
			authorized.GET("/", api.GetUser)
		}

	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}

func main() {
	initRoute()
}
