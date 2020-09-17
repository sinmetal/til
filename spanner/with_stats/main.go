package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/k0kubun/pp"
	"google.golang.org/api/iterator"
)

var sc *spanner.Client

func main() {
	ctx := context.Background()

	var err error
	sc, err = createSpannerClient(ctx, fmt.Sprintf("projects/%s/instances/%s/databases/%s", "gcpug-public-spanner", "merpay-sponsored-instance", "sinmetal"))
	if err != nil {
		panic(err)
	}

	PrintStats(ctx)
}

// PrintStats is QueryWithStatsのStatsを見る
/*
map[string]interface {}{
  "elapsed_time":                 "36.85 msecs",
  "cpu_time":                     "34.84 msecs",
  "query_plan_creation_time":     "34.64 msecs",
  "deleted_rows_scanned":         "0",
  "optimizer_version":            "2",
  "remote_server_calls":          "0/0",
  "rows_returned":                "1",
  "rows_scanned":                 "0",
  "query_text":                   "SELECT 1",
  "runtime_creation_time":        "0 msecs",
  "optimizer_statistics_package": "",
  "filesystem_delay_seconds":     "0 msecs",
  "bytes_returned":               "8",
}
*/
func PrintStats(ctx context.Context) {
	sql := `
SELECT 1
`
	iter := sc.Single().QueryWithStats(ctx, spanner.NewStatement(sql))
	defer iter.Stop()
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			panic(err)
		}
	}

	_, err := pp.Println(iter.QueryStats)
	if err != nil {
		panic(err)
	}
}

func createSpannerClient(ctx context.Context, database string) (*spanner.Client, error) {
	cli, err := spanner.NewClientWithConfig(ctx, database, spanner.ClientConfig{
		SessionPoolConfig: spanner.SessionPoolConfig{
			MinOpened:     1,  // 1query投げておしまいので、1でOK
			MaxOpened:     10, // 1query投げておしまいなので、そんなにたくさんは要らない
			WriteSessions: 0,  // Readしかしないので、WriteSessionsをPoolする必要はない
		},
	})
	if err != nil {
		return nil, err
	}

	return cli, nil
}
