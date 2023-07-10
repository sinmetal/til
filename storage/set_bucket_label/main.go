package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/apstndb/adcplus/tokensource"
	"github.com/k0kubun/pp/v3"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"
)

func main() {
	ctx := context.Background()

	var projectID string
	if len(os.Args) > 1 {
		projectID = os.Args[1]
	}
	if len(projectID) < 1 {
		log.Fatalln("ProjectID required")
	}
	fmt.Printf("projectID=%s\n\n", projectID)

	ts, err := tokensource.SmartAccessTokenSource(ctx)
	if err != nil {
		panic(err)
	}

	gcsClient, err := storage.NewClient(ctx, option.WithTokenSource(ts))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := gcsClient.Close(); err != nil {
			fmt.Printf("failed storage.Client.Close() err=%s\n", err)
		}
	}()

	iter := gcsClient.Buckets(ctx, projectID)
	for {
		attrs, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", attrs.Name)
		pp.Println(attrs)
		v, ok := attrs.Labels["bucket"]
		if ok {
			// bucket というkeyのlabelがすでにある場合はskip
			fmt.Printf("already bucket label is %s\n", v)
			continue
		}
		var au storage.BucketAttrsToUpdate
		au.SetLabel("bucket", strings.ReplaceAll(attrs.Name, ".", "-"))
		updatedAttrs, err := gcsClient.Bucket(attrs.Name).If(storage.BucketConditions{MetagenerationMatch: attrs.MetaGeneration}).Update(ctx, au)
		if err != nil {
			panic(err)
		}
		fmt.Println("--UPDATED--")
		pp.Println(updatedAttrs)
		fmt.Println()
	}
}
