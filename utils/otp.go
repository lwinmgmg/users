package utils

import (
	"github.com/pquerna/otp/totp"
)

func GenerateOtpSecret(username string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "user",
		AccountName: username,
	})
	if err != nil {
		return "", err
	}
	return key.Secret(), err
}
