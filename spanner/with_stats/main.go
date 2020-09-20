package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/spanner"
	"github.com/k0kubun/pp"
	"google.golang.org/api/iterator"
)

const DatastoreProjectID = "sinmetal-ci"

var sc *spanner.Client
var ds *datastore.Client

func main() {
	ctx := context.Background()

	var err error
	sc, err = createSpannerClient(ctx, fmt.Sprintf("projects/%s/instances/%s/databases/%s", "gcpug-public-spanner", "merpay-sponsored-instance", "sinmetal_benchmark_b"))
	if err != nil {
		panic(err)
	}
	ds, err = datastore.NewClient(ctx, DatastoreProjectID)
	if err != nil {
		panic(err)
	}

	PrintStats(ctx)
	AnalyzeQuery(ctx)
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

elapsed_time:38.52 secs(string)
query_plan_creation_time:39.04 msecs(string)
statistics_load_time:0(string)
cpu_time:292.82 secs(string)
runtime_creation_time:1.17 secs(string)
deleted_rows_scanned:0(string)
remote_server_calls:4771/4771(string)
bytes_returned:17480(string)
rows_scanned:26018674(string)
optimizer_statistics_package:(string)
query_text:
SELECT OrderId, OrderDetailId, OD.Price, Item.ItemID, Item.Name As ItemName, Item.Price As ItemPrice, OD.CommitedAt
FROM OrderDetail1M OD JOIN Item1K Item ON OD.ItemId = Item.ItemId
ORDER BY ItemId
LIMIT 100
(string)
data_bytes_read:19758968877(string)
rows_returned:100(string)
optimizer_version:2(string)
filesystem_delay_seconds:30.4 secs(string)

*/
func PrintStats(ctx context.Context) {
	sql := `
SELECT OrderId, OrderDetailId, OD.Price, Item.ItemID, Item.Name As ItemName, Item.Price As ItemPrice, OD.CommitedAt
FROM OrderDetail1M OD JOIN Item1K Item ON OD.ItemId = Item.ItemId
ORDER BY ItemId
LIMIT 100
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
	for k, v := range iter.QueryStats {
		fmt.Printf("%s:%v(%T)\n", k, v, v)
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
