package main

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	isa := os.Getenv("CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT")
	ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: isa,
		Scopes:          []string{storage.ScopeFullControl},
	})
	if err != nil {
		panic(err)
	}

	gcs, err := storage.NewClient(ctx, option.WithTokenSource(ts))
	if err != nil {
		panic(err)
	}
	w := gcs.Bucket("hoge").Object(time.Now().String()).NewWriter(ctx)
	_, err = w.Write([]byte("Hello"))
	if err != nil {
		panic(err)
	}
	if err := w.Close(); err != nil {
		panic(err)
	}

}