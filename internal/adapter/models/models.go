package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint `gorm:"primarykey"`
	Login        string
	Email        string `gorm:"index"`
	Tournaments  []Tournament
	Teams        []Team
	Players      []Player
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
	ID                uint `gorm:"primarykey"`
	UserID            uint `gorm:"index;not null"`
	Title             string
	LogoURL           string
	Applications      []Application
	StartDate         *time.Time `gorm:"default:null"`
	EndDate           *time.Time `gorm:"default:null"`
	RegisterStartDate *time.Time `gorm:"default:null"`
	RegisterEndDate   *time.Time `gorm:"default:null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

type Team struct {
	ID           uint `gorm:"primarykey"`
	UserID       uint `gorm:"index;not null"`
	Title        string
	Players      []Player `gorm:"many2many:team_players"`
	Applications []Application
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Player struct {
	ID         uint `gorm:"primarykey"`
	UserID     uint `gorm:"index;not null"`
	FirstName  string
	SecondName string
	LastName   string
	BDay       *time.Time `gorm:"default:null"`
	PhotoURL   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type ApplicationStatus string

const (
	Draft      ApplicationStatus = "draft"
	InProgress ApplicationStatus = "inprogress"
	Accepted   ApplicationStatus = "accepted"
	Rejected   ApplicationStatus = "rejected"
	Canceled   ApplicationStatus = "canceled"
)

type Application struct {
	ID           uint              `gorm:"primarykey"`
	TeamID       uint              `gorm:"index:idx_application,unique;not null"`
	TournamentID uint              `gorm:"index:idx_application,unique;not null"`
	Players      []Player          `gorm:"many2many:application_players"`
	Status       ApplicationStatus `gorm:"index:idx_status;not null"`
	StatusDate   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
