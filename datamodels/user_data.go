package datamodels

import "errors"

type UserData struct {
}

type UserLoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type ReAuthTokenRequest struct {
	AccessToken string `json:"access_token"`
}

type UserSignUpData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
}

func (signUpDt *UserSignUpData) Validate() error {
	return nil
}

type OtpData struct {
	PassCode    string `json:"passcode"`
	AccessToken string `json:"access_token"`
}
