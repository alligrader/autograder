package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func main() {
	const port = ":80"
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(os.Stdout)

	r := mux.NewRouter()
	r.HandleFunc("/github", GitHubHandler)
	http.Handle("/", r)
	log.Printf("Serving on port %s", port)
	http.ListenAndServe(port, nil)
}
