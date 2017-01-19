package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(os.Stdout)

	r := mux.NewRouter()
	r.HandleFunc("/github", GitHubHandler)
	http.Handle("/", r)
}
