package datamodels

type UserData struct {
	Code            string      `json:"code"`
	Username        string      `json:"username"`
	IsAuthenticator bool        `json:"is_authenticator"`
	Is2FA           bool        `json:"is_2fa"`
	PartnerData     PartnerData `json:"partner_data"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type TokenAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Image       string `json:"image"`
	Key         string `json:"key"`
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
