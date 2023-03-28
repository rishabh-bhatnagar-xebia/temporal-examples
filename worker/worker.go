package worker

import (
	"async/utils"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func SpawnActivityWorker(queueName string, activity any, workflow any) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		utils.LogRed(err)
		return
	}
	defer c.Close()

	w := worker.New(c, queueName, worker.Options{})
	w.RegisterActivity(activity)
	w.RegisterWorkflow(workflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		utils.LogRed("error running a worker:", err)
	}
}

func SpawnWorkflowWorker(queueName string, workflow interface{}) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		utils.LogRed(err)
	}
	defer c.Close()

	w := worker.New(c, queueName, worker.Options{})
	w.RegisterWorkflow(workflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		utils.LogRed("error running a worker:", err)
	}
}
