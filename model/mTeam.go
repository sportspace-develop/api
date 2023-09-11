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
	TeamInvite            []TeamInvite
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt `gorm:"index"`
}

/*
 * Инвайты в команду
 */
type teamInviteStatus string

const (
	TIWait     teamInviteStatus = "wait"
	TISended   teamInviteStatus = "sended"
	TISuccess  teamInviteStatus = "success"
	TIRejected teamInviteStatus = "rejected"
	TICancel   teamInviteStatus = "cancel"
)

func (s teamInviteStatus) ToString() string {
	return string(s)
}

func (s teamInviteStatus) Title() string {
	switch s {
	case TISended:
		return "отправлен"
	case TISuccess:
		return "подтвержден"
	case TICancel:
		return "отменен"
	case TIRejected:
		return "откланен"
	default:
		return "ожидает"
	}
}

type TeamInvite struct {
	ID        uint   `gorm:"primarykey"`
	Email     string `gorm:"index:,unique,composite:invite"`
	Code      string
	TeamID    uint             `gorm:"index:,unique,composite:invite"`
	Status    teamInviteStatus `gorm:"type:enum('wait', 'sended', 'success', 'rejected', 'cancel')";"column:teamInviteStatus"`
	CreatedAt time.Time
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

//TeamInvite

func GetInvitesToTeamByTeam(teamId uint) ([]TeamInvite, error) {
	var invites []TeamInvite
	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return invites, err
	}

	result := db.Where("team_id = ?", teamId).Find(&invites)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return invites, result.Error
	}

	return invites, nil

}

func CreateInvitesToTeam(invites []TeamInvite) ([]TeamInvite, error) {
	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return invites, err
	}

	result := db.Create(&invites)
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return invites, result.Error
	}

	return invites, nil
}

func UpdateInvitesToTeam(invites []TeamInvite) ([]TeamInvite, error) {
	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return invites, err
	}

	result := db.Save(&invites)
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return invites, result.Error
	}

	return invites, nil
}

func CreateOrUpdateInvitesToTeam(invites []TeamInvite, teamId uint) ([]TeamInvite, error) {

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return invites, err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	createdInvites, err := GetInvitesToTeamByTeam(teamId)
	if err != nil {
		return nil, err
	}

	var reInvites []TeamInvite
	var newInvites []TeamInvite

	for _, invite := range invites {
		var new = true
		for _, createdInvite := range createdInvites {
			if invite.Email == createdInvite.Email {
				reInvites = append(reInvites, createdInvite)
				new = false
				break
			}
		}
		if new {
			newInvites = append(newInvites, invite)
		}
	}

	if len(createdInvites) == 0 {
		newInvites = invites
	}

	if len(newInvites) > 0 {
		result := tx.Create(&newInvites)
		if result.Error != nil {
			tx.Rollback()
			log.ERROR(result.Error.Error())
			return invites, result.Error
		}
	}

	if len(reInvites) > 0 {
		result := tx.Save(&reInvites)
		if result.Error != nil {
			tx.Rollback()
			log.ERROR(result.Error.Error())
			return invites, result.Error
		}
	}

	return append(newInvites, reInvites...), tx.Commit().Error
}
