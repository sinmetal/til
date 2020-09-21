package storage

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"github.com/google/uuid"
	"google.golang.org/api/iam/v1"
)

func TestStorageSignedURLService_CreatePutObjectURL(t *testing.T) {
	ctx := context.Background()

	signedURLService := createStorageSignedURLService(t)

	object := uuid.New().String()
	u, err := signedURLService.CreatePutObjectURL(ctx, "sinmetal-ci-signed-url", object, "image/jpg", time.Now().Add(10*time.Minute))
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile("/Users/sinmetal/go/src/github.com/sinmetal/til/storage/signedurl/sinmetal.jpg")
	if err != nil {
		t.Fatal(err)
	}
	// Generates *http.Request to request with PUT method to the Signed URL.
	req, err := http.NewRequest("PUT", u, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "image/jpg")
	// req.Header.Add("Content-Type", "image/jpg")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp.Status)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(respBody))
}

func TestStorageSignedURLService_CreateDownloadURL(t *testing.T) {
	ctx := context.Background()

	signedURLService := createStorageSignedURLService(t)

	const bucket = "sinmetal-ci-signed-url"
	object := uuid.New().String()
	putURL, err := signedURLService.CreatePutObjectURL(ctx, bucket, object, "image/jpg", time.Now().Add(10*time.Minute))
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile("/Users/sinmetal/go/src/github.com/sinmetal/til/storage/signedurl/sinmetal.jpg")
	if err != nil {
		t.Fatal(err)
	}
	// Generates *http.Request to request with PUT method to the Signed URL.
	req, err := http.NewRequest("PUT", putURL, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	const contentType = "image/jpg"
	req.Header.Set("Content-Type", contentType)
	client := new(http.Client)
	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	// ここまで File Upload

	responseContentDisposition := fmt.Sprintf(`attachment;filename*=UTF-8''%s`, url.PathEscape("いやっほー .jpeg"))
	qp := url.Values{}
	qp.Add("response-content-disposition", responseContentDisposition)
	qp.Add("response-content-type", contentType)
	dlURL, err := signedURLService.CreateDownloadURL(ctx, bucket, object, qp, time.Now().Add(10*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dlURL)
}

func createStorageSignedURLService(t *testing.T) *StorageSignedURLService {
	ctx := context.Background()

	iamService, err := iam.NewService(ctx)
	if err != nil {
		t.Fatal(err)
	}

	iamCredentialsClient, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	signedURLService, err := NewStorageSignedURLService(ctx, "signedurl@sinmetal-ci.iam.gserviceaccount.com", fmt.Sprintf("projects/%s/serviceAccounts/%s", "sinmetal-ci", "signedurl@sinmetal-ci.iam.gserviceaccount.com"), iamService, iamCredentialsClient)
	if err != nil {
		t.Fatal(err)
	}

	return signedURLService
}
