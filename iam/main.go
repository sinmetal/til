package iam

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/iam/credentials/apiv1"
	credentialspb "google.golang.org/genproto/googleapis/iam/credentials/v1"
)

func GetSignedJwt(ctx context.Context, sa string) {
	c, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		panic(err)
	}

	d := time.Now().Add(3000 * time.Second).Unix()
	name := fmt.Sprintf("projects/-/serviceAccounts/%s", sa)
	payload := fmt.Sprintf("{ \"iss\": \"%s\", \"sub\": \"%s\", \"aud\": \"https://firestore.googleapis.com/google.firestore.v1beta1.Firestore\", \"iat\": %d, \"exp\": %d }", sa, sa, d, d)
	fmt.Println(name)
	fmt.Println(payload)
	resp, err := c.SignJwt(ctx, &credentialspb.SignJwtRequest{
		Name:      name,
		Delegates: []string{},
		Payload:   payload,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.SignedJwt)
}
