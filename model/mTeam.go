package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

/*
 * Игроки
 */
type Member struct {
	gorm.Model
	FirstName                       string
	SecondName                      string
	LastName                        string
	BDay                            sql.NullTime
	TeamID                          uint
	MemberOfTournamentApplicationID uint
}

/*
 * Команды
 */
type Team struct {
	ID                    uint   `gorm:"primarykey"`
	Title                 string `gorm:"unique"`
	DGameID               uint
	UserID                uint
	TournamentApplication []TournamentApplication
	Member                []Member
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt `gorm:"index"`
}

func GetTeams(id uint) ([]Team, error) {
	team := []Team{}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return team, err
	}

	result := db
	if id > 0 {
		result = result.Where("id = ?", id)
	}
	result = result.Find(&team)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return team, result.Error
	}

	return team, nil
}

func GetTeamById(Id uint) (Team, error) {
	team := Team{}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return team, err
	}

	result := db.Where("id = ?", Id).First(&team)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return team, result.Error
	}

	return team, nil
}

func GetTeamByUser(userId uint) (Team, error) {
	team := Team{}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return team, err
	}

	result := db.Where("user_id = ?", userId).First(&team)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return team, result.Error
	}

	return team, nil
}

func CreateTeam(title string, gameTypeId uint, userId uint) (Team, error) {
	team := Team{
		Title:   title,
		DGameID: gameTypeId,
		UserID:  userId,
	}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return team, err
	}

	result := db.Create(&team)
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return team, result.Error
	}

	return team, nil
}

func UpdateTeam(t Team) (Team, error) {
	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return t, err
	}

	result := db.Save(&t)
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return t, result.Error
	}

	return t, nil
}
