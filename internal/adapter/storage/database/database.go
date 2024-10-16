package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"sport-space/internal/adapter/models"
	"sport-space/internal/adapter/storage/errstore"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Storage struct {
	db  *gorm.DB
	log *zap.Logger
}

type Config struct {
	DSN string `env:"DATABASE_URI"`
}

type option func(s *Storage)

func New(ctx context.Context, cfg Config, options ...option) (*Storage, error) {
	var err error
	s := &Storage{
		log: zap.NewNop(),
	}
	lgr := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: lgr,
	})
	if err != nil {
		return nil, fmt.Errorf("failed connect to database: %w", err)
	}

	s.db = db.WithContext(ctx)

	for _, opt := range options {
		opt(s)
	}

	err = s.db.AutoMigrate(
		&models.User{},
		&models.OTPUser{},
		&models.Tournament{},
		&models.Team{},
		&models.Player{},
		// &models.TeamPlayer{},
		&models.Application{},
		// &models.ApplicationPlayer{},
	)

	if err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return s, nil
}

func (s *Storage) NewUser(ctx context.Context, login, email, passwordHash string) (*models.User, error) {
	user := &models.User{
		Login:        login,
		Email:        email,
		PasswordHash: passwordHash,
	}

	err := s.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed create user: %w", err)
	}

	return user, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{
		Email: email,
	}
	err := s.db.WithContext(ctx).Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Join(err, errstore.ErrNotFoundData)
		}
		return nil, err
	}
	return user, nil
}

func (s *Storage) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	user := &models.User{
		ID: userID,
	}
	err := s.db.WithContext(ctx).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Join(err, errstore.ErrNotFoundData)
		}
		return nil, err
	}
	return user, nil
}

func (s *Storage) NewOTP(ctx context.Context, otp *models.OTPUser) error {
	err := s.db.WithContext(ctx).Save(otp).Error
	if err != nil {
		return fmt.Errorf("failed create otp: %w", err)
	}
	return nil
}

func (s *Storage) GetOTP(ctx context.Context, user *models.User) (*models.OTPUser, error) {
	otp := &models.OTPUser{
		UserID: user.ID,
	}
	err := s.db.WithContext(ctx).Where("user_id = ?", user.ID).First(otp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return otp, errors.Join(err, errstore.ErrNotFoundData)
		}
		return otp, fmt.Errorf("failed get otp: %w", err)
	}
	return otp, nil
}

func (s *Storage) RemoveOTP(ctx context.Context, user *models.User) error {
	otp := &models.OTPUser{}
	err := s.db.WithContext(ctx).Where("user_id = ?", user.ID).Delete(otp).Error
	if err != nil {
		return fmt.Errorf("failed remove opt: %w", err)
	}

	return nil
}

func (s *Storage) GetAllTournaments(ctx context.Context) (tournaments *[]models.Tournament, err error) {
	tournaments = &[]models.Tournament{}
	err = s.db.Find(tournaments).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tournaments, errors.Join(err, errstore.ErrNotFoundData)
		}
		return nil, fmt.Errorf("failed get all torurnaments: %w", err)
	}
	return tournaments, nil
}

func (s *Storage) NewTournament(ctx context.Context, tournament *models.Tournament) (
	*models.Tournament, error,
) {
	err := s.db.Create(tournament).Error
	if err != nil {
		return nil, fmt.Errorf("failed create torunament: %w", err)
	}
	return tournament, nil
}

func (s *Storage) GetTournaments(ctx context.Context, userID uint) (*[]models.Tournament, error) {
	tournaments := &[]models.Tournament{}
	err := s.db.Where("user_id = ?", userID).Find(tournaments).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Join(err, errstore.ErrNotFoundData)
		}
		return nil, fmt.Errorf("failed found tournaments by user: %w", err)
	}

	return tournaments, nil
}

func (s *Storage) GetTournamentByID(ctx context.Context, tournamentID uint) (*models.Tournament, error) {
	tournament := &models.Tournament{}
	err := s.db.Where("id = ?", tournamentID).Preload("Applications").First(tournament).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Join(err, errstore.ErrNotFoundData)
		}
		return nil, fmt.Errorf("failed found tournament by id: %w", err)
	}

	return tournament, nil
}

