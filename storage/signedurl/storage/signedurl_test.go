package storage

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"github.com/google/uuid"
	"google.golang.org/api/iam/v1"
)

func TestStorageSignedURLService_CreatePutObjectURL(t *testing.T) {
	ctx := context.Background()

	signedURLService := createStorageSignedURLService(t)

	b, fi := openFile(t, "/Users/sinmetal/go/src/github.com/sinmetal/til/storage/signedurl/sinmetal.jpg")

	object := uuid.New().String()
	u, err := signedURLService.CreatePutObjectURL(ctx, "sinmetal-ci-signed-url", object, fi.ContentType, fi.Size, time.Now().Add(10*time.Minute))
	if err != nil {
		panic(err)
	}

	// Generates *http.Request to request with PUT method to the Signed URL.
	req, err := http.NewRequest("PUT", u, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", fi.ContentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", fi.Size))
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
	b, fi := openFile(t, "/Users/sinmetal/go/src/github.com/sinmetal/til/storage/signedurl/sinmetal.jpg")

	object := uuid.New().String()
	u, err := signedURLService.CreatePutObjectURL(ctx, "sinmetal-ci-signed-url", object, fi.ContentType, fi.Size, time.Now().Add(10*time.Minute))
	if err != nil {
		panic(err)
	}

	// Generates *http.Request to request with PUT method to the Signed URL.
	req, err := http.NewRequest("PUT", u, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", fi.ContentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", fi.Size))
	client := new(http.Client)
	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	// ここまで File Upload

	responseContentDisposition := fmt.Sprintf(`attachment;filename*=UTF-8''%s`, url.PathEscape("いやっほー .jpeg"))
	qp := url.Values{}
	qp.Add("response-content-disposition", responseContentDisposition)
	qp.Add("response-content-type", fi.ContentType)
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

type FileInfo struct {
	Name        string
	ContentType string
	Size        int64
}

func openFile(t *testing.T, path string) ([]byte, *FileInfo) {
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}

	// Reset the read pointer if necessary.
	_, err = file.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	var ret FileInfo
	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	ret.ContentType = contentType

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	ret.Name = fi.Name()
	ret.Size = fi.Size()

	body, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	return body, &ret
}
