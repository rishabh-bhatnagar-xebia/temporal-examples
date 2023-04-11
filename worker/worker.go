package worker

import (
	"learn_temporal/utils"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func SpawnActivityWorker(queueName string, workflows []any, activity any) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		utils.LogRed(err)
		return
	}
	defer c.Close()

	w := worker.New(c, queueName, worker.Options{})
	w.RegisterActivity(activity)
	utils.LogDebug("serving activity:", activity)
	for _, wf := range workflows {
		utils.LogDebug("associated workflow:", wf)
		w.RegisterWorkflow(wf)
	}

	err = w.Run(worker.InterruptCh())
	if err != nil {
		utils.LogRed("error running a worker:", err)
	}
}
