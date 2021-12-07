package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	filterbox "github.com/sinmetal/til/pubsub/filter"
)

func main() {
	ctx := context.Background()

	pubsubClient, err := pubsub.NewClient(ctx, "sinmetal-ci")
	if err != nil {
		panic(err)
	}
	s, err := filterbox.NewService(ctx, pubsubClient)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			msg := time.Now().String()
			fmt.Printf("Send Mesasge:%s\n", msg)
			_, err := s.SendMessage(ctx, msg)
			if err != nil {
				fmt.Printf("failed Send Message: %s\n", err)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			msg := fmt.Sprintf("by sinmetal: %s", time.Now().String())
			fmt.Printf("Send Mesasge: %s\n", msg)
			_, err := s.SendMessageBySinmetal(ctx, msg)
			if err != nil {
				fmt.Printf("failed Send Message: %s\n", err)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			msg := fmt.Sprintf("by gold: %s", time.Now().String())
			fmt.Printf("Send Mesasge: %s\n", msg)
			_, err := s.SendMessageOnlyGold(ctx, msg)
			if err != nil {
				fmt.Printf("failed Send Message: %s\n", err)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			msg, err := s.PullMessage(ctx, "non-filter-sub")
			if err != nil {
				fmt.Printf("failed Pull Message: %s\n", err)
			}
			fmt.Printf("non-filter-sub: %s\n", msg)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			msg, err := s.PullMessage(ctx, "filter-on")
			if err != nil {
				fmt.Printf("failed Pull Message: %s\n", err)
			}
			fmt.Printf("filter-on: %s\n", msg)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			msg, err := s.PullMessage(ctx, "filter-not")
			if err != nil {
				fmt.Printf("failed Pull Message: %s\n", err)
			}
			fmt.Printf("filter-not: %s\n", msg)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			msg, err := s.PullMessage(ctx, "filter-only-gold")
			if err != nil {
				fmt.Printf("failed Pull Message: %s\n", err)
			}
			fmt.Printf("filter-only-gold: %s\n", msg)
			time.Sleep(2 * time.Second)
		}
	}()

	var errChan chan error
	<-errChan
}