func (s *Storage) UpdTournamentByUser(ctx context.Context, tournament *models.Tournament) (*models.Tournament, error) {
	res := s.db.Model(tournament).
		Where("user_id = ? and id = ?", tournament.UserID, tournament.ID).
		Save(tournament)
	if res.RowsAffected == 0 {
		return nil, errstore.ErrNotFoundData
	}
	if err := res.Error; err != nil {
		return nil, fmt.Errorf("failed update tournament by user: %w", err)
	}

	return tournament, nil
}

func (s *Storage) NewTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	err := s.db.Create(team).Error
	if err != nil {
		return nil, fmt.Errorf("failed create team: %w", err)
	}
	return team, err
}

func (s *Storage) GetTeams(ctx context.Context, user *models.User) (*[]models.Team, error) {
	teams := &[]models.Team{}
	err := s.db.Where("user_id = ?", user.ID).Find(teams).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return teams, errstore.ErrNotFoundData
		}
		return nil, fmt.Errorf("failed find teams by user: %w", err)
	}

	return teams, nil
}

func (s *Storage) GetTeamByID(ctx context.Context, teamID uint) (*models.Team, error) {
	team := &models.Team{}
	err := s.db.Where("id = ?", teamID).Preload("Players").First(team).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errstore.ErrNotFoundData
		}
		return nil, fmt.Errorf("failed get team by id: %w", err)
	}
	return team, nil
}

func (s *Storage) UpdTeam(ctx context.Context, team *models.Team, playersIDs *[]uint) (*models.Team, *[]models.Player, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if playersIDs != nil {
			_players := &[]models.Player{}
			err := tx.Where("id IN ?", *playersIDs).Find(_players).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed find players: %w", err)
			}
			err = tx.Model(team).Association("Players").Replace(_players)
			if err != nil {
				return fmt.Errorf("failed replace players: %w", err)
			}
			team.Players = *_players
		}
		err := tx.Save(team).Error
		if err != nil {
			return fmt.Errorf("failed update team: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed update team by user: %w", err)
	}
	return team, &team.Players, nil
}

func (s *Storage) NewPlayer(ctx context.Context, player *models.Player) (*models.Player, error) {
	err := s.db.Create(player).Error
	if err != nil {
		return nil, fmt.Errorf("failed create player: %w", err)
	}
	return player, nil
}
func (s *Storage) NewPlayerBatch(ctx context.Context, players *[]models.Player) (*[]models.Player, error) {
	err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"first_name", "second_name", "last_name", "b_day", "photo_url"}),
	}).Create(players).Error
	if err != nil {
		return nil, fmt.Errorf("failed create batch players: %w", err)
	}
	return players, nil
}

func (s *Storage) GetPlayers(ctx context.Context, userID uint) (*[]models.Player, error) {
	players := &[]models.Player{}
	err := s.db.Where("user_id = ?", userID).Find(players).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return players, errstore.ErrNotFoundData
		}
		return nil, fmt.Errorf("failed find players: %w", err)
	}

	return players, nil
}

func (s *Storage) GetPlayerByID(ctx context.Context, playerID uint) (*models.Player, error) {
	player := &models.Player{}
	err := s.db.Where("id = ?", playerID).First(player).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Join(err, errstore.ErrNotFoundData)
		}
		return nil, fmt.Errorf("failed find player: %w", err)
	}
	return player, nil
}

func (s *Storage) GetPlayersByIDs(ctx context.Context, playerIDs []uint) (*[]models.Player, error) {
	players := &[]models.Player{}
	err := s.db.Where("id IN ?", playerIDs).Find(players).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed find player: %w", err)
	}
	return players, nil
}

func (s *Storage) UpdPlayer(ctx context.Context, player *models.Player) (*models.Player, error) {
	res := s.db.Where("id = ? and user_id = ?", player.ID, player.UserID).Save(player)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Join(err, errstore.ErrNotFoundData)
		}
		return nil, fmt.Errorf("failed update player: %w", err)
	}
	if res.RowsAffected == 0 {
		return nil, errstore.ErrNotFoundData
	}
	return player, nil
}

func (s *Storage) GetPlayersFromTeam(ctx context.Context, teamID uint) (*[]models.Player, error) {
	players := &[]models.Player{}
	err := s.db.Joins("JOIN team_players on team_players.player_id = players.id", s.db.Where("team_players.team_id = ?", teamID)).Find(players).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errstore.ErrNotFoundData
		}
		return nil, fmt.Errorf("failed found players from team: %w", err)
	}

	return players, nil
}

