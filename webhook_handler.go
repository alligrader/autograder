package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
)

const secretKey = "hello_alligrader"

func GitHubHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(secretKey))
	if err != nil {
		log.Warn(err)
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Warn(err)
	}

	switch event := event.(type) {
	case *github.PushEvent:
		processPushEvent(w, event)
	case *github.PullRequestEvent:
		processPullRequestEvent(w, event)
	}
}

func processPushEvent(w http.ResponseWriter, e *github.PushEvent) {
	const s = "Received a push event!"
	log.Info(s)
	w.Write([]byte(s))
}

func processPullRequestEvent(w http.ResponseWriter, e *github.PullRequestEvent) {
	const s = "Receieved a pull request event!"
	log.Info(s)
	w.Write([]byte(s))
}
