package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	fmt.Println("Starting app...")

	var githubToken = conf.GetString("GH_ACCESS_TOKEN")
	var secretKey = conf.GetString("GH_SECRET_KEY")
	var port = conf.GetString("PORT")

	githubHandler := &githubHandler{
		log:         logger,
		secretKey:   secretKey,
		accessToken: githubToken,
		conf:        conf,
	}

	fmt.Println("Created github handler")

	r := mux.NewRouter()
	r.Handle("/github", githubHandler)
	http.Handle("/", r)
	logger.Infof("Serving on port %s", port)
	logger.Fatal(http.ListenAndServe(port, nil))
}
