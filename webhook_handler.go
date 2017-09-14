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
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type environment int

const (
	Production environment = iota
	Test
	Staging
	Development

	ProductionEnv  = "production"
	TestEnv        = "test"
	StagingEnv     = "staging"
	DevelopmentEnv = "development"

	configFileName = "configuration"

	environmentKey = "ENV"
)

var (
	conf   = viper.New()
	env    environment
	logger *logrus.Logger
)

func init() {

	conf.SetConfigName(configFileName)
	conf.AddConfigPath(".")
	conf.ReadInConfig()
	conf.AutomaticEnv()
	conf.WatchConfig()
	setEnvironment(conf)

	logger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.TextFormatter{ForceColors: true},
		Level:     logrus.InfoLevel,
	}

	switch env {
	case Production:
		logger.Formatter = &logrus.JSONFormatter{}
		logger.Level = logrus.WarnLevel
	case Development:
		logger.Level = logrus.DebugLevel
	}
}

func setEnvironment(config *viper.Viper) {

	enumMapper := map[string]environment{
		ProductionEnv:  Production,
		TestEnv:        Test,
		DevelopmentEnv: Development,
		StagingEnv:     Staging,
	}

	env = enumMapper[config.GetString(environmentKey)]
}

type githubHandler struct {
	log         *logrus.Logger
	secretKey   string
	accessToken string
	conf        *viper.Viper
}

func (g *githubHandler) getClient() *http.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return tc
}

func (g *githubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var (
		payload = g.getPayload(r, g.secretKey)
		event   = g.getEvent(payload, r)
	)

	g.handleEvent(w, event, payload)
}

func (g *githubHandler) getPayload(r *http.Request, secretKey string) []byte {
	payload, err := github.ValidatePayload(r, []byte(g.secretKey))
	if err != nil {
		g.log.Fatal(err)
	}
	return payload
}

func (g *githubHandler) getEvent(payload []byte, r *http.Request) interface{} {
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		g.log.Warn(err)
	}
	return event
}

func (g *githubHandler) handleEvent(w http.ResponseWriter, event interface{}, payload []byte) {
	switch event := event.(type) {
	case *github.PushEvent:
		err := json.Unmarshal(payload, event)
		if err != nil {
			g.log.Fatal(err)
		}
		g.processPushEvent(w, event)
	case *github.PullRequestEvent:
		g.processPullRequestEvent(w, event)
	}
}

// TODO
// processPushEvent needs to connect to kubernetes,
// set the metadata to have the right environment variables
// then launch those containers
func (g *githubHandler) processPushEvent(w http.ResponseWriter, e *github.PushEvent) {
	const s = "Received a push event!"
	g.log.Info(s)
	defer w.Write([]byte(s))

	const (
		pipelineName = "checkstyleLinter"
		jarLoc       = "../jobs/lib/checkstyle-7.6.1-all.jar"
		srcDir       = ""
		checks       = "../jobs/.test/checkstyle.xml"
	)

	if e.Repo == nil {
		g.log.Info(e.String())
		g.log.Fatal("e.Repo is nil.")
	}

	if e.Repo.Owner == nil {
		g.log.Fatal("e.Repo.Owner is nil.")
	}

	if e.Repo.Owner.GetName() == "" {
		g.log.Fatal("e.Repo.Owner.GetName() is empty")
	}

	if e.GetAfter() == "" {
		g.log.Warn("e.GetAfter() is empty!")
	} else {
		g.log.Infof("e.GetAfter() == %v", e.GetAfter())
	}

	var (
		httpClient   = g.getClient()
		client       = github.NewClient(httpClient)
		repo         = e.Repo.GetName()
		owner        = e.Repo.Owner.GetName()
		ref          = e.GetAfter()
		outputLoc, _ = ioutil.TempFile("", "findbugs.out")
		fetchStep    = jobs.NewGithubStep(owner, repo, ref, g.log)
		checkStep    = jobs.NewCheckstyleStep(jarLoc, outputLoc.Name(), srcDir, checks, false, g.log)
		commentStep  = jobs.NewCommentStep(owner, repo, ref, client, g.log)
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
		g.log.Fatal("No response!")
	}

	if res.Error != nil {
		g.log.Info("Returning an error in the stage result.")
		g.log.Fatal(res.Error)
	}
}

// TODO
// processPushEvent needs to connect to kubernetes,
// set the metadata to have the right environment variables
// then launch those containers
func (g *githubHandler) processPullRequestEvent(w http.ResponseWriter, e *github.PullRequestEvent) {
	w.Write([]byte("Yay!"))
}
