package main

import (
	"learn_temporal/activities"
	"learn_temporal/shared"
	"learn_temporal/utils"
	"learn_temporal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func SpawnActivityWorker(queueName string) {
	utils.LogGreen("spawning worker on", queueName)
	c, err := client.Dial(client.Options{})
	if err != nil {
		utils.LogRed(err)
		return
	}
	defer c.Close()

	w := worker.New(c, queueName, worker.Options{})
	w.RegisterWorkflow(workflows.GitWorkflow)
	w.RegisterWorkflow(workflows.AsyncWithChild)
	w.RegisterActivity(activities.WriteToGit)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		utils.LogRed("error running a worker:", err)
	}
}

func main() {
	SpawnActivityWorker(shared.QueueNameAsyncV1)
}
