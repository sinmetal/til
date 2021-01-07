package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Printf("Defaulting to port %s", port)
	}

	fmt.Printf("Listening on port %s", port)
	http.HandleFunc("/run/idtokenreq", ToRunHandler)
	http.HandleFunc("/", HelloHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), http.DefaultServeMux); err != nil {
		fmt.Printf("failed ListenAndServe err=%+v", err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("request.header:%+v", r.Header)

	_, err := w.Write([]byte(fmt.Sprintf("Hello GAEClientInvoker : %s", time.Now().String())))
	if err != nil {
		fmt.Println(err.Error())
	}
}
