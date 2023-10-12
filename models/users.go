package models

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/utils"
	"gorm.io/gorm"
)

const (
	USER_CODE_LENGTH int = 5
)

type User struct {
	DefaultModel
	Code            string  `gorm:"uniqueIndex; index; not null; size:10;"`
	Username        string  `gorm:"uniqueIndex; index; not null; size:32;"`
	Password        []byte  `gorm:"size:256;"`
	PartnerID       uint    `gorm:"uniqueIndex; not null; index"`
	Partner         Partner `gorm:"foreignKey:PartnerID"`
	OtpUrl          string  `gorm:"size:256; index;"`
	IsAuthenticator bool    `gorm:"default:false;"`
}

func (user *User) GetSequence() string {
	return "user_sequence"
}

func (user *User) NextCode(db *gorm.DB) string {
	var nextSequence int
	db.Raw(fmt.Sprintf("SELECT nextval('%v');", user.GetSequence())).Scan(&nextSequence)
	return UuidCode.ConvertCode(nextSequence, USER_CODE_LENGTH)
}

func (user *User) Create(tx *gorm.DB) error {
	if strings.TrimSpace(user.Username) == "" {
		return utils.ErrInvalid
	}
	user.Code = user.NextCode(tx)
	return tx.Create(user).Error
}

func (user *User) Authenticate(tx *gorm.DB, userLoginData *datamodels.UserLoginData) error {
	if err := tx.Where(&User{
		Username: userLoginData.Username,
	}).First(user).Error; err != nil {
		return err
	}
	if user.Username == "" {
		return errors.New("user not found")
	}
	if !bytes.Equal(user.Password, utils.Hash256(userLoginData.Password)) {
		return errors.New("wrong password")
	}
	return nil
}

func (user *User) GetPartnerByUsername(username string, tx *gorm.DB) (*Partner, error) {
	partner := Partner{}
	if err := tx.Where(&User{
		Username: username,
	}).First(user).Error; err != nil {
		return &partner, err
	}
	if err := tx.First(&partner, user.PartnerID).Error; err != nil {
		return &partner, err
	}
	return &partner, nil
}

func (user *User) GetUserByUsername(username string, tx *gorm.DB) error {
	return tx.Where(&User{
		Username: username,
	}).First(user).Error
}

func (user *User) Login() error {
	return nil
}

func (user *User) SetOtpUrl(url string, tx *gorm.DB) error {
	user.OtpUrl = url
	return tx.Save(user).Error
}

func (user *User) ChangePassword(newPassword string) error {
	return nil
}

func (user *User) ChangeEmail(newEmail string) error {
	return nil
}

func (user *User) ChangePhone(newPhone string) error {
	return nil
}

func (user *User) SetIsAuthenticator(input bool, tx *gorm.DB) error {
	user.IsAuthenticator = input
	return tx.Save(user).Error
}
