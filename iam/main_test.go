package iam

import (
	"context"
	"os"
	"testing"
)

func TestGetSignedJwt(t *testing.T) {
	sa := os.Getenv("SA")
	GetSignedJwt(context.Background(), sa)
}
