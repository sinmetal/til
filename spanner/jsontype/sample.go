package jsontype

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

// InsertJsonSample1 is interface{}突っ込んだら、雑にJsonにしてくれたりしないかなと思ったけど、ダメだった
// panic: spanner: code = "InvalidArgument", desc = "client doesn't support type *main.JsonBody"
func InsertJsonSample1(ctx context.Context) error {
	spa, err := spanner.NewClient(ctx, "projects/gcpug-public-spanner/instances/merpay-sponsored-instance/databases/sinmetal")
	if err != nil {
		return err
	}
	defer spa.Close()

	type JsonBody struct {
		ID    string
		Count int64
		Date  time.Time
	}
	body := &JsonBody{
		ID:    "helloJson",
		Count: 100,
		Date:  time.Now(),
	}

	type Entity struct {
		ID   string
		Json interface{}
	}

	e := &Entity{
		ID:   "sample",
		Json: body,
	}

	is, err := spanner.InsertOrUpdateStruct("JsonSample", e)
	if err != nil {
		return err
	}
	_, err = spa.Apply(ctx, []*spanner.Mutation{is})
	if err != nil {
		return err
	}
	return nil
}

// InsertJsonSample1 is interface{}突っ込んだら、雑にJsonにしてくれたりしないかなと思ったけど、ダメだった
// panic: spanner: code = "InvalidArgument", desc = "client doesn't support type *main.JsonBody"
func InsertJsonSample2(ctx context.Context) error {
	spa, err := spanner.NewClient(ctx, "projects/gcpug-public-spanner/instances/merpay-sponsored-instance/databases/sinmetal")
	if err != nil {
		return err
	}
	defer spa.Close()

	type JsonBody struct {
		ID    string
		Count int64
		Date  time.Time
	}
	body := &JsonBody{
		ID:    "helloJson",
		Count: 100,
		Date:  time.Now(),
	}

	type Entity struct {
		ID   string
		Json []byte
	}
	j, err := json.Marshal(body)
	if err != nil {
		return err
	}

	e := &Entity{
		ID:   "sample",
		Json: j,
	}

	is, err := spanner.InsertOrUpdateStruct("JsonSample", e)
	if err != nil {
		return err
	}
	_, err = spa.Apply(ctx, []*spanner.Mutation{is})
	if err != nil {
		return err
	}
	return nil
}

// InsertJsonSample1 is interface{}突っ込んだら、雑にJsonにしてくれたりしないかなと思ったけど、ダメだった
// panic: spanner: code = "InvalidArgument", desc = "client doesn't support type *main.JsonBody"
func InsertJsonSample3(ctx context.Context) error {
	spa, err := spanner.NewClient(ctx, "projects/gcpug-public-spanner/instances/merpay-sponsored-instance/databases/sinmetal")
	if err != nil {
		return err
	}
	defer spa.Close()

	const tableName = "JsonSample"
	const id = "sample"
	columns := []string{"ID", "Json"}

	type JsonBody struct {
		ID    string
		Count int64
		Date  time.Time
	}
	body := &JsonBody{
		ID:    "helloJson",
		Count: 100,
		Date:  time.Date(2021, 8, 30, 13, 0, 0, 0, time.UTC),
	}

	type Entity struct {
		ID   string
		Json spanner.NullJSON
	}

	e := &Entity{
		ID:   id,
		Json: spanner.NullJSON{Value: *body, Valid: true},
	}

	is, err := spanner.InsertOrUpdateStruct(tableName, e)
	if err != nil {
		return err
	}
	_, err = spa.Apply(ctx, []*spanner.Mutation{is})
	if err != nil {
		return err
	}

	iter := spa.Single().Read(ctx, tableName, spanner.KeySets(spanner.Key{id}), columns)
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		var e Entity
		if err := row.ToStruct(&e); err != nil {
			return err
		}
		fmt.Printf("%#v", e)
	}
	return nil
}
