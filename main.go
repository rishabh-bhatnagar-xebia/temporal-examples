package main

import (
	"async/activities"
	"async/shared"
	"async/utils"
	custom_worker "async/worker"
	"async/workflows"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		utils.LogRed("worker type and queue type must be specified, exitting")
		os.Exit(1)
	}
	queueType := os.Args[2]
	workerName := os.Args[1]
	var workflow any

	var queueName string
	switch queueType {
	case "basic":
		queueName = shared.QueueNameBasic
		workflow = workflows.Basic
	case "async_v1":
		queueName = shared.QueueNameAsyncV1
		workflow = workflows.AsyncWithChild
	default:
		utils.LogRed(fmt.Sprintf("unknown queue type %s", queueType))
		os.Exit(1)
	}

	switch strings.ToLower(workerName) {
	case "db":
		custom_worker.SpawnActivityWorker(queueName, activities.WriteToDB, workflow)
	case "git":
		custom_worker.SpawnActivityWorker(queueName, activities.WriteToGit, workflow)
	default:
		utils.LogRed(fmt.Sprintf("unknown activity/workflow %s", workerName))
		os.Exit(1)
	}
}
