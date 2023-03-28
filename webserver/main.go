package main

import (
	"async/shared"
	"async/utils"
	"async/workflows"
	"context"
	"encoding/json"
	"fmt"
	"go.temporal.io/sdk/client"
	"net/http"
)

func TriggerWorkflow[workflowOutVar any](queueName string, workflow any) (workflowOutVar, error) {
	utils.LogDebug("will be triggering a workflow on", queueName, "queue")
	c, err := client.Dial(client.Options{})
	if err != nil {
		utils.LogRed(err)
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
	}
	return result, err
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/workflow", func(w http.ResponseWriter, r *http.Request) {
		pathQuery := r.URL.Query()
		if !pathQuery.Has(shared.HttpWorkflowTypeParamName) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("missing workflow type"))
			return
		}
		workflowType := pathQuery.Get(shared.HttpWorkflowTypeParamName)
		switch workflowType {
		case "basic":
			result, err := TriggerWorkflow[shared.WorkflowBasicOut](shared.QueueNameBasic, workflows.Basic)
			utils.LogGreen(fmt.Sprintf("DB Output:%+v", *result.DBOut))
			utils.LogGreen(fmt.Sprintf("Git Output:%+v", *result.GitOut))
			if err != nil {
				utils.LogRed(err)
			}
			out, _ := json.Marshal(result)
			_, _ = w.Write(out)
		case "async_v1":
			result, err := TriggerWorkflow[shared.WorkflowAsyncV1Out](shared.QueueNameAsyncV1, workflows.AsyncWithChild)
			if err != nil {
				utils.LogRed(err)
			}
			utils.LogGreen(fmt.Sprintf("DB Output:%+v", *result.DBOut))
			out, _ := json.Marshal(result)
			_, _ = w.Write(out)
		default:
			http.Error(w, "invalid workflow type "+workflowType, http.StatusBadRequest)
		}
	})
	port := "8001"
	fmt.Println("serving on port:", port)
	utils.LogRed(http.ListenAndServe(":"+port, router))
}
