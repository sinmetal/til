package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/compute/metadata"
)

func ToRunHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := makeGetRequest("https://gcpboxtest-73zry4yfvq-an.a.run.app")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	_, err = w.Write(b)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// makeGetRequest makes a GET request to the specified Cloud Run endpoint in
// serviceURL (must be a complete URL) by authenticating with the ID token
// obtained from the Metadata API.
func makeGetRequest(serviceURL string) (*http.Response, error) {
	// query the id_token with ?audience as the serviceURL
	tokenURL := fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", serviceURL)
	idToken, err := metadata.Get(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("metadata.Get: failed to query id_token: %+v", err)
	}
	fmt.Printf("idToken:%s\n", idToken)

	req, err := http.NewRequest("GET", serviceURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", idToken))
	return http.DefaultClient.Do(req)
}
