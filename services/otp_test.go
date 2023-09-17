package services

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func TestABCD(t *testing.T) {
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "abc",
		AccountName: "def",
	})
	url := key.URL()
	key1, _ := otp.NewKeyFromURL(url)
	img, _ := key1.Image(300, 300)
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Error(err)
	}
	if err := os.WriteFile("abcd.png", buf.Bytes(), 0666); err != nil {
		t.Error(err)
	}
}
