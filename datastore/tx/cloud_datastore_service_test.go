package tx

import (
	"context"
	"testing"

	"cloud.google.com/go/datastore"
	"google.golang.org/grpc/codes"
)

func TestDatastoreService_PutMultiByLegacy(t *testing.T) {
	ctx := context.Background()

	legacyDS, err := datastore.NewClient(ctx, "gcpugjp-dev")
	if err != nil {
		t.Fatal(err)
	}

	dss := &DatastoreService{
		legacyDS,
	}

	cases := []struct {
		name  string
		count int
		want  codes.Code
	}{
		{"normal", 25, 0},
		{"too many entity group error", 26, codes.InvalidArgument},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := dss.PutMulti(ctx, "k", tt.count)
			if tt.want != 0 {
				sts, ok := UnwrapGRPCError(err)
				if !ok {
					t.Errorf("want status.Status but got %T%v", err, err)
				}
				if e, g := tt.want, sts.Code(); e != g {
					t.Errorf("want error.code %v but got %v", e, g)
				}
			}
		})
	}
}
