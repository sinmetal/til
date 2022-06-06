package main

import (
	"fmt"
	"net/http"
	"time"

	"go.mercari.io/datastore/aedatastore"
	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/", indexHandler)

	appengine.Main()
}

type Entity struct {
	CreatedAt time.Time
}

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.WithContext(r.Context(), r)

	ds, err := aedatastore.FromContext(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed aedatastore.FromContext. %s", err)))
		return
	}

	v := Entity{
		CreatedAt: time.Now(),
	}
	_, err = ds.Put(ctx, ds.NameKey("Sample1", time.Now().String(), nil), &v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed datastore.Put() %s", err)))
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}
