package utils

import (
	"fmt"
	"testing"
)

func TestGenerateOtpSecret(t *testing.T){
	fmt.Println(GenerateOtpSecret("admin"))
}
