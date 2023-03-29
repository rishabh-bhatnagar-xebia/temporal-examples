package main

import (
	"async/shared"
	"async/utils"
	"context"
	"go.temporal.io/sdk/client"
)

func TriggerWorkflow[workflowOutVar any](queueName string, workflow any, workerName string) (workflowOutVar, error) {
	utils.LogDebug("will be triggering a workflow on", queueName, "queue")
	c, err := client.Dial(client.Options{})
	if err != nil {
		utils.LogRed(err)
	}
	defer c.Close()

	workflowId := "worker-" + workerName
	utils.LogDebug(workflowId, "is the workflow id")
	input := shared.WorkflowIn{Data: "this is the input data to be persisted"}
	options := client.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: queueName,
	}

	run, err := c.ExecuteWorkflow(context.Background(), options, workflow, input)
	if err != nil {
		utils.LogDebug(err)
	}
	utils.LogGreen("Run ID:", run.GetRunID())

	var result workflowOutVar
	err = run.Get(context.Background(), &result)
	if err != nil {
		utils.LogRed("unable to get workflow result", err)
	}
	return result, err
}
