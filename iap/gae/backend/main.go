package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	iapbox "github.com/sinmetalcraft/gcpbox/iap/appengine"
)

func main() {
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	_, user := iapbox.CurrentUserWithContext(r.Context(), r)
	if user != nil {
		fmt.Fprintf(w, "Hello, %s", user.Email)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}
