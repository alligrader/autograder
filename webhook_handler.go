package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/RobbieMcKinstry/pipeline"
	"github.com/alligrader/jobs"
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const secretKey = "hello_alligrader"

func GitHubHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(secretKey))
	if err != nil {
		log.Fatal(err)
	}

	if len(payload) == 0 {
		log.Fatal("We have an empty payload problem!")
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Warn(err)
	}

	switch event := event.(type) {
	case *github.PushEvent:
		err = json.Unmarshal(payload, event)
		if err != nil {
			log.Info(string(payload))
			log.Fatal(err)
		}
		processPushEvent(w, event)
	case *github.PullRequestEvent:
		processPullRequestEvent(w, event)
	}
}

// processPushEvent needs to connect to kubernetes,
// set the metadata to have the right environment variables
// then launch those containers
func processPushEvent(w http.ResponseWriter, e *github.PushEvent) {
	const s = "Received a push event!"
	log.Info(s)
	defer w.Write([]byte(s))

	const (
		pipelineName = "checkstyleLinter"
		jarLoc       = "../jobs/lib/checkstyle-7.6.1-all.jar"
		srcDir       = ""
		checks       = "../jobs/.test/checkstyle.xml"
	)

	if e.Repo == nil {
		log.Info(e.String())
		log.Fatal("e.Repo is nil.")
	}

	if e.Repo.Owner == nil {
		log.Fatal("e.Repo.Owner is nil.")
	}

	if e.Repo.Owner.GetName() == "" {
		log.Fatal("e.Repo.Owner.GetName() is empty")
	}

	if e.GetAfter() == "" {
		log.Warn("e.GetAfter() is empty!")
	} else {
		log.Infof("e.GetAfter() == %v", e.GetAfter())
	}

	var (
		httpClient   = getClient()
		client       = github.NewClient(httpClient)
		repo         = e.Repo.GetName()
		owner        = e.Repo.Owner.GetName()
		ref          = e.GetAfter()
		outputLoc, _ = ioutil.TempFile("", "findbugs.out")
		fetchStep    = jobs.NewGithubStep(owner, repo, ref)
		checkStep    = jobs.NewCheckstyleStep(jarLoc, outputLoc.Name(), srcDir, checks, "", false)
		commentStep  = jobs.NewCommentStep(owner, repo, ref, client)
		pipe         = pipeline.New(pipelineName, 1000)
		stage        = pipeline.NewStage(pipelineName, false, false)
	)
	defer os.Remove(outputLoc.Name())

	stage.AddStep(fetchStep)
	stage.AddStep(checkStep)
	stage.AddStep(commentStep)
	pipe.AddStage(stage)

	res := pipe.Run()
	if res == nil {
		log.Fatal("No response!")
	}

	if res.Error != nil {
		log.Info("Returning an error in the stage result.")
		log.Fatal(res.Error)
	}
}

// processPushEvent needs to connect to kubernetes,
// set the metadata to have the right environment variables
// then launch those containers
func processPullRequestEvent(w http.ResponseWriter, e *github.PullRequestEvent) {
	const s = "Receieved a pull request event!"
	log.Info(s)
	w.Write([]byte(s))
}

// TODO really need to inject environment variables
func getClient() *http.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "8f23d9e3b9cc22d3be326928ee73c4880996de65"},
	)
	tc := oauth2.NewClient(ctx, ts)
	return tc
}
