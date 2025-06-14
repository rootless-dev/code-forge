package main

import "github.com/carlosealves2/code-forge/oidc-service/internal/bootstrap"

var (
	Version = "0.0.1"
	Commit  = "short-commit"
)

func main() {
	bootstrap.New(Version, Commit).Run()
}
