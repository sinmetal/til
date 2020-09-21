package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/iam/v1"

	"github.com/sinmetal/til/storage/signedurl/storage"
)

func main() {
	ctx := context.Background()

	iamService, err := iam.NewService(ctx)
	if err != nil {
		panic(err)
	}
	signedURLService, err := storage.NewStorageSignedURLService(ctx, "signedurl@sinmetal-ci.iam.gserviceaccount.com", fmt.Sprintf("projects/%s/serviceAccounts/%s", "sinmetal-ci", "signedurl@sinmetal-ci.iam.gserviceaccount.com"), iamService)
	if err != nil {
		panic(err)
	}
	h := &Handlers{
		signedURLService: signedURLService,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	http.HandleFunc("/signedURL", h.HandleGetSignedURL)
	http.HandleFunc("/", StaticContentsHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), http.DefaultServeMux); err != nil {
		log.Printf("failed ListenAndServe err=%+v", err)
	}
}

type Handlers struct {
	signedURLService *storage.StorageSignedURLService
}

func (h *Handlers) HandleGetSignedURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contentType := r.FormValue("contentType")
	object := uuid.New().String()
	url, err := h.signedURLService.CreatePutObjectURL(ctx, "sinmetal-ci-signed-url", object, contentType, time.Now().Add(10*time.Minute))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := struct {
		URL string `json:"url"`
	}{
		URL: url,
	}
	e := json.NewEncoder(w)
	if err := e.Encode(res); err != nil {
		fmt.Println(err.Error())
	}
}
