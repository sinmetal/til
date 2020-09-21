package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
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
