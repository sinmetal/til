package main

import (
	"context"
	"testing"

	"google.golang.org/api/spanner/v1"
)

func TestOldestActiveQueriesService_ExecuteStreamingSql(t *testing.T) {
	ctx := context.Background()

	spa, err := spanner.NewService(ctx)
	if err != nil {
		t.Fatal(err)
	}
	s := &OldestActiveQueriesService{
		spa: spa,
	}
	session, err := s.CreateSession(ctx)
	if err != nil {
		t.Fatal(err)
	}

	sql := `SELECT * FROM OrderDetail1M`
	if err := s.ExecuteStreamingSql(ctx, session.Name, sql); err != nil {
		t.Fatal(err)
	}
}
