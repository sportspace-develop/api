package model

import (
	"time"

	"gorm.io/gorm"
)

type Organization struct {
	ID         uint `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Title      string         `json:"title"`
	Address    string         `json:"addess"`
	Tournament []Tournament
	UserID     uint `json:"-"`
}

func CreateOrganization(title string, userId uint, address string) (Organization, error) {
	organization := Organization{
		Title:   title,
		UserID:  userId,
		Address: address,
	}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return organization, err
	}

	result := db.Create(&organization)
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return organization, result.Error
	}

	return organization, nil
}

func GetOrganizations() ([]Organization, error) {
	organization := []Organization{}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return organization, err
	}

	result := db.Find(&organization)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return organization, result.Error
	}

	return organization, nil
}

func GetOrganizationById(organizationId uint) (Organization, error) {
	organization := Organization{}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return organization, err
	}

	result := db.Where("id = ?", organizationId).First(&organization)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return organization, result.Error
	}

	return organization, nil
}

func GetOrganizationByUserId(userId uint) (Organization, error) {
	organization := Organization{}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return organization, err
	}

	result := db.Where("user_id = ?", userId).First(&organization)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return organization, result.Error
	}

	return organization, nil
}

func UpdateOrganization(o Organization) (Organization, error) {
	// organization := Organization{
	// 	Title:   title,
	// 	UserID:  userId,
	// 	Address: address,
	// }

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return o, err
	}

	result := db.Save(&o)
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return o, result.Error
	}

	return o, nil
}
