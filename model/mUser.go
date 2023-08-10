package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string
	Username sql.NullString
	Password sql.NullString
	Session  []Session
	Code     []UserAuthCode
}

type Session struct {
	gorm.Model
	UserId       uint
	RefreshToken string
	ExpiresIn    time.Time
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
		return user, err
	}

	result := db.Create(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func FindUserByEmail(email string) (User, error) {
	user := User{
		Email: email,
	}
	db, err := Connect()
	if err != nil {
		return user, err
	}

	result := db.First(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return user, result.Error
	}

	return user, nil
}

func FindUserById(id uint) (User, error) {
	user := User{}
	db, err := Connect()
	if err != nil {
		return user, err
	}

	result := db.First(&user, id)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
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
		return false, err
	}

	result := db.Create(&uCode)
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func ActivateUserAuthCode(authCode UserAuthCode) (UserAuthCode, error) {

	db, err := Connect()
	if err != nil {
		return authCode, err
	}

	authCode.Status = CSActivated
	result := db.Save(&authCode)
	if result.Error != nil {
		return authCode, result.Error
	}

	return authCode, nil
}

func FindCodeNotActivatedByUserCode(user User, code string) (UserAuthCode, error) {
	var dataAuthCode UserAuthCode

	db, err := Connect()
	if err != nil {
		return dataAuthCode, err
	}

	result := db.Where("expires_in > ? and status = ? and code = ? and user_id = ?", time.Now().UTC(), string(CSWait), code, user.ID).First(&dataAuthCode)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return dataAuthCode, result.Error
	}

	return dataAuthCode, nil
}

func FindCodeNotActivatedByUser(user User) (UserAuthCode, error) {
	var dataAuthCode UserAuthCode

	db, err := Connect()
	if err != nil {
		return dataAuthCode, err
	}

	result := db.Where("expires_in > ? and status = ? and user_id = ?", time.Now().UTC(), string(CSWait), user.ID).First(&dataAuthCode)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return dataAuthCode, result.Error
	}

	return dataAuthCode, nil
}

func SaveCodeNotActivatedByUser(data UserAuthCode) (UserAuthCode, error) {

	db, err := Connect()
	if err != nil {
		return data, err
	}

	result := db.Save(&data)
	if result.Error != nil {
		return data, result.Error
	}

	return data, nil
}

func NewSession(user User, refreshToken string) (Session, error) {
	session := Session{
		UserId:       user.ID,
		RefreshToken: refreshToken,
		ExpiresIn:    time.Now().UTC().Add(time.Hour * 24 * 30),
	}

	db, err := Connect()
	if err != nil {
		return session, err
	}

	result := db.Create(&session)
	if result.Error != nil {
		return session, result.Error
	}

	return session, nil
}

func DeleteSession(refreshToken string) (bool, error) {
	db, err := Connect()
	if err != nil {
		return false, err
	}

	result := db.Where("refresh_token = ?", refreshToken).Delete(&Session{})
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil

}

func FindSession(refreshToken string) (Session, error) {
	session := Session{}

	db, err := Connect()
	if err != nil {
		return session, err
	}

	result := db.Where("refresh_token = ?", refreshToken).First(&session)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return session, result.Error
	}

	return session, nil
}
