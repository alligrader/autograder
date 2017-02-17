package main

import (
	"net/http"
)

func GitHubHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Success!`))
}
