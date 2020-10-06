package serviceusage_test

import (
	"context"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/sinmetal/til/serviceusage"
	orgsus "google.golang.org/api/serviceusage/v1"
)

func TestServiceUsageService_ListAll(t *testing.T) {
	ctx := context.Background()

	sus := newTestServiceUsageService(t)

	_, err := sus.ListAll(ctx, 401580979819)
	if err != nil {
		t.Fatal(err)
	}
}

func TestServiceUsageService_ListByDiff(t *testing.T) {
	ctx := context.Background()

	sus := newTestServiceUsageService(t)

	list, err := sus.ListByDiff(ctx, 401580979819, 942245768052)
	if err != nil {
		t.Fatal(err)
	}
	pp.Println(list)
}

func newTestServiceUsageService(t *testing.T) *serviceusage.ServiceUsageService {
	ctx := context.Background()

	orgService, err := orgsus.NewService(ctx)
	if err != nil {
		t.Fatal(err)
	}
	sus, err := serviceusage.NewService(ctx, orgService)
	if err != nil {
		t.Fatal(err)
	}
	return sus
}
