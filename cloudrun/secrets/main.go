package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello World")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func ReadSecretEnvHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := os.Getenv(key)
	_, err := fmt.Fprintf(w, "%s:%s", key, value)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func ReadSecretFileHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	fileBody, err := os.ReadFile(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err2 := fmt.Fprintf(w, "failed read file %s, err=%s", key, err)
		if err2 != nil {
			fmt.Println(err2.Error())
			return
		}
		return
	}

	w.Header().Set("content-type", "application/json;charset=utf8")
	m := map[string]string{
		"key":  key,
		"body": string(fileBody),
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Start Listen PORT:%s\n", port)

	http.HandleFunc("/secrets/env", ReadSecretEnvHandler)
	http.HandleFunc("/secrets/file", ReadSecretFileHandler)
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		fmt.Println(err.Error())
	}
}
