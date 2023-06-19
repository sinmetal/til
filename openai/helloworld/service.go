package helloworld

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type Service struct {
	Client *openai.Client
}

func NewService(ctx context.Context, client *openai.Client) (*Service, error) {
	return &Service{
		Client: client,
	}, nil
}

func (s *Service) CreateImage(ctx context.Context) error {
	resp, err := s.Client.CreateImage(ctx, openai.ImageRequest{
		Prompt:         "Hell's watchdog Cerberus with three heads.\nHis head colors are red, blue and yellow.\nPlease draw in 32bit pixel art.",
		N:              10,
		Size:           "256x256",
		ResponseFormat: "url",
		User:           "",
	})
	if err != nil {
		return err
	}

	for _, v := range resp.Data {
		fmt.Println(v.URL)
	}

	return nil
}

func (s *Service) Hoge() error {
	return nil
}
