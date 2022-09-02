package filter_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	filterbox "github.com/sinmetal/til/pubsub/filter"
)

func TestService_SendMessage(t *testing.T) {
	ctx := context.Background()

	pubsubClient, err := pubsub.NewClient(ctx, "sinmetal-ci")
	if err != nil {
		t.Fatal(err)
	}
	s, err := filterbox.NewService(ctx, pubsubClient)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.SendMessage(ctx, time.Now().String())
	if err != nil {
		t.Fatal(err)
	}
}
