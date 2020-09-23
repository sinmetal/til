package storage

import (
	"context"
	"fmt"
	"net/url"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/storage"
	"golang.org/x/xerrors"
	"google.golang.org/api/iam/v1"
	credentialspb "google.golang.org/genproto/googleapis/iam/credentials/v1"
)

// StorageSignedURLService is Storage Signed URL Util Service
type StorageSignedURLService struct {
	ServiceAccountName   string
	ServiceAccountID     string
	IAMService           *iam.Service
	IAMCredentialsClient *credentials.IamCredentialsClient
}

// NewStorageSignedURLService is StorageServiceを生成する
//
// 利用するServiceAccountの roles/iam.serviceAccountTokenCreator https://cloud.google.com/iam/docs/service-accounts?hl=en#the_service_account_token_creator_role を持っている必要がある
// serviceAccountName is SignedURLを発行するServiceAccountの @ より前の値。ex. hoge@projectid.iam.gserviceaccount.com の場合は "hoge"
// serviceAccountID is serviceAccountNameに指定したものと同じServiceAccountのID。format "projects/%s/serviceAccounts/%s"。
// iamService is iamService
func NewStorageSignedURLService(ctx context.Context, serviceAccountName string, serviceAccountID string, iamService *iam.Service, iamCredentialsClient *credentials.IamCredentialsClient) (*StorageSignedURLService, error) {
	return &StorageSignedURLService{
		ServiceAccountName:   serviceAccountName,
		ServiceAccountID:     serviceAccountID,
		IAMService:           iamService,
		IAMCredentialsClient: iamCredentialsClient,
	}, nil
}

// CreateSignedURLForPutObject is ObjectをPutするSignedURLを発行する
//
// https://cloud.google.com/blog/ja/products/gcp/uploading-images-directly-to-cloud-storage-by-using-signed-url を参考に作られている
func (s *StorageSignedURLService) CreatePutObjectURL(ctx context.Context, bucket string, object string, contentType string, contentLength int64, expires time.Time) (string, error) {
	u, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		GoogleAccessID: s.ServiceAccountName,
		Method:         "PUT",
		Expires:        expires,
		ContentType:    contentType,
		Headers:        []string{fmt.Sprintf("Content-Length:%d", contentLength)},
		Scheme:         storage.SigningSchemeV4,
		SignBytes: func(b []byte) ([]byte, error) {
			req := &credentialspb.SignBlobRequest{
				Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", s.ServiceAccountName),
				Payload: b,
			}
			resp, err := s.IAMCredentialsClient.SignBlob(ctx, req)
			if err != nil {
				return nil, err
			}
			return resp.SignedBlob, nil
		},
	})
	if err != nil {
		return "", xerrors.Errorf("failed PutObjectSignedURL: saName=%s,saID=%s,bucket=%s,object=%s : %w", s.ServiceAccountName, s.ServiceAccountID, bucket, object, err)
	}
	return u, nil
}

func (s *StorageSignedURLService) CreateDownloadURL(ctx context.Context, bucket string, object string, queryParameters url.Values, expires time.Time) (string, error) {
	u, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		GoogleAccessID:  s.ServiceAccountName,
		Method:          "GET",
		Expires:         expires,
		Scheme:          storage.SigningSchemeV4,
		QueryParameters: queryParameters,
		SignBytes: func(b []byte) ([]byte, error) {
			req := &credentialspb.SignBlobRequest{
				Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", s.ServiceAccountName),
				Payload: b,
			}
			resp, err := s.IAMCredentialsClient.SignBlob(ctx, req)
			if err != nil {
				return nil, err
			}
			return resp.SignedBlob, nil
		},
	})
	if err != nil {
		return "", xerrors.Errorf("failed CreateDownloadURL: saName=%s,saID=%s,bucket=%s,object=%s : %w", s.ServiceAccountName, s.ServiceAccountID, bucket, object, err)
	}
	return u, nil
}
