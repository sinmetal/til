package txnest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

func TestNestTX(t *testing.T) {
	ctx := context.Background()

	client, err := spanner.NewClientWithConfig(ctx, "projects/gcpug-public-spanner/instances/merpay-sponsored-instance/databases/sinmetal", spanner.ClientConfig{
		SessionPoolConfig: spanner.SessionPoolConfig{
			MaxOpened:           10,
			MinOpened:           1,
			MaxIdle:             0,
			MaxBurst:            0,
			WriteSessions:       0,
			HealthCheckWorkers:  0,
			HealthCheckInterval: 0,
			TrackSessionHandles: false,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, rwtx *spanner.ReadWriteTransaction) error {
		iter := client.Single().Query(context.Background(), spanner.NewStatement("SELECT 1"))
		defer iter.Stop()
		for {
			row, err := iter.Next()
			if err == iterator.Done {
				break
			} else if err != nil {
				return fmt.Errorf("failed read : %w", err)
			}
			v := row.String()
			fmt.Printf("row:%s\n", v)
		}

		mu := spanner.Insert("Session", []string{"Id"}, []interface{}{uuid.New().String()})
		if err := rwtx.BufferWrite([]*spanner.Mutation{mu}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

type AccessCounter struct {
	ID         string
	Count      int64
	CommitedAt time.Time
}

func TestTx(t *testing.T) {
	ctx := context.Background()

	client, err := spanner.NewClientWithConfig(ctx, "projects/gcpug-public-spanner/instances/merpay-sponsored-instance/databases/sinmetal", spanner.ClientConfig{
		SessionPoolConfig: spanner.SessionPoolConfig{
			MaxOpened:           10,
			MinOpened:           1,
			MaxIdle:             0,
			MaxBurst:            0,
			WriteSessions:       0,
			HealthCheckWorkers:  0,
			HealthCheckInterval: 0,
			TrackSessionHandles: false,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	CounterID := "10f7b715-21ae-47f3-b5a7-f6fcf886c59b"
	allColumns := []string{"ID", "Count", "CommitedAt"}
	const table = "AccessCounter"

	rotx1 := client.ReadOnlyTransaction()
	{
		iter := rotx1.Read(ctx, table, spanner.KeySets(spanner.Key{CounterID}), allColumns)
		defer iter.Stop()
		row, err := iter.Next()
		if err != nil {
			t.Fatal(err)
		}
		var v1 AccessCounter
		if err := row.ToStruct(&v1); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("v1.Count:%d\n", v1.Count)
	}

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, rwtx *spanner.ReadWriteTransaction) error {
		iter := rwtx.Read(ctx, "AccessCounter", spanner.KeySets(spanner.Key{CounterID}), allColumns)
		defer iter.Stop()

		row, err := iter.Next()
		if err != nil {
			return err
		}
		var v AccessCounter
		if err := row.ToStruct(&v); err != nil {
			return err
		}
		fmt.Printf("write before:%d\n", v.Count)
		v.Count++

		mu, err := spanner.UpdateStruct("AccessCounter", &v)
		if err != nil {
			return err
		}
		// mu := spanner.Update("AccessCounter", []string{"ID", "Count", "CommitedAt"}, []interface{}{uuid.New().String(), 10, spanner.CommitTimestamp})
		if err := rwtx.BufferWrite([]*spanner.Mutation{mu}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	{
		iter := rotx1.Read(ctx, table, spanner.KeySets(spanner.Key{CounterID}), allColumns)
		defer iter.Stop()
		row, err := iter.Next()
		if err != nil {
			t.Fatal(err)
		}
		var v1 AccessCounter
		if err := row.ToStruct(&v1); err != nil {
			t.Fatal(err)
		}
		fmt.Printf("v2.Count:%d\n", v1.Count)
	}
	rotx1.Close()
}
