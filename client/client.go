package client

import (
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

func New() *http.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)

		return tc
	}
	return http.DefaultClient
}
