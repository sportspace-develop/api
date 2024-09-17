package rest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"sport-space/docs"
	"sport-space/internal/adapter/models"
	"sport-space/internal/adapter/storage/errstore"
	"sport-space/internal/core/sportspace"
	"sport-space/pkg/jwt"
	"sport-space/pkg/tools"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "sport-space/docs"
)

var (
	msgErrorCloseBody = "failed close body request"
)

type sport interface {
	NewOTP(ctx context.Context, email string) error
	LoginWithOTP(ctx context.Context, email, otp string) (*models.User, error)
	GetAllTournaments(ctx context.Context) (*[]models.Tournament, error)
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
	NewTournament(ctx context.Context, tournament *models.Tournament) (*models.Tournament, error)
	GetTournaments(ctx context.Context, user *models.User) (*[]models.Tournament, error)
	GetTournamentByID(ctx context.Context, tournamentID uint) (*models.Tournament, error)
	UpdTournament(ctx context.Context, tournament *models.Tournament) (*models.Tournament, error)
	NewTeam(ctx context.Context, team *models.Team) (*models.Team, error)
	GetTeams(ctx context.Context, user *models.User) (*[]models.Team, error)
	GetTeamByID(ctx context.Context, teamID uint) (*models.Team, error)
	UpdTeam(ctx context.Context, team *models.Team, playersIDs *[]uint) (*models.Team, *[]models.Player, error)
	NewPlayer(ctx context.Context, player *models.Player) (*models.Player, error)
	GetPlayers(ctx context.Context, userID uint) (*[]models.Player, error)
	GetPlayerByID(ctx context.Context, playerID uint) (*models.Player, error)
	UpdPlayer(ctx context.Context, player *models.Player) (*models.Player, error)
	NewApplicationTeam(ctx context.Context, playerIDs *[]uint, tournamentID, teamID, userID uint) (
		*models.Application, *[]models.Player, error,
	)
	UpdApplicationTeam(ctx context.Context, applicationID uint, playerIDs *[]uint, status models.ApplicationStatus, teamID uint, userID uint) (
		*models.Application, *[]models.Player, error,
	)
	GetApplicationsTeam(ctx context.Context, teamID uint) (*[]models.Application, error)
	GetApplicationByID(ctx context.Context, applicationID uint) (*models.Application, error)
	GetPlayersFromApplication(ctx context.Context, applicationID uint) (*[]models.Player, error)
	GetApplicationsFromTournament(ctx context.Context, tournamentID uint) (*[]models.Application, error)
	UpdApplicationTournament(ctx context.Context, applicationID uint, status models.ApplicationStatus, tournamentID uint, userID uint) (*models.Application, error)
}

type Server struct {
	srv           *http.Server
	log           *zap.Logger
	sport         sport
	secret        string
	uploadPath    string
	uploadURL     string
	uploadMaxSize int64
	baseURL       string
}

type option func(s *Server)

func SetLogger(l *zap.Logger) option {
	return func(s *Server) {
		s.log = l
	}
}

func SetAddress(address string) option {
	return func(s *Server) {
		s.srv.Addr = address
	}
}

func SetSecretKey(secret string) option {
	return func(s *Server) {
		s.secret = secret
	}
}

func SetUploadPath(path string) option {
	return func(s *Server) {
		s.uploadPath = path
	}
}

func SetBaseURL(url string) option {
	return func(s *Server) {
		s.baseURL = url
	}
}

