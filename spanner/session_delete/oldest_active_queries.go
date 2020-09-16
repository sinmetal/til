package main

import (
	"context"
	"fmt"

	"google.golang.org/api/spanner/v1"
)

type OldestActiveQueriesService struct {
	spa *spanner.Service
}

func NewOldestActiveQueriesService(ctx context.Context) (*OldestActiveQueriesService, error) {
	spa, err := spanner.NewService(ctx)
	if err != nil {
		return nil, err
	}
	return &OldestActiveQueriesService{
		spa: spa,
	}, nil
}

func (s *OldestActiveQueriesService) CreateSession(ctx context.Context) (*spanner.Session, error) {
	session, err := s.spa.Projects.Instances.Databases.Sessions.Create(GCPUGPublicSpannerDatabase, &spanner.CreateSessionRequest{}).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *OldestActiveQueriesService) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := s.spa.Projects.Instances.Databases.Sessions.Delete(fmt.Sprintf("%s/sessions/%s", GCPUGPublicSpannerDatabase, sessionID)).Context(ctx).Do()
	if err != nil {
		return err
	}
	return nil
}

func (s *OldestActiveQueriesService) ExecuteStreamingSql(ctx context.Context, session string, sql string) error {
	result, err := s.spa.Projects.Instances.Databases.Sessions.ExecuteStreamingSql(session, &spanner.ExecuteSqlRequest{
		Sql: sql,
	}).Context(ctx).Do()
	if err != nil {
		return err
	}
	fmt.Printf("result.ServerResponse:%+v\n", result.ServerResponse)
	fmt.Printf("result.Metadata:%+v\n", result.Metadata)
	for _, v := range result.Values {
		fmt.Println(v)
	}

	return nil
}
