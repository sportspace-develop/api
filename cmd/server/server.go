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
	gormsessions "github.com/gin-contrib/sessions/gorm"
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
	r.Use(gin.Recovery())
	// r.Use(api.TimeoutMiddleware(2 * time.Second))

	db, err := model.Connect()
	if err != nil {
		panic(err)
	}
	store := gormsessions.NewStore(db, true, []byte(config.App.CookieSecret))
	r.Use(sessions.Sessions("sport-space-session", store))

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		auth.Use(func(c *gin.Context) {})
		{
			auth.POST("/otp", api.GetAuthCode)
			auth.POST("/login", api.Authorize)
			auth.POST("/refresh", api.Refresh)
			auth.POST("/logout", api.Logout).Use(api.AuthRequiredMiddleware())
		}

		profile := v1.Group("/profile")
		profile.Use(api.AuthRequiredMiddleware())
		{
			profile.GET("/", api.GetProfile)

			profile.POST("/setPassword", api.SetPassword)

			profile.GET("/organization", api.GetOrganization)
			profile.POST("/organization", api.CreateOrganization)
			profile.PUT("/organization", api.UpdateOrganization)

			profile.GET("/tournament", api.GetTournament)
			profile.POST("/tournament", api.CreateTournament)
			profile.PUT("/tournament", api.UpdateTournament)

			profile.GET("/team", api.GetTeam)
			profile.POST("/team", api.CreateTeam)
			profile.PUT("/team", api.UpdateTeam)

			profile.GET("/team/invite", api.GetInviteToTeam)
			profile.POST("/team/invite", api.CreateInviteToTeam)
			profile.PUT("/team/invite", api.UpdateInviteToTeam)

			profile.GET("/player", api.GetPlayer)
			profile.PUT("/player", api.UpdatePlayer)

			profile.GET("/player/invite", api.GetPlayerInvite)
			profile.PUT("/player/invite", api.UpdatePlayerInvite)
		}

		guest := v1.Group("/")
		{
			guest.GET("/home", api.GetHome) // - Общая информация для пользователя + справочники
			guest.GET("/organization", api.GetAllOrganization)
			guest.GET("/tournament", api.GetTournaments)
			guest.GET("/team", api.GetTeams)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}

func main() {
	log.INFO("Start app")
	initRoute()
}
