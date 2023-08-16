package datamodels

import "errors"

type UserData struct {
}

type UserLoginData struct {
	Username string
	Password string
}

func (userLoginData *UserLoginData) Validate() error {
	if userLoginData.Username == "" {
		return errors.New("wrong username")
	}
	if userLoginData.Password == "" {
		return errors.New("wrong password")
	}
	return nil
}

type TokenResponse struct {
	AccessToken string
	TokenType   string
}

type ReAuthTokenRequest struct {
	Token string
}

type UserSignUpData struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	UserName  string
	Password  string
}

func (signUpDt *UserSignUpData) Validate() error {
	return nil
}
