package formatter

import (
	"github.com/google/go-github/v26/github"
	"sort"
	"time"
)

func Format(issues []*github.Issue) []*github.Issue {
	return sortByCreatedAtDesc(removeOld(uniquely(issues)))
}

func uniquely(issues []*github.Issue) []*github.Issue {
	issuesById := make(map[int64]*github.Issue)
	for _, issue := range issues {
		issuesById[issue.GetID()] = issue
	}
	values := make([]*github.Issue, 0, len(issuesById))
	for _, value := range issuesById {
		values = append(values, value)
	}
	return values
}

func removeOld(issues []*github.Issue) []*github.Issue {
	twoYearsAgo := time.Now().AddDate(-2, 0, 0)
	freshIssues := make([]*github.Issue, 0, len(issues))
	for _, issue := range issues {
		if issue.GetCreatedAt().After(twoYearsAgo) {
			freshIssues = append(freshIssues, issue)
		}
	}
	return freshIssues
}

func sortByCreatedAtDesc(issues []*github.Issue) []*github.Issue {
	sort.Slice(issues, func(i, j int) bool {
		// TODO: sort by UpdatedAt?
		return issues[i].GetCreatedAt().After(issues[j].GetCreatedAt())
	})
	return issues
}
