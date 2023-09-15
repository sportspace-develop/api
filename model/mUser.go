package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string
	Username     sql.NullString
	Password     sql.NullString
	Session      []UserSession
	Code         []UserAuthCode
	Tournament   []Tournament
	Organization []Organization
	Player       []Player
}

type UserSession struct {
	ID           uint `gorm:"primarykey"`
	UserId       uint
	RefreshToken string
	ExpiresIn    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type userAuthCodeStatus string

const (
	CSWait      userAuthCodeStatus = "wait"
	CSActivated userAuthCodeStatus = "activated"
)

type UserAuthCode struct {
	gorm.Model
	UserId       uint
	Code         string
	Status       userAuthCodeStatus `gorm:"type:enum('wait', 'activated')";"column:userAuthCodeStatus"`
	ExpiresIn    time.Time
	AttemptCount uint8
	AttemptDate  time.Time
	BlockDate    sql.NullTime
}

func Registration(email string) (User, error) {
	user := User{
		Email: email,
	}
	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return user, err
	}

	result := db.Create(&user)
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return user, result.Error
	}

	return user, nil
}

func FindUserByEmail(email string) (User, error) {
	user := User{}
	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return user, err
	}

	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return user, result.Error
	}

	return user, nil
}

func FindUserById(id uint) (User, error) {
	user := User{}
	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return user, err
	}

	result := db.First(&user, id)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return user, result.Error
	}

	return user, nil
}

func CreateUserAuthCode(user User, code string) (bool, error) {
	uCode := UserAuthCode{
		UserId:       user.ID,
		Code:         code,
		ExpiresIn:    time.Now().UTC().Add(time.Duration(time.Minute * 60)),
		AttemptCount: 1,
		AttemptDate:  time.Now().UTC(),
		Status:       CSWait,
	}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return false, err
	}

	result := db.Create(&uCode)
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return false, result.Error
	}

	return true, nil
}

func ActivateUserAuthCode(authCode UserAuthCode) (UserAuthCode, error) {

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return authCode, err
	}

	authCode.Status = CSActivated
	result := db.Save(&authCode)
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return authCode, result.Error
	}

	return authCode, nil
}

func FindCodeNotActivatedByUserCode(user User, code string) (UserAuthCode, error) {
	var dataAuthCode UserAuthCode

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return dataAuthCode, err
	}

	result := db.Where("expires_in > ? and status = ? and code = ? and user_id = ?", time.Now().UTC(), string(CSWait), code, user.ID).First(&dataAuthCode)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return dataAuthCode, result.Error
	}

	return dataAuthCode, nil
}

func FindCodeNotActivatedByUser(user User) (UserAuthCode, error) {
	var dataAuthCode UserAuthCode

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return dataAuthCode, err
	}

	result := db.Where("expires_in > ? and status = ? and user_id = ?", time.Now().UTC(), string(CSWait), user.ID).First(&dataAuthCode)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return dataAuthCode, result.Error
	}

	return dataAuthCode, nil
}

func UpdateCodeNotActivatedByUser(data UserAuthCode) (UserAuthCode, error) {

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return data, err
	}

	result := db.Save(&data)
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return data, result.Error
	}

	return data, nil
}

func NewSession(user User, refreshToken string) (UserSession, error) {
	session := UserSession{
		UserId:       user.ID,
		RefreshToken: refreshToken,
		ExpiresIn:    time.Now().UTC().Add(time.Hour * 24 * 30),
	}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return session, err
	}

	result := db.Create(&session)
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return session, result.Error
	}

	return session, nil
}

func DeleteSession(refreshToken string) (bool, error) {
	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return false, err
	}

	result := db.Where("refresh_token = ?", refreshToken).Delete(&UserSession{})
	if result.Error != nil {
		log.ERROR(result.Error.Error())
		return false, result.Error
	}

	return true, nil
}

func FindSession(refreshToken string) (UserSession, error) {
	session := UserSession{}

	db, err := Connect()
	if err != nil {
		log.ERROR(err.Error())
		return session, err
	}

	result := db.Where("refresh_token = ?", refreshToken).First(&session)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.ERROR(result.Error.Error())
		return session, result.Error
	}

	return session, nil
}
