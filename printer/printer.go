package printer

import (
	"fmt"
	"github.com/IgorPerikov/mighty-watcher/data"
	"github.com/google/go-github/v26/github"
	"time"
)

func Print(issues []*github.Issue) {
	for _, str := range convertToTimeGroups(issues) {
		fmt.Printf("%s\n", str)
	}
}

// TODO: make it better and more flexible, support today, yesterday, this week, this month, everything else
func convertToTimeGroups(issues []*github.Issue) []string {
	res := make([]string, 0, len(issues)+3)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastWeek := today.AddDate(0, 0, -7)

	res = append(res, "today:")
	isToday := true
	isLastWeek := false

	for _, issue := range issues {
		if isToday && issue.GetCreatedAt().Before(today) {
			isLastWeek = true
			isToday = false
			res = append(res, "\nlast week:")
		}
		if isLastWeek && issue.GetCreatedAt().Before(lastWeek) {
			isLastWeek = false
			res = append(res, "\nolder:")
		}
		resultString := fmt.Sprintf(
			"%v \"%v\" %v",
			data.GetRepoNameFromIssue(issue),
			issue.GetTitle(),
			issue.GetHTMLURL(),
		)
		res = append(res, resultString)
	}
	return res
}
