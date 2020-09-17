package main

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/k0kubun/pp"
	"go.mercari.io/datastore/clouddatastore"
)

type TilSpannerProfileV1 struct {
	QueryProfileV1 QueryPlan
}

type TilSpannerProfileV2 struct {
	QueryProfileV2 QueryPlan
}

func AnalyzeQuery(ctx context.Context) {
	sql := `
SELECT 1
`

	plan, err := sc.Single().AnalyzeQuery(ctx, spanner.NewStatement(sql))
	if err != nil {
		panic(err)
	}

	mds, err := clouddatastore.FromClient(ctx, ds)
	if err != nil {
		panic(err)
	}
	qp, err := QueryPlanFromPB(plan)
	if err != nil {
		panic(err)
	}

	e := TilSpannerProfileV1{
		QueryProfileV1: qp,
	}
	key, err := mds.Put(ctx, mds.NameKey("TilSpannerProfile", "v1", nil), &e)
	if err != nil {
		panic(err)
	}

	{
		var stored TilSpannerProfileV1
		if err := mds.Get(ctx, key, &stored); err != nil {
			panic(err)
		}
		_, err = pp.Println(stored)
		if err != nil {
			panic(err)
		}
	}
	{
		// mercari/datastoreって存在しないPropertyは無視でいいんだっけ？と思って、試してみた
		// 無視してた
		var stored TilSpannerProfileV2
		if err := mds.Get(ctx, key, &stored); err != nil {
			panic(err)
		}
		_, err = pp.Println(stored)
		if err != nil {
			panic(err)
		}
	}
}
