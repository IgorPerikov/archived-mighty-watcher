package data

import (
	"github.com/google/go-github/v26/github"
	"strings"
)

type RepoName struct {
	Owner string
	Name  string
}

func (r RepoName) String() string {
	return r.Owner + "/" + r.Name
}

func GetRepoNameFromRepo(repo *github.Repository) RepoName {
	segments := strings.Split(repo.GetFullName(), "/")
	return RepoName{Owner: segments[0], Name: segments[1]}
}

func GetRepoNameFromIssue(issue *github.Issue) RepoName {
	urlSegments := strings.Split(issue.GetRepositoryURL(), "/")
	length := len(urlSegments)
	return RepoName{Owner: urlSegments[length-2], Name: urlSegments[length-1]}
}
