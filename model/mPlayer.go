package model

import (
	"database/sql"
	"strings"
	"time"

	"gorm.io/gorm"
)

/*
 * Игроки
 */
type Player struct {
	ID                           uint `gorm:"primarykey"`
	FirstName                    string
	SecondName                   string
	LastName                     string
	BDay                         sql.NullTime
	Teams                        []Team `gorm:"many2many:team_players"`
	TournamentApplicationPlayers []TournamentApplicationPlayer
	UserID                       uint
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	// DeletedAt gorm.DeletedAt `gorm:"index"`
}

func GetPlayerByUserId(userId uint) (Player, error) {
	player := Player{}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return player, err
	}

	result := db.Where("user_id = ?", userId).First(&player)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return player, result.Error
	}

	return player, nil
}

func GetPlayerInviteToTeamByEmail(email string, status TeamInviteStatus) ([]TeamInvite, error) {
	var invites []TeamInvite

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return invites, err
	}

	result := db.Where("email = ?", strings.ToLower(email))
	if status != "" {
		result = result.Where("status = ?", status)
	}
	result = result.Find(&invites)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return invites, result.Error
	}

	return invites, nil
}

func GetPlayerInviteToTeamById(id uint, status TeamInviteStatus) (TeamInvite, error) {
	var invites TeamInvite

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return invites, err
	}

	result := db.Where("id = ?", id)
	result = result.Where("status = ?", status.ToString())
	result = result.First(&invites)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return invites, result.Error
	}

	return invites, nil
}

func UpdatePlayerInviteToTeam(invite TeamInvite) (TeamInvite, error) {

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return invite, err
	}

	result := db.Save(invite)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return invite, result.Error
	}

	return invite, nil
}

func AcceptedPlayerInviteToTeam(invite TeamInvite, userId uint) (TeamInvite, error) {
	invite.Status = TIAccepted

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return invite, err
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// user, err := FindUserByEmail(strings.ToLower(invite.Email))
	// if err != nil {
	// 	return invite, err
	// }
	// if user.ID == 0 {
	// 	return invite, fmt.Errorf("User not found")
	// }

	result := tx.Save(invite)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return invite, result.Error
	}

	team, err := GetTeamById(invite.TeamID)
	if err != nil {
		return invite, err
	}

	player, err := GetPlayerByUserId(userId)
	if err != nil {
		return invite, err
	}
	if player.ID == 0 {
		player.UserID = userId
	}
	// } else {
	// }

	result = tx.Save(&player)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return invite, result.Error
	}

	player.Teams = append(player.Teams, team)
	result = tx.Save(&player)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return invite, result.Error
	}

	return invite, tx.Commit().Error
}
