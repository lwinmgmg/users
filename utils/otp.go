package utils

import (
	"time"

	"github.com/pquerna/otp/totp"
)

func GenerateOtpUrl(username string, duration time.Duration) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "user",
		AccountName: username,
		Period:      uint(duration.Seconds()),
	})
	if err != nil {
		return "", err
	}
	return key.URL(), err
}
