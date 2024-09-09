package sportspace

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sport-space/internal/adapter/models"
	"sport-space/internal/adapter/storage/errstore"
	"sport-space/pkg/tools"

	"go.uber.org/zap"
)

var (
	lengthShortPassword uint = 8
)

type storage interface {
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
	UpdTeam(ctx context.Context, team *models.Team) (*models.Team, error)
	NewPlayer(ctx context.Context, player *models.Player) (*models.Player, error)
	GetPlayers(ctx context.Context, userID uint) (*[]models.Player, error)
	UpdPlayer(ctx context.Context, player *models.Player) (*models.Player, error)
	AddPlayersTeam(ctx context.Context, playerIDs *[]uint, teamID uint, userID uint) error
	GetPlayersTeam(ctx context.Context, team *models.Team) (*[]models.Player, error)
	RemovePlayersTeam(ctx context.Context, playerIDs *[]uint, teamID uint, userID uint) error
}

type sender interface {
	SendCodeToEmail(email string, code string) (bool, error)
}

type SportSpace struct {
	log    *zap.Logger
	store  storage
	sender sender
}

type option func(s *SportSpace)

func SetLogger(l *zap.Logger) option {
	return func(s *SportSpace) {
		s.log = l
	}
}

func New(store storage, sender sender, options ...option) (*SportSpace, error) {
	s := &SportSpace{
		log:    zap.NewNop(),
		store:  store,
		sender: sender,
	}

	for _, opt := range options {
		opt(s)
	}

	return s, nil
}

func (s *SportSpace) LoginWithOTP(ctx context.Context, email, otp string) (*models.User, error) {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed getting user: %w", err)
	}

	otpStored, err := s.store.GetOTP(ctx, user)
	if err != nil {
		return user, fmt.Errorf("failed getting otp by user: %w", err)
	}

	if otp != otpStored.Password {
		return user, errors.New("otp is not equal")
	}

	err = s.store.RemoveOTP(ctx, user)
	if err != nil {
		// позволяем пользователю войти, пишем в лог ошибку
		s.log.Error("failed remove otp", zap.Uint("userID", user.ID), zap.String("email", email), zap.Error(err))
	}

	return user, nil
}

func (s *SportSpace) NewOTP(ctx context.Context, email string) error {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			// задаем хеш пароль пустым что бы нельзя было авторизоваться по постоянному паролю
			user, err = s.store.NewUser(ctx, email, email, "")
			if err != nil {
				return fmt.Errorf("failed create user: %w", err)
			}
		} else {
			s.log.Error("getting user", zap.String("email", email))
			return fmt.Errorf("failed getting user by email: %w", err)
		}
	}

	otpStore, err := s.store.GetOTP(ctx, user)
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			otpStore = &models.OTPUser{}
			otpStore.UserID = user.ID
		} else {
			return fmt.Errorf("failed get otp: %w", err)
		}
	}

	otpStore.Password = tools.RandomString(lengthShortPassword)
	otpStore.Attempt += 1
	otpStore.UpdatedAt = time.Now()

	err = s.store.NewOTP(ctx, otpStore)
	if err != nil {
		s.log.Error("save otp", zap.Error(err), zap.String("email", email), zap.Uint("userID", user.ID))
		return fmt.Errorf("failed save otp: %w", err)
	}

	_, err = s.sender.SendCodeToEmail(email, otpStore.Password)
	if err != nil {
		return fmt.Errorf("failed send otp to email `%s`: %w", email, err)
	}

	return nil
}

func (s *SportSpace) GetAllTournaments(ctx context.Context) (*[]models.Tournament, error) {
	return s.store.GetAllTournaments(ctx)
}

func (s *SportSpace) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	return s.store.GetUserByID(ctx, userID)
}

func (s *SportSpace) NewTournament(ctx context.Context, tournament *models.Tournament) (
	*models.Tournament, error,
) {
	tournament, err := s.store.NewTournament(ctx, tournament)
	if err != nil {
		return nil, fmt.Errorf("failed create tournament: %w", err)
	}

	return tournament, nil
}

func (s *SportSpace) GetTournaments(ctx context.Context, user *models.User) (*[]models.Tournament, error) {
	return s.store.GetTournaments(ctx, user.ID)
}

func (s *SportSpace) GetTournamentByID(ctx context.Context, tournamentID uint) (*models.Tournament, error) {
	tournament, err := s.store.GetTournamentByID(ctx, tournamentID)
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			return &models.Tournament{}, nil
		}
		return nil, fmt.Errorf("failed get tournamet: %w", err)
	}
	return tournament, nil
}

func (s *SportSpace) UpdTournament(ctx context.Context, tournament *models.Tournament) (*models.Tournament, error) {
	tournament, err := s.store.UpdTournamentByUser(ctx, tournament)
	if err != nil {
		return nil, fmt.Errorf("failed update tournament: %w", err)
	}

	return tournament, nil
}

func (s *SportSpace) NewTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	team, err := s.store.NewTeam(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("failed create team: %w", err)
	}

	return team, nil
}

func (s *SportSpace) GetTeams(ctx context.Context, user *models.User) (*[]models.Team, error) {
	teams, err := s.store.GetTeams(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed get teams: %w", err)
	}

	return teams, nil
}

func (s *SportSpace) GetTeamByID(ctx context.Context, teamID uint) (*models.Team, error) {
	team, err := s.store.GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed get team: %w", err)
	}
	return team, nil
}

func (s *SportSpace) UpdTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	team, err := s.store.UpdTeam(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("failed update team: %w", err)
	}
	return team, nil
}

func (s *SportSpace) NewPlayer(ctx context.Context, player *models.Player) (*models.Player, error) {
	return s.store.NewPlayer(ctx, player)
}

func (s *SportSpace) GetPlayers(ctx context.Context, userID uint) (*[]models.Player, error) {
	return s.store.GetPlayers(ctx, userID)
}

func (s *SportSpace) UpdPlayer(ctx context.Context, player *models.Player) (*models.Player, error) {
	return s.store.UpdPlayer(ctx, player)
}

func (s *SportSpace) AddPlayersTeam(ctx context.Context, playerIDs *[]uint, teamID uint, userID uint) error {
	return s.store.AddPlayersTeam(ctx, playerIDs, teamID, userID)
}

func (s *SportSpace) GetPlayersTeam(ctx context.Context, team *models.Team) (*[]models.Player, error) {
	return s.store.GetPlayersTeam(ctx, team)
}

func (s *SportSpace) RemovePlayersTeam(ctx context.Context, playerIDs *[]uint, teamID uint, userID uint) error {
	err := s.store.RemovePlayersTeam(ctx, playerIDs, teamID, userID)
	if err != nil {
		return fmt.Errorf("failed remove players from team: %w", err)
	}
	return nil
}
