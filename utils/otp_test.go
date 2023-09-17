package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateOtpUrl(t *testing.T) {
	fmt.Println(GenerateOtpUrl("admin", time.Minute))
}
