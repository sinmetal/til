package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mercari.io/datastore"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var _ datastore.PropertyTranslator = QueryPlan(nil)
var _ json.Marshaler = QueryPlan(nil)

type QueryPlan []byte

func QueryPlanFromPB(qppb *sppb.QueryPlan) (QueryPlan, error) {
	if qppb == nil {
		return QueryPlan(nil), nil
	}

	return protojson.Marshal(qppb)
}

func (qp QueryPlan) ToQueryPlanPB() (*sppb.QueryPlan, error) {
	if len(qp) == 0 || string(qp) == "null" {
		return nil, nil
	}
	qppb := &sppb.QueryPlan{}
	err := protojson.Unmarshal(qp, qppb)
	if err != nil {
		return nil, err
	}
	return qppb, nil
}

func (qp QueryPlan) MarshalJSON() ([]byte, error) {
	if len(qp) == 0 {
		return []byte("null"), nil
	}
	return json.RawMessage(qp).MarshalJSON()
}

func (qp QueryPlan) ToPropertyValue(ctx context.Context) (interface{}, error) {
	if len(qp) == 0 {
		return "null", nil
	}
	return string(qp), nil
}

func (qp QueryPlan) FromPropertyValue(ctx context.Context, p datastore.Property) (dst interface{}, err error) {
	s, ok := p.Value.(string)
	if !ok {
		return nil, fmt.Errorf("%T is not datastore.Key", p.Value)
	}

	return QueryPlan(s), nil
}
