package main

import (
	"async/shared"
	"async/utils"
	"async/workflows"
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"os"
)

func TriggerWorkflow[workflowOutVar any](queueName string, workflow any) workflowOutVar {
	utils.LogDebug("will be triggering a workflow on", queueName, "queue")
	c, err := client.Dial(client.Options{})
	if err != nil {
		utils.LogRed(err)
		os.Exit(1)
	}
	defer c.Close()

	input := shared.WorkflowIn{Data: "this is the input data to be persisted"}
	options := client.StartWorkflowOptions{
		ID:        "random-id",
		TaskQueue: queueName,
	}

	run, err := c.ExecuteWorkflow(context.Background(), options, workflow, input)
	if err != nil {
		utils.LogDebug(err)
	}

	var result workflowOutVar
	err = run.Get(context.Background(), &result)
	if err != nil {
		utils.LogRed("unable to get workflow result", err)
		return result
	}
	return result
}

func main() {
	res := TriggerWorkflow[shared.WorkflowBasicOut](shared.QueueNameBasic, workflows.Basic)
	fmt.Println(res.DBOut.ID)
	fmt.Println(res.GitOut.ID)
}
