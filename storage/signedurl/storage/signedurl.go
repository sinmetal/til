package storage

import (
	"context"
	"encoding/base64"
	"net/url"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/xerrors"
	"google.golang.org/api/iam/v1"
)

// StorageSignedURLService is Storage Signed URL Util Service
type StorageSignedURLService struct {
	ServiceAccountName string
	ServiceAccountID   string
	IAMService         *iam.Service
}

// NewStorageSignedURLService is StorageServiceを生成する
//
// 利用するServiceAccountの roles/iam.serviceAccountTokenCreator https://cloud.google.com/iam/docs/service-accounts?hl=en#the_service_account_token_creator_role を持っている必要がある
// serviceAccountName is SignedURLを発行するServiceAccountの @ より前の値。ex. hoge@projectid.iam.gserviceaccount.com の場合は "hoge"
// serviceAccountID is serviceAccountNameに指定したものと同じServiceAccountのID。format "projects/%s/serviceAccounts/%s"。
// iamService is iamService
func NewStorageSignedURLService(ctx context.Context, serviceAccountName string, serviceAccountID string, iamService *iam.Service) (*StorageSignedURLService, error) {
	return &StorageSignedURLService{
		ServiceAccountName: serviceAccountName,
		ServiceAccountID:   serviceAccountID,
		IAMService:         iamService,
	}, nil
}

// CreateSignedURLForPutObject is ObjectをPutするSignedURLを発行する
// https://cloud.google.com/blog/ja/products/gcp/uploading-images-directly-to-cloud-storage-by-using-signed-url を参考に作られている
func (s *StorageSignedURLService) CreatePutObjectURL(ctx context.Context, bucket string, object string, contentType string, expires time.Time) (string, error) {
	url, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		GoogleAccessID: s.ServiceAccountName,
		Method:         "PUT",
		Expires:        expires,
		ContentType:    contentType,
		Scheme:         storage.SigningSchemeV4,
		// To avoid management for private key, use SignBytes instead of PrivateKey.
		// In this example, we are using the `iam.serviceAccounts.signBlob` API for signing bytes.
		// If you hope to avoid API call for signing bytes every time,
		// you can use self hosted private key and pass it in Privatekey.
		SignBytes: func(b []byte) ([]byte, error) {
			// この関数 deprecated になってるから、移行する必要があるな July 1, 2021
			resp, err := s.IAMService.Projects.ServiceAccounts.SignBlob(
				s.ServiceAccountID,
				&iam.SignBlobRequest{BytesToSign: base64.StdEncoding.EncodeToString(b)},
			).Context(ctx).Do()
			if err != nil {
				return nil, err
			}
			return base64.StdEncoding.DecodeString(resp.Signature)
		},
	})
	if err != nil {
		return "", xerrors.Errorf("failed PutObjectSignedURL: saName=%s,saID=%s,bucket=%s,object=%s : %w", s.ServiceAccountName, s.ServiceAccountID, bucket, object, err)
	}
	return url, nil
}

func (s *StorageSignedURLService) CreateDownloadURL(ctx context.Context, bucket string, object string, queryParameters url.Values, expires time.Time) (string, error) {
	u, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		GoogleAccessID:  s.ServiceAccountName,
		Method:          "GET",
		Expires:         expires,
		Scheme:          storage.SigningSchemeV4,
		QueryParameters: queryParameters,
		// To avoid management for private key, use SignBytes instead of PrivateKey.
		// In this example, we are using the `iam.serviceAccounts.signBlob` API for signing bytes.
		// If you hope to avoid API call for signing bytes every time,
		// you can use self hosted private key and pass it in Privatekey.
		SignBytes: func(b []byte) ([]byte, error) {
			// この関数 deprecated になってるから、移行する必要があるな July 1, 2021
			resp, err := s.IAMService.Projects.ServiceAccounts.SignBlob(
				s.ServiceAccountID,
				&iam.SignBlobRequest{BytesToSign: base64.StdEncoding.EncodeToString(b)},
			).Context(ctx).Do()
			if err != nil {
				return nil, err
			}
			return base64.StdEncoding.DecodeString(resp.Signature)
		},
	})
	if err != nil {
		return "", xerrors.Errorf("failed CreateDownloadURL: saName=%s,saID=%s,bucket=%s,object=%s : %w", s.ServiceAccountName, s.ServiceAccountID, bucket, object, err)
	}
	return u, nil
}
