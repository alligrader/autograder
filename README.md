# autograder

This service does not have a UI. It just listens for Git hooks to come in and launches the containers. It speaks directly with Kubernetes.

[Here is a tutorial](https://groob.io/tutorial/go-github-webhook/) that describes how to deal with the incoming git hooks. This tutorial is like 50% of all the code this repo expects to have.

#Summary
Listens for Git hooks and talks to kubernetes

# TODO

I need to make the Rakefile pull the version from the environment variables, then break the tasks into smaller chunks.
