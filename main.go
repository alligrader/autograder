package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	githubHandler := &githubHandler{logger, secretKey, githubToken}

	r := mux.NewRouter()
	r.Handle("/github", githubHandler)
	http.Handle("/", r)
	logger.Printf("Serving on port %s", port)
	logger.Fatal(http.ListenAndServe(port, nil))
}
