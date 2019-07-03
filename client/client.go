package client

import (
	"context"
	"fmt"
	"github.com/IgorPerikov/mighty-watcher/data"
	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
	"log"
	"os"
)

var client *github.Client

const tokenEnvName = "MIGHTY_WATCHER_GITHUB_TOKEN"

func init() {
	token, exists := os.LookupEnv(tokenEnvName)
	if !exists {
		log.Panicf("Environment variable %s not found", tokenEnvName)
	}
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client = github.NewClient(oauth2.NewClient(context.TODO(), tokenSource))
}

func GetStarred(ctx context.Context, perPage int) []data.RepoName {
	starredRepositories, _, err := client.Activity.ListStarred(
		ctx,
		"",
		&github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				PerPage: perPage,
			},
		},
	)
	if err != nil {
		log.Panicf("Getting starred repos failed with %v", err)
	}
	var repoFullNames []data.RepoName
	for _, starredRepo := range starredRepositories {
		repoFullNames = append(repoFullNames, data.GetRepoNameFromRepo(starredRepo.Repository))
	}
	return repoFullNames
}

func GetIssues(ctx context.Context, owner string, name string, label string) ([]*github.Issue, error) {
	issues, _, err := client.Issues.ListByRepo(
		ctx,
		owner,
		name,
		&github.IssueListByRepoOptions{
			State:    "open",
			Assignee: "none",
			Labels:   []string{label},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get issues by label, reason=%v", err)
	}
	return issues, nil
}
