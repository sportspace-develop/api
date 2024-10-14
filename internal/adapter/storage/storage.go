package storage

import (
	"context"
	"errors"

	"sport-space/internal/adapter/models"
	"sport-space/internal/adapter/storage/database"
)

type Store interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
	GetAllTournaments(ctx context.Context) (*[]models.Tournament, error)
	NewUser(ctx context.Context, login, email, passwordHash string) (*models.User, error)
	NewOTP(ctx context.Context, otp *models.OTPUser) error
	GetOTP(ctx context.Context, user *models.User) (*models.OTPUser, error)
	RemoveOTP(ctx context.Context, user *models.User) error
	NewTournament(ctx context.Context, tournament *models.Tournament) (*models.Tournament, error)
	GetTournaments(ctx context.Context, userID uint) (*[]models.Tournament, error)
	GetTournamentByID(ctx context.Context, tournamentID uint) (*models.Tournament, error)
	UpdTournamentByUser(ctx context.Context, tournament *models.Tournament) (*models.Tournament, error)
	NewTeam(ctx context.Context, team *models.Team) (*models.Team, error)
	GetTeams(ctx context.Context, user *models.User) (*[]models.Team, error)
	GetTeamByID(ctx context.Context, teamID uint) (*models.Team, error)
	UpdTeam(ctx context.Context, team *models.Team, playersIDs *[]uint) (*models.Team, *[]models.Player, error)
	NewPlayer(ctx context.Context, player *models.Player) (*models.Player, error)
	NewPlayerBatch(ctx context.Context, players *[]models.Player) (*[]models.Player, error)
	GetPlayers(ctx context.Context, userID uint) (*[]models.Player, error)
	GetPlayerByID(ctx context.Context, playerID uint) (*models.Player, error)
	GetPlayersByIDs(ctx context.Context, playerIDs []uint) (*[]models.Player, error)
	GetPlayersFromTeam(ctx context.Context, teamID uint) (*[]models.Player, error)
	UpdPlayer(ctx context.Context, player *models.Player) (*models.Player, error)
	NewApplication(ctx context.Context, application *models.Application, player *[]models.Player) (
		*models.Application, *[]models.Player, error,
	)
	UpdApplication(ctx context.Context, application *models.Application, players *[]models.Player) (
		*models.Application, *[]models.Player, error,
	)
	GetApplicationFromTeamTournament(ctx context.Context, tournamentID, teamID uint) (*models.Application, error)
	GetApplicationByID(ctx context.Context, applicationID uint) (*models.Application, error)
	GetApplicationsByTeamID(ctx context.Context, teamID uint) (*[]models.Application, error)
	GetPlayersFromApplication(ctx context.Context, applicationID uint) (*[]models.Player, error)
	GetApplicationsFromTournament(ctx context.Context, tournamentID uint) (*[]models.Application, error)
	UpdApplicationTournament(ctx context.Context, application *models.Application) (*models.Application, error)
}

type Config struct {
	Database *database.Config
}

func New(ctx context.Context, cfg Config) (Store, error) {
	if cfg.Database != nil {
		return database.New(ctx, *cfg.Database)
	}

	return nil, errors.New("storage setting is empty")
}
