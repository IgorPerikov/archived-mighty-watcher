package main // TODO: names, packages everywhere

import (
	"context"
	"github.com/IgorPerikov/mighty-watcher/client"
	"github.com/IgorPerikov/mighty-watcher/formatter"
	"github.com/IgorPerikov/mighty-watcher/labels"
	"github.com/IgorPerikov/mighty-watcher/printer"
	"github.com/google/go-github/v26/github"
	"golang.org/x/sync/semaphore"
	"log"
	"sync"
)

const parallelismLevel = 40

func main() {
	var result []*github.Issue
	var wg sync.WaitGroup
	var mutex sync.Mutex

	weighted := semaphore.NewWeighted(parallelismLevel) // TODO: name

	for _, task := range createTasks() {
		wg.Add(1)
		go func(task *searchTask) {
			defer wg.Done()
			err := weighted.Acquire(context.TODO(), 1)
			if err == nil {
				defer weighted.Release(1)
				issues := task.execute()
				mutex.Lock()
				defer mutex.Unlock()
				result = append(result, issues...)
			} else {
				log.Printf("Error acquiring semaphore %v", err)
			}
		}(task)
	}
	wg.Wait()

	formattedIssues := formatter.Format(result)
	printer.Print(formattedIssues)
}

func createTasks() []*searchTask {
	var tasks []*searchTask

	repoNames := client.GetStarred(context.TODO(), 2000)
	easyLabels := labels.GetEasy()

	for _, repoName := range repoNames {
		for _, label := range easyLabels {
			tasks = append(tasks, &searchTask{repoName.Owner, repoName.Name, label})
		}
	}
	return tasks
}

type searchTask struct {
	owner string
	name  string
	label string
}

func (task *searchTask) execute() []*github.Issue {
	issues, err := client.GetIssues(context.TODO(), task.owner, task.name, task.label)
	if err != nil {
		log.Println(err)
	}
	return issues
}
