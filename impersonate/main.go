package main

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	var ts oauth2.TokenSource
	// Google Cloud SDKと同じ環境変数を見ている
	// 環境変数の名前は CLOUDSDK_SECTION_NAME_PROPERTY_NAME というパターンに従う https://cloud.google.com/sdk/docs/properties#setting_properties_via_environment_variables
	// impersonate-service-accountのプロパティはauth/impersonate_service_account https://cloud.google.com/sdk/gcloud/reference#--impersonate-service-account
	isa := os.Getenv("CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT")
	if isa != "" {
		var err error
		ts, err = WithImpersonateTokenSource(ctx, isa, []string{storage.ScopeFullControl})
		if err != nil {
			panic(err)
		}
	}

	// GCS Sample
	gcs, err := NewGCSClient(ctx, ts)
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

// WithImpersonateTokenSource is 指定したService AccountになりすますTokenSourceを作る
// 指定したService Accountに対してroles/iam.serviceAccountTokenCreatorが必要
// https://cloud.google.com/iam/docs/creating-short-lived-service-account-credentials
func WithImpersonateTokenSource(ctx context.Context, serviceAccountEmail string, scopes []string) (oauth2.TokenSource, error) {
	ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: serviceAccountEmail,
		Scopes:          scopes,
	})
	if err != nil {
		return nil, err
	}
	return ts, nil
}

// NewGCSClient is 指定したTokenSourceを設定したGCSClientを作る
func NewGCSClient(ctx context.Context, ts oauth2.TokenSource) (*storage.Client, error) {
	var ops []option.ClientOption
	if ts != nil {
		ops = append(ops, option.WithTokenSource(ts))
	}

	gcs, err := storage.NewClient(ctx, ops...)
	if err != nil {
		return nil, err
	}
	return gcs, nil
}