func New(service sport, options ...option) (*Server, error) {
	s := &Server{
		srv:           &http.Server{},
		log:           zap.NewNop(),
		sport:         service,
		uploadPath:    "./uploads", // Directory upload.
		uploadURL:     "/uploads",  // URL upload.
		uploadMaxSize: 2 << 20,
		secret:        "",
		baseURL:       "",
	}

	for _, opt := range options {
		opt(s)
	}

	return s, nil
}

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func (s *Server) Run() error {
	// swagger info.
	docs.SwaggerInfo.Host = s.baseURL
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}
	docs.SwaggerInfo.Title = "SportSpace API"
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.Description = "sport-space api documentation"

	r := gin.New()
	r.MaxMultipartMemory = s.uploadMaxSize
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))
	r.Use(s.middlewareLogger())
	r.GET("/ping", s.handlerPing)
	r.Static(s.uploadURL, s.uploadPath)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/otp", s.handlerAuthOTP)
			auth.POST("/login", s.handlerLogin)
			auth.GET("/logout", s.handlerLogout)
		}
		user := api.Group("/user")
		user.Use(s.middlewareAuthentication())
		{
			user.GET("/", s.handlerUser)
			user.POST("/tournaments", s.handlerUserNewTournament)
			user.GET("/tournaments", s.handlerUserTournaments)
			user.GET("/tournaments/:id", s.handlerUserTournament)
			user.PUT("/tournaments/:id", s.handlerUserUpdTournament)
			user.PUT("/tournaments/:id/upload", s.handlerUserUploadTournament)

			user.POST("/teams", s.handlerUserNewTeam)
			user.GET("/teams", s.handlerUserTeams)
			user.GET("/teams/:id", s.handlerUserTeam)
			user.PUT("/teams/:id", s.handlerUserUptTeam)

			user.POST("/players", s.handlerUserNewPlayer)
			user.GET("/players", s.handlerUserPlayers)
			user.PUT("/players/:id", s.handlerUserUpdatePlayer)
			user.PUT("/players/:id/upload", s.handlerUserUploadPlayer)

			// заявки турнира
			user.GET("/tournaments/:id/applications", s.handlerGetTournamentApplications)
			user.GET("/tournaments/:id/applications/:aid", s.handlerGetTournamentApplication)
			user.PUT("/tournaments/:id/applications/:aid", s.handlerUpdTournamentApplication)

			// заявки команды
			user.POST("/teams/:id/applications", s.handlerNewTeamApplication)
			user.PUT("/teams/:id/applications/:aid", s.handlerUpdStatusTeamApplication)
			user.GET("/teams/:id/applications", s.handlerGetTeamApplications)
			user.GET("/teams/:id/applications/:aid", s.handlerGetApplication)
		}

		guest := api.Group("/")
		{
			guest.GET("/tournaments", s.handlerGetAllTournament)
		}

	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.srv.Handler = r.Handler()

	if err := s.srv.ListenAndServe(); err != nil {
		return fmt.Errorf("server stopped with error: %w", err)
	}

	return nil
}

func unauthorize(c *gin.Context) {
	userCookie := &http.Cookie{
		Name:  cookieName,
		Value: "",
		Path:  "/",
	}
	c.Request.AddCookie(userCookie)
	http.SetCookie(c.Writer, userCookie)
}

func (s *Server) authorization(c *gin.Context, login, password string) error {
	var err error
	var user *models.User
	ctx := c.Request.Context()
	if user, err = s.sport.LoginWithOTP(ctx, login, password); err != nil {
		return fmt.Errorf("failed authorization: %w", err)
	}

	jwtRest := jwt.New([]byte(s.secret))
	signedCookie, err := jwtRest.Create(cookieKey, strconv.Itoa(int(user.ID)))
	if err != nil {
		return fmt.Errorf("can't create cookie data: %w", err)
	}

	userCookie := &http.Cookie{
		Name:  cookieName,
		Value: signedCookie,
		Path:  "/",
	}
	c.Request.AddCookie(userCookie)
	http.SetCookie(c.Writer, userCookie)

	return nil
}

func (s *Server) readBody(c *gin.Context) ([]byte, int) {
	bBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.log.Error("failed read body", zap.Error(err))
		return []byte{}, http.StatusInternalServerError
	}
	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			s.log.Error(msgErrorCloseBody, zap.Error(err))
		}
	}()
	return bBody, 0
}

func (s *Server) login(c *gin.Context, login, password string) (int, string) {
	if err := s.authorization(c, login, password); err != nil {
		if errors.Is(err, sportspace.ErrLoginNotValid) || errors.Is(err, sportspace.ErrPasswordNotValid) {
			if errors.Is(err, sportspace.ErrLoginNotValid) {
				return http.StatusBadRequest, "Не верный формат логина"
			}
			if errors.Is(err, sportspace.ErrPasswordNotValid) {
				return http.StatusBadRequest, "Не верный формат пароля"
			}
			return http.StatusInternalServerError, ""
		}
		if errors.Is(err, sportspace.ErrPasswordNotEquale) || errors.Is(err, errstore.ErrNotFoundData) {
			return http.StatusUnauthorized, ""
		}
		s.log.Error("authorization failed", zap.Error(err))
		return http.StatusInternalServerError, ""
	}
	return http.StatusOK, ""
}

func (s *Server) genUploadName(name string) string {
	return "/" + tools.RandomString(20) + filepath.Ext(name)
}

func (s *Server) saveFile(file *multipart.FileHeader, dst string) error {
	return tools.SaveUploadedFile(file, s.uploadPath+dst)
}

func (s *Server) removeFile(dst string) error {
	return os.Remove(s.uploadPath + dst)
}

func (s Server) getFullUploadPath(name string) string {
	if name == "" {
		return ""
	}
	return s.uploadURL + name
}

func (s Server) getFullUploadURL(name string) string {
	if name == "" {
		return ""
	}
	return s.baseURL + s.uploadURL + name
}

func (s Server) isValidImgExtension(file *multipart.FileHeader) bool {
	return strings.HasSuffix(file.Filename, ".png") ||
		strings.HasSuffix(file.Filename, ".jpg") ||
		strings.HasSuffix(file.Filename, ".jpeg") ||
		strings.HasSuffix(file.Filename, ".gif")
}
