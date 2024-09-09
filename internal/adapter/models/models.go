package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint `gorm:"primarykey"`
	Login        string
	Email        string `gorm:"index"`
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type OTPUser struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint `gorm:"index;not null"`
	User      User
	Password  string
	Attempt   uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Tournament struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint `gorm:"index;not null"`
	User      User
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Team struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint `gorm:"index;not null"`
	User      User
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Player struct {
	ID         uint `gorm:"primarykey"`
	UserID     uint `gorm:"index;not null"`
	User       User
	FirstName  string
	SecondName string
	LastName   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type TeamPlayer struct {
	ID        uint `gorm:"primarykey"`
	PlayerID  uint `gorm:"index:idx_team_player,unique;not null"`
	Player    Player
	TeamID    uint `gorm:"index:idx_team_player,unique;not null"`
	Team      Team
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TournamentApplication struct {
	ID           uint `gorm:"primarykey"`
	TeamID       uint `gorm:"index:idx_application_team;not null"`
	Team         Team
	TournamentID uint `gorm:"index:idx_application_tournament;not null"`
	Tournament   Tournament
	Accepted     bool
	AcceptedDate time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// type ApplicationPlayer struct {
// 	ID            uint `gorm:"primarykey"`
// 	ApplicationID uint `gorm:"index:idx_application_player;not null"`
// 	Application   Application
// 	PlayerID      uint
// 	Player        Player
// 	CreatedAt     time.Time
// 	UpdatedAt     time.Time
// 	DeletedAt     gorm.DeletedAt `gorm:"index"`
// }
