package main

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

func Hoge() {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	// Project ID for this request.
	project := "sinmetal-grpc-plaground" // TODO: Update placeholder value.

	// Name of the SslCertificate resource to return.
	sslCertificate := "greeting-sinmetalcraft-jp" // TODO: Update placeholder value.

	resp, err := computeService.SslCertificates.Get(project, sslCertificate).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Change code below to process the `resp` object:
	fmt.Printf("%v\n", resp.PrivateKey)
}
