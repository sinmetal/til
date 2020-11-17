package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/sinmetalcraft/silverdile/v2"
	log "github.com/vvakame/sdlog/aelog"
	"golang.org/x/xerrors"
	"google.golang.org/api/idtoken"
)

func main() {
	ctx := context.Background()

	var (
		u      = flag.String("url", "url", "url")
		bucket = flag.String("bucket", "bucket", "bucket")
		object = flag.String("object", "object", "object")
	)
	flag.Parse()

	w := httptest.NewRecorder()

	err := downloadImageFromIronLizard(ctx, *u, w, &silverdile.ImageOption{
		Bucket: *bucket,
		Object: *object,
		Size:   10,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(w.Code)
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

// downloadImageFromIronLizard is IronLizard on Cloud Run から Image を Resize して Download する
func downloadImageFromIronLizard(ctx context.Context, u string, w http.ResponseWriter, imgInfo *silverdile.ImageOption) error {
	serviceURL := fmt.Sprintf("%s/image/resize/%s/%s", u, imgInfo.Bucket, imgInfo.Object)
	if imgInfo.Size > 0 {
		serviceURL = fmt.Sprintf("%s/=s%d", serviceURL, imgInfo.Size)
	}

	tokenSource, err := idtoken.NewTokenSource(ctx, u)
	if err != nil {
		return xerrors.Errorf("failed idtoken.NewTokenSource %#v : %w", imgInfo, err)
	}
	token, err := tokenSource.Token()
	if err != nil {
		return xerrors.Errorf("failed tokenSource.Token() %#v : %w", imgInfo, err)
	}
	fmt.Println(token.Expiry)
	fmt.Println(token.AccessToken)

	fmt.Println(serviceURL)
	req, err := http.NewRequest("GET", serviceURL, nil)
	if err != nil {
		return xerrors.Errorf("failed http.NewRequest %#v : %w", imgInfo, err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return xerrors.Errorf("failed request to IronLizard : %s : %w", serviceURL, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Warningf(ctx, "failed response.body.close from IronLizard. err=%v", err)
		}
	}()
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
