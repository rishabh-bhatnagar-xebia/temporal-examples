package main

import (
	"fmt"
	"learn_temporal/activities"
	"learn_temporal/shared"
	"learn_temporal/utils"
	custom_worker "learn_temporal/worker"
	"learn_temporal/workflows"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		utils.LogRed("worker type and queue type must be specified, exitting")
		os.Exit(1)
	}

	queueType := os.Args[2]
	var wfs []any
	var queueName string
	switch queueType {
	case "basic":
		queueName = shared.QueueNameBasic
		wfs = []any{workflows.Basic}
	case "async_v1":
		queueName = shared.QueueNameAsyncV1
		wfs = []any{workflows.AsyncWithChild, workflows.GitWorkflow}
	case "async_v2":
		queueName = shared.QueueNameAsyncV2
		wfs = []any{workflows.AsyncWithQueries}
	default:
		utils.LogRed(fmt.Sprintf("unknown queue type %s", queueType))
		os.Exit(1)
	}

	workerName := os.Args[1]
	var activity any
	switch strings.ToLower(workerName) {
	case "db":
		activity = activities.WriteToDB
	case "git":
		activity = activities.WriteToGit
	default:
		utils.LogRed(fmt.Sprintf("unknown activity/workflow %s", workerName))
		os.Exit(1)
	}

	custom_worker.SpawnActivityWorker(queueName, wfs, activity)
}
