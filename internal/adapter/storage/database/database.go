package database

import (
	"context"
	"errors"
	"fmt"

	"sport-space/internal/adapter/models"
	"sport-space/internal/adapter/storage/errstore"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
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
		&models.TeamPlayer{},
		&models.TournamentApplication{},
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
	err := s.db.WithContext(ctx).First(user).Error
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
	err := s.db.WithContext(ctx).First(otp).Error
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
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed found tournaments by user: %w", err)
	}

	return tournaments, nil
}

func (s *Storage) GetTournamentByID(ctx context.Context, tournamentID uint) (*models.Tournament, error) {
	tournament := &models.Tournament{}
	err := s.db.Where("id = ?", tournamentID).First(tournament).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed found tournament by id: %w", err)
	}

	return tournament, nil
}

func (s *Storage) UpdTournamentByUser(ctx context.Context, tournament *models.Tournament) (*models.Tournament, error) {
	res := s.db.Model(tournament).
		Where("user_id = ? and id = ?", tournament.UserID, tournament.ID).
		Updates(models.Tournament{Title: tournament.Title})
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
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed find teams by user: %w", err)
	}

	return teams, nil
}

func (s *Storage) GetTeamByID(ctx context.Context, teamID uint) (*models.Team, error) {
	team := &models.Team{}
	err := s.db.Where("id = ?", teamID).First(team).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed get team by id: %w", err)
	}
	return team, nil
}

func (s *Storage) UpdTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	res := s.db.Model(team).
		Where("user_id = ? and id = ?", team.UserID, team.ID).
		Updates(models.Tournament{Title: team.Title})
	if res.RowsAffected == 0 {
		return nil, errstore.ErrNotFoundData
	}
	if err := res.Error; err != nil {
		return nil, fmt.Errorf("failed update team by user: %w", err)
	}
	return team, nil
}

func (s *Storage) NewPlayer(ctx context.Context, player *models.Player) (*models.Player, error) {
	err := s.db.Create(player).Error
	if err != nil {
		return nil, fmt.Errorf("failed create player: %w", err)
	}
	return player, nil
}

func (s *Storage) GetPlayers(ctx context.Context, userID uint) (*[]models.Player, error) {
	players := &[]models.Player{}
	err := s.db.Where("user_id = ?", userID).Find(players).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed find players: %w", err)
	}

	return players, nil
}

func (s *Storage) UpdPlayer(ctx context.Context, player *models.Player) (*models.Player, error) {
	res := s.db.Where("id = ? and user_id = ?", player.ID, player.UserID).Updates(player)
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

func (s *Storage) AddPlayersTeam(ctx context.Context, playerIDs *[]uint, teamID uint, userID uint) error {
	if len(*playerIDs) == 0 {
		return fmt.Errorf("list IDs is empty: %w", errstore.ErrNotFoundData)
	}

	team, err := s.GetTeamByID(ctx, teamID)
	if err != nil {
		return fmt.Errorf("failed find team by id `%d`: %w", teamID, err)
	}

	if team.UserID != userID {
		return fmt.Errorf("not found team by user: %w", errstore.ErrNotFoundData)
	}

	tx := s.db.Begin().WithContext(ctx)
	defer tx.Rollback()
	players := &[]models.Player{}
	err = tx.Where("user_id = ?", userID).Find(players, playerIDs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Join(err, errstore.ErrNotFoundData)
		}
		return fmt.Errorf("failed get players by IDs: %w", err)
	}

	teamPlayer := []models.TeamPlayer{}
	for _, id := range *playerIDs {
		teamPlayer = append(teamPlayer, models.TeamPlayer{PlayerID: id, TeamID: teamID})
	}
	err = tx.CreateInBatches(&teamPlayer, len(teamPlayer)).Error
	if err != nil {
		var sqlError *pgconn.PgError
		if errors.As(err, &sqlError) && sqlError.Code == pgerrcode.UniqueViolation {
			return errors.Join(err, errstore.ErrConflictData)
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.Join(err, errstore.ErrConflictData)
		}
		return fmt.Errorf("failed add players to team: %w", err)
	}
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("failed transaction: %w", err)
	}
	return nil
}

func (s *Storage) GetPlayersTeam(ctx context.Context, team *models.Team) (*[]models.Player, error) {
	var players *[]models.Player
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("id = ? and user_id = ?", team.ID, team.UserID).First(team).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.Join(err, errstore.ErrNotFoundData)
			}
			return fmt.Errorf("failed get team: %w", err)
		}
		teamPlayers := &[]models.TeamPlayer{}
		err = tx.Where("team_id =?", team.ID).Find(teamPlayers).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("failed get players from team: %w", err)
		}
		IDs := []uint{}
		for _, player := range *teamPlayers {
			IDs = append(IDs, player.PlayerID)
		}

		err = tx.Where("id in ?", IDs).Find(&players).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed find players: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed get players by team: %w", err)
	}

	return players, nil
}

func (s *Storage) RemovePlayersTeam(ctx context.Context, playerIDs *[]uint, teamID uint, userID uint) error {
	tx := s.db.Begin().WithContext(ctx)
	defer tx.Rollback()

	team, err := s.GetTeamByID(ctx, teamID)
	if err != nil {
		return fmt.Errorf("failed get players by team: %w", err)
	}
	if team.UserID != userID {
		return fmt.Errorf("not found team by user: %w", errstore.ErrNotFoundData)
	}

	teamPlayers := &[]models.TeamPlayer{}
	err = tx.Where("team_id = ? and player_id IN ?", team.ID, *playerIDs).Find(teamPlayers).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("failed get players from team: %w", err)
	}
	if len(*teamPlayers) == 0 {
		return fmt.Errorf("found players in team `%d` is empty: %w", teamID, errstore.ErrNotFoundData)
	}

	err = tx.Delete(teamPlayers).Error
	if err != nil {
		return fmt.Errorf("failed remove players from team: %w", err)
	}

	if err = tx.Commit().Error; err != nil {
		return fmt.Errorf("failed commit deleting: %w", err)
	}
	return nil
}
