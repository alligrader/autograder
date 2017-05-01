package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	const port = ":80"
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	r := mux.NewRouter()
	r.HandleFunc("/github", GitHubHandler)
	http.Handle("/", r)
	log.Printf("Serving on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
