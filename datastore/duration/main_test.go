package main

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
)

type Duration struct {
	Du time.Duration `datastore:",noindex"` // time.Durationは 経過時間を nanosecound count (int64) で保持しているので、Datastore上はIntegerで保存される
}

// TestDurationToDatastore
// DurationってDatastoreに普通に入るのか？と、ふと思ってやってみたら、普通に入った
func TestDurationToDatastore(t *testing.T) {
	ctx := context.Background()

	ds, err := datastore.NewClient(ctx, "sinmetal-ci")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := ds.Close(); err != nil {
			t.Log(err)
		}
	}()
	now := time.Now()

	key, err := ds.Put(ctx, datastore.NameKey("Duration", "test", nil), &Duration{
		Du: time.Since(now),
	})
	if err != nil {
		t.Fatal(err)
	}

	var stored Duration
	if err := ds.Get(ctx, key, &stored); err != nil {
		t.Fatal(err)
	}
	t.Log(stored.Du) // 183ns と出てくる
}
