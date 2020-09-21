package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"google.golang.org/api/iam/v1"
)

func TestStorageSignedURLService_CreatePutObjectURL(t *testing.T) {
	ctx := context.Background()

	iamService, err := iam.NewService(ctx)
	if err != nil {
		panic(err)
	}
	signedURLService, err := NewStorageSignedURLService(ctx, "signedurl@sinmetal-ci.iam.gserviceaccount.com", fmt.Sprintf("projects/%s/serviceAccounts/%s", "sinmetal-ci", "signedurl@sinmetal-ci.iam.gserviceaccount.com"), iamService)
	if err != nil {
		panic(err)
	}

	object := uuid.New().String()
	url, err := signedURLService.CreatePutObjectURL(ctx, "sinmetal-ci-signed-url", object, "image/jpg", time.Now().Add(10*time.Minute))
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile("/Users/sinmetal/go/src/github.com/sinmetal/til/storage/signedurl/sinmetal.jpg")
	if err != nil {
		t.Fatal(err)
	}
	// Generates *http.Request to request with PUT method to the Signed URL.
	req, err := http.NewRequest("PUT", url, bytes.NewReader(b))
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

	iamService, err := iam.NewService(ctx)
	if err != nil {
		panic(err)
	}
	signedURLService, err := NewStorageSignedURLService(ctx, "signedurl@sinmetal-ci.iam.gserviceaccount.com", fmt.Sprintf("projects/%s/serviceAccounts/%s", "sinmetal-ci", "signedurl@sinmetal-ci.iam.gserviceaccount.com"), iamService)
	if err != nil {
		panic(err)
	}

	const bucket = "sinmetal-ci-signed-url"
	object := uuid.New().String()
	putURL, err := signedURLService.CreatePutObjectURL(ctx, bucket, object, "image/jpg", time.Now().Add(10*time.Minute))
	if err != nil {
		panic(err)
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
	req.Header.Set("Content-Type", "image/jpg")
	// req.Header.Add("Content-Type", "image/jpg")
	client := new(http.Client)
	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	// ここまで File Upload

	headers := []string{}
	dlURL, err := signedURLService.CreateDownloadURL(ctx, bucket, object, headers, time.Now().Add(10*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dlURL)

	// 以下は Content-Disposition, Content-Type を設定するバージョン
	downloadURL, err := url.Parse(dlURL)
	if err != nil {
		t.Fatal(err)
	}
	vs, err := url.ParseQuery(downloadURL.RawQuery)
	if err != nil {
		t.Fatal(err)
	}
	vs.Add("response-content-disposition", fmt.Sprintf(`attachment; filename*=UTF-8''%s`, url.PathEscape("いえーいふぁいる")))
	vs.Add("response-content-type", "image/jpg")
	downloadURL.RawQuery = vs.Encode()
	t.Log(downloadURL.String())
}