func (s *Storage) NewApplication(ctx context.Context, application *models.Application, players *[]models.Player) (
	*models.Application, *[]models.Player, error,
) {
	application.Status = models.Draft
	application.StatusDate = time.Now()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(application).Error
		if err != nil {
			return fmt.Errorf("failed create application: %w", err)
		}
		err = tx.Model(application).Association("Players").Replace(players)
		if err != nil {
			return fmt.Errorf("failed create application batch players: %w", err)
		}

		return nil
	})
	if err != nil {
		var sqlError *pgconn.PgError
		if errors.As(err, &sqlError) && sqlError.Code == pgerrcode.UniqueViolation {
			return nil, nil, errors.Join(err, errstore.ErrConflictData)
		}
		return nil, nil, fmt.Errorf("failed create application with transactions: %w", err)
	}

	return application, players, nil
}

func (s *Storage) GetApplicationFromTeamTournament(ctx context.Context, tournamentID, teamID uint) (*models.Application, error) {
	application := &models.Application{}
	err := s.db.WithContext(ctx).Where("tournament_id = ? and team_id = ?", tournamentID, teamID).First(application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errstore.ErrNotFoundData
		}
		return nil, fmt.Errorf("failed get appliction: %w", err)
	}
	return application, nil
}

func (s *Storage) GetApplicationByID(ctx context.Context, applicationID uint) (*models.Application, error) {
	application := &models.Application{}
	err := s.db.Where("id = ?", applicationID).Preload("Players").First(application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("not found application: %w", errors.Join(err, errstore.ErrNotFoundData))
		}
		return nil, fmt.Errorf("failed find application: %w", err)
	}
	return application, nil
}

func (s *Storage) UpdApplication(ctx context.Context, application *models.Application, players *[]models.Player) (
	*models.Application, *[]models.Player, error,
) {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("id = ?", application.ID).Updates(application).Error
		if err != nil {
			return fmt.Errorf("failed update application: %w", err)
		}
		if *players != nil {
			err = tx.Model(application).Association("Players").Replace(players)
			if err != nil {
				return fmt.Errorf("failed create application batch players: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("not found application: %w", errors.Join(err, errstore.ErrNotFoundData))
		}
		var sqlError *pgconn.PgError
		if errors.As(err, &sqlError) && sqlError.Code == pgerrcode.UniqueViolation {
			return nil, nil, errors.Join(err, errstore.ErrConflictData)
		}
		return nil, nil, fmt.Errorf("failed update application with transactions: %w", err)
	}

	return application, players, nil
}

func (s *Storage) GetApplicationsByTeamID(ctx context.Context, teamID uint) (*[]models.Application, error) {
	applications := &[]models.Application{}
	err := s.db.Where("team_id = ?", teamID).Find(applications).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return applications, fmt.Errorf("not found applications: %w", errors.Join(err, errstore.ErrNotFoundData))
		}
		return nil, fmt.Errorf("failed get applications: %w", err)
	}

	return applications, nil
}

func (s *Storage) GetPlayersFromApplication(ctx context.Context, applicationID uint) (*[]models.Player, error) {
	application := &models.Application{}
	err := s.db.Where("id = ?", applicationID).Preload("Players").First(application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &application.Players, fmt.Errorf("not found players from application: %w", errors.Join(err, errstore.ErrNotFoundData))
		}
		return nil, fmt.Errorf("failed get players from application: %w", err)
	}
	return &application.Players, nil
}

func (s *Storage) GetApplicationsFromTournament(ctx context.Context, tournamentID uint) (*[]models.Application, error) {
	applications := &[]models.Application{}
	statuses := []models.ApplicationStatus{models.InProgress, models.Accepted, models.Rejected}
	err := s.db.Model(&models.Tournament{}).Where("tournament_id = ? and status in ?", tournamentID, statuses).Association("Applications").Find(applications)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("not found applications: %w", errors.Join(err, errstore.ErrNotFoundData))
		}
		return nil, fmt.Errorf("failed get applications: %w", err)
	}
	return applications, nil
}

func (s *Storage) UpdApplicationTournament(ctx context.Context, application *models.Application) (*models.Application, error) {
	err := s.db.WithContext(ctx).Where("id = ?", application.ID).Updates(application).Error
	if err != nil {
		return nil, fmt.Errorf("failed update application: %w", err)
	}
	return application, nil
}
