package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/option"
)

const (
	GCPUGPublicSpannerDatabase = "projects/gcpug-public-spanner/instances/merpay-sponsored-instance/databases/sinmetal_benchmark_b"
)

func main() {
	ctx := context.Background()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	sc, err := createClient(ctx, GCPUGPublicSpannerDatabase)
	if err != nil {
		panic(err)
	}

	//go func() {
	//	spa, err := spaaaa.NewService(ctx)
	//	if err != nil {
	//		panic(err)
	//	}
	//	s := &OldestActiveQueriesService{
	//		spa: spa,
	//	}
	//	session, err := s.CreateSession(ctx)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("session:%+v\n", session)
	//
	//	sql := `SELECT * FROM OrderDetail1M`
	//	if err := s.ExecuteStreamingSql(ctx, session.Name, sql); err != nil {
	//		panic(err)
	//	}
	//}()

	go func() {
		sql := `SELECT "Single", OrderId, OrderDetailId, OD.Price, Item.ItemID, Item.Name As ItemName, Item.Price As ItemPrice, OD.CommitedAt
FROM OrderDetail1M OD JOIN Item1K Item ON OD.ItemId = Item.ItemId
ORDER BY ItemId
LIMIT 100000`

		oaqs, err := NewOldestActiveQueriesService(ctx)
		if err != nil {
			panic(err)
		}
		ss := NewSpannerService(ctx, sc, oaqs)
		for {
			if err := ss.RunSingleQuery(ctx, spanner.NewStatement(sql)); err != nil {
				fmt.Printf("failed RunSingleQuery() err=%+v\n", err)
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	go func() {
		sql := `SELECT "RWTx", OrderId, OrderDetailId, OD.Price, Item.ItemID, Item.Name As ItemName, Item.Price As ItemPrice, OD.CommitedAt
FROM OrderDetail1M OD JOIN Item1K Item ON OD.ItemId = Item.ItemId
ORDER BY ItemId
LIMIT 100000`

		oaqs, err := NewOldestActiveQueriesService(ctx)
		if err != nil {
			panic(err)
		}
		ss := NewSpannerService(ctx, sc, oaqs)
		for {
			if err := ss.RunRWTxQuery(ctx, spanner.NewStatement(sql)); err != nil {
				fmt.Printf("failed RunRWTxQuery() err=%+v\n", err)
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	go func() {
		time.Sleep(60 * time.Second)
		sc, err := createClient(ctx, GCPUGPublicSpannerDatabase)
		if err != nil {
			panic(err)
		}
		oaqs, err := NewOldestActiveQueriesService(ctx)
		if err != nil {
			panic(err)
		}
		ss := NewSpannerService(ctx, sc, oaqs)

		for {
			if err := ss.QueryOldestActiveQueries(ctx); err != nil {
				fmt.Printf("failed QueryOldestActiveQueries() err=%+v\n", err)
			}
			time.Sleep(60 * time.Second)
		}
	}()

	sig := <-quit

	sc.Close()
	fmt.Printf("\n%v : Spanner.Client.Close()\n", sig)
}

func createClient(ctx context.Context, db string, o ...option.ClientOption) (*spanner.Client, error) {
	config := spanner.ClientConfig{
		SessionPoolConfig: spanner.SessionPoolConfig{
			MinOpened:           10,
			TrackSessionHandles: true,
		},
	}
	dataClient, err := spanner.NewClientWithConfig(ctx, db, config, o...)
	if err != nil {
		return nil, err
	}

	return dataClient, nil
}
