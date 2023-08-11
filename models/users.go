package models

import (
	"bytes"
	"errors"

	"github.com/lwinmgmg/user/datamodels"
	"github.com/lwinmgmg/user/utils"
	"gorm.io/gorm"
)

type User struct {
	DefaultModel
	Username  string  `gorm:"uniqueIndex; not null; size:32;"`
	Password  []byte  `gorm:"size:256;"`
	PartnerID uint    `gorm:"uniqueIndex; not null;"`
	Partner   Partner `gorm:"foreignKey:PartnerID"`
}

func (user *User) Create(tx *gorm.DB) error {
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

func (user *User) Exist(tx *gorm.DB) error {
	if err := tx.Where(&User{
		Username: user.Username,
	}).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return utils.ErrRecordAlreadyExist
}

func (user *User) Login() error {
	return nil
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
