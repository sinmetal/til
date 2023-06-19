package helloworld_test

import (
	"context"
	"os"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/sinmetal/til/openai/helloworld"
)

func TestService_CreateImage(t *testing.T) {
	ctx := context.Background()

	s := newService(t)

	if err := s.CreateImage(ctx); err != nil {
		t.Fatal(err)
	}
}

func newService(t *testing.T) *helloworld.Service {
	ctx := context.Background()

	apiKey := os.Getenv("OPENAI_APIKEY")
	client := openai.NewClient(apiKey)

	s, err := helloworld.NewService(ctx, client)
	if err != nil {
		t.Fatal(err)
	}
	return s
}
