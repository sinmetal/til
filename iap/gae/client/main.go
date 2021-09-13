package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/apstndb/adcplus/tokensource"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	// audience には IAPのClientIDを入れる
	ts, err := tokensource.SmartIDTokenSource(ctx, "1057248013684-gjni4id2eisrsesi5b4tv0jm3luo76pp.apps.googleusercontent.com")
	if err != nil {
		panic(err)
	}
	client := oauth2.NewClient(ctx, ts)

	resp, err := client.Get("https://gcpbox-gae.an.r.appspot.com/")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf(string(body))
	fmt.Println()
}
