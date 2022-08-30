package filter

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type Service struct {
	pubsubClient *pubsub.Client
}

func NewService(ctx context.Context, pubsubClient *pubsub.Client) (*Service, error) {
	return &Service{
		pubsubClient: pubsubClient,
	}, nil
}

func (s *Service) SendMessage(ctx context.Context, message string) (string, error) {
	ret := s.pubsubClient.TopicInProject("filtertest", "sinmetal-ci").Publish(ctx, &pubsub.Message{
		Data: []byte(message),
	})
	id, err := ret.Get(ctx)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Service) SendMessageBySinmetal(ctx context.Context, message string) (string, error) {
	ret := s.pubsubClient.TopicInProject("filtertest", "sinmetal-ci").Publish(ctx, &pubsub.Message{
		Data: []byte(message),
		Attributes: map[string]string{
			"domain": "sinmetal.jp",
		},
	})
	id, err := ret.Get(ctx)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Service) SendMessageOnlyGold(ctx context.Context, message string) (string, error) {
	ret := s.pubsubClient.TopicInProject("filtertest", "sinmetal-ci").Publish(ctx, &pubsub.Message{
		Data: []byte(message),
		Attributes: map[string]string{
			"gold": "",
		},
	})
	id, err := ret.Get(ctx)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Service) PullMessage(ctx context.Context, subscriber string) (msg string, err error) {
	err = s.pubsubClient.SubscriptionInProject(subscriber, "sinmetal-ci").Receive(ctx, func(ctx context.Context, message *pubsub.Message) {
		message.Ack()
		msg = string(message.Data)
		fmt.Printf("%s: %s\n", subscriber, msg)
	})
	if err != nil {
		err = err
	}
	return msg, err
}
