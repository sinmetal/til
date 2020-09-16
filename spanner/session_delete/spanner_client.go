package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

type SpannerService struct {
	sc   *spanner.Client
	oaqs *OldestActiveQueriesService
}

func NewSpannerService(ctx context.Context, spannerClient *spanner.Client, oaqs *OldestActiveQueriesService) *SpannerService {
	return &SpannerService{
		sc:   spannerClient,
		oaqs: oaqs,
	}
}

func (s *SpannerService) RunSingleQuery(ctx context.Context, statement spanner.Statement) error {
	//ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	//defer cancel()

	fmt.Println("Start RunSingleQuery()")
	iter := s.sc.Single().Query(ctx, statement)
	defer iter.Stop()

	var rowCount int64
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			fmt.Printf("Finish RunSingleQuery() RowCount:%d, ElapsedTime:%v\n", rowCount, iter.QueryStats["elapsed_time"])
			return nil
		}
		if err != nil {
			return err
		}
		rowCount++
	}
}

func (s *SpannerService) RunRWTxQuery(ctx context.Context, statement spanner.Statement) error {
	fmt.Println("Start RunRWTxQuery()")
	_, err := s.sc.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		iter := tx.QueryWithStats(ctx, statement)
		defer iter.Stop()

		var rowCount int64
		for {
			_, err := iter.Next()
			if err == iterator.Done {
				fmt.Printf("Finish RunRWTxQuery() RowCount:%d, ElapsedTime:%v\n", rowCount, iter.QueryStats["elapsed_time"])
				return nil
			}
			if err != nil {
				return err
			}
			rowCount++
		}
	})
	return err
}

type OldestActiveQuery struct {
	StartTime time.Time `spanner:"start_time"`
	Text      string    `spanner:"text"`
	SessionID string    `spanner:"session_id"`
}

func (s *SpannerService) QueryOldestActiveQueries(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	const sql = `
SELECT start_time,
       text,
       session_id
FROM spanner_sys.oldest_active_queries
ORDER BY start_time ASC;
`
	iter := s.sc.Single().Query(ctx, spanner.NewStatement(sql))
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}

		var v OldestActiveQuery
		if err := row.ToStruct(&v); err != nil {
			return err
		}
		if strings.Contains(v.Text, "spanner_sys") {
			continue
		}
		if err := s.oaqs.DeleteSession(ctx, v.SessionID); err != nil {
			return err
		}
		fmt.Printf("DeleteSession:%s, text=%s\n", v.SessionID, strings.ReplaceAll(v.Text, "\n", " "))
	}
}
