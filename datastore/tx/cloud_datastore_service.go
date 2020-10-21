package tx

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/datastore"
)

type DatastoreService struct {
	ds *datastore.Client
}

type Entity struct {
}

func (s *DatastoreService) Kind() string {
	return "TxTest"
}

func (s *DatastoreService) PutMulti(ctx context.Context, keyPrefix string, count int) error {
	var keys []*datastore.Key
	var entities []*Entity
	for i := 0; i < count; i++ {
		keys = append(keys, datastore.NameKey(s.Kind(), fmt.Sprintf("%s%d", keyPrefix, i), nil))
		entities = append(entities, &Entity{})
	}

	_, err := s.ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		_, err := tx.PutMulti(keys, entities)
		if err != nil {
			// return err
			return fmt.Errorf(": %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func UnwrapGRPCError(err error) (*status.Status, bool) {
	cerr := err
	for {
		sts, ok := status.FromError(cerr)
		if ok {
			return sts, true
		}
		nerr := errors.Unwrap(cerr)
		if nerr == nil {
			return nil, false
		}
		cerr = nerr
	}
}
