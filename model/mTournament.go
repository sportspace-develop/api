package model

import (
	"time"

	"gorm.io/gorm"
)

/*
 * Справочник игр
 */
type DGame struct {
	gorm.Model
	Title      string `gorm:"unique" `
	Name       string `gorm:"unique"`
	Team       []Team
	Tournament []Tournament
}

type Tournament struct {
	ID                    uint      `gorm:"primarykey" json:"id"`
	Title                 string    `json:"title"`
	StartDate             time.Time `json:"start_date" gorm:"index"`
	EndDate               time.Time `json:"end_date" gorm:"index"`
	StartRegistrationDate time.Time `json:"start_registration_date"`
	EndRegistrationDate   time.Time `json:"end_registration_date"`
	DGameID               uint      `json:"game_type_id"` // вид спорта
	UserID                uint      `json:"-" gorm:"index"`
	IsTeam                bool      `gorm:"default:true" json:"is_team"` // команды или индивидульно

	OrganizationID        uint `json:"organization_id" gorm:"index"`
	TournamentApplication []TournamentApplication

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// заявки на турнир
type TournamentApplication struct {
	gorm.Model
	TournamentID                 uint `gorm:"index"`
	TeamID                       uint `gorm:"index"`
	Status                       string
	TournamentApplicationPlayers []TournamentApplicationPlayer
}

type TournamentApplicationPlayer struct {
	gorm.Model
	TournamentApplicationID uint
	PlayerID                uint
}

func GetGameTypes() ([]DGame, error) {
	var gTypes []DGame

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return gTypes, err
	}

	result := db.Find(&gTypes)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return gTypes, result.Error
	}
	return gTypes, nil
}

func CreateTournament(title string, userId uint, startDate, endDate, startRegistrationDate, endRegistrationDate time.Time, isTeam bool, gameType uint) (Tournament, error) {
	tournament := Tournament{
		Title:                 title,
		UserID:                userId,
		StartDate:             startDate,
		EndDate:               endDate,
		StartRegistrationDate: startRegistrationDate,
		EndRegistrationDate:   endRegistrationDate,
		DGameID:               gameType,
	}
	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return tournament, err
	}

	result := db.Create(&tournament)
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return tournament, result.Error
	}

	return tournament, nil
}

func GetTournamentById(tournamentId uint) (Tournament, error) {
	tournament := Tournament{}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return tournament, err
	}

	result := db.Where("id = ?", tournamentId).First(&tournament)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return tournament, result.Error
	}

	return tournament, nil
}

func GetTournaments[ID uint | int](id ID) ([]Tournament, error) {
	tournaments := []Tournament{}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return tournaments, err
	}

	result := db
	if id > 0 {
		result = result.Where("id = ?", id)
	}
	result = result.Find(&tournaments)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return tournaments, result.Error
	}

	return tournaments, nil
}

func GetTournamentsByUser(userId uint) ([]Tournament, error) {
	tournaments := []Tournament{}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return tournaments, err
	}

	result := db.Where("user_id = ?", userId).Find(&tournaments)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return tournaments, result.Error
	}

	return tournaments, nil
}

func UpdateTournament(tournament Tournament) (Tournament, error) {
	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return tournament, err
	}

	result := db.Save(&tournament)
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return tournament, result.Error
	}

	return tournament, nil
}
