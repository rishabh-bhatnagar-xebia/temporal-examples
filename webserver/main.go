package main

import (
	"context"
	"encoding/json"
	"fmt"
	"learn_temporal/shared"
	"learn_temporal/utils"
	"learn_temporal/workflows"
	workflowtype "learn_temporal/workflowtype"
	"net/http"
	"os"
	"time"

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
	data := "this is the input data to be persisted"
	input := workflowtype.WorkflowIn{Data: &data}
	options := client.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: queueName,
	}

	// fmt.Println("before shared.C *****************************")
	// fmt.Println("shared.C", shared.C.Get(), shared.C.Get())
	run, err := c.ExecuteWorkflow(context.Background(), options, workflow, input)
	if err != nil {
		utils.LogDebug(err)
		os.Exit(0)
	}
	utils.LogGreen("Run ID:", run.GetRunID())

	var result workflowOutVar
	err = run.Get(context.Background(), &result)
	if err != nil {
		utils.LogRed("unable to get workflow result", err)
	}
	return result, err
}

func TriggerWorkflowAsync2(queueName string, workflow any, workerName string, data string) (workflowtype.WorkflowAsyncV2Out, error) {
	utils.LogDebug("will be triggering a workflow on", queueName, "queue")
	c, err := client.Dial(client.Options{})
	if err != nil {
		utils.LogRed(err)
	}
	defer c.Close()

	workflowId := "worker-" + workerName
	utils.LogDebug(workflowId, "is the workflow id")
	input := workflowtype.WorkflowIn{Data: &data}
	options := client.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: queueName,
	}

	run, err := c.ExecuteWorkflow(context.Background(), options, workflow, input)
	if err != nil {
		utils.LogDebug(err)
	}
	utils.LogGreen("Run ID:", run.GetRunID())

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	var result workflowtype.WorkflowAsyncV2Out
	for range ticker.C {
		query := "current_state"
		resp, err := c.QueryWorkflow(context.Background(), workflowId, run.GetRunID(), query)
		if err != nil {
			utils.LogRed(err)
			continue
		}
		if !resp.HasValue() {
			continue
		}
		var state workflowtype.WorkflowAsyncV2Status
		if err := resp.Get(&state); err != nil {
			utils.LogRed(err)
			continue
		}
		printResult(state)
		if *state.Completed || (state.Result != nil && state.Result.DBOut != nil) {
			result = *state.Result
			break
		}
	}
	return result, err
}

func printResult(state workflowtype.WorkflowAsyncV2Status) {
	out := ""
	dbOut := &out
	if state.Result.DBOut != nil {
		dbOut = state.Result.DBOut.ID
	}
	gitOut := &out
	utils.LogGreen("current status:", fmt.Sprintf("{DBID: %+v; GitID: %+v}", *dbOut, *gitOut))
}

func main() {
	data := "this is the input data to be persisted"
	router := http.NewServeMux()
	router.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if !q.Has("run_id") {
			http.Error(w, "expected run_id as a param in the url", http.StatusBadRequest)
			return
		}
		runId := q.Get("run_id")

		c, err := client.Dial(client.Options{})
		if err != nil {
			utils.LogRed(err)
		}
		defer c.Close()

		resp, err := c.QueryWorkflowWithOptions(context.Background(), &client.QueryWorkflowWithOptionsRequest{
			WorkflowID: "worker-async_v2",
			RunID:      runId,
			QueryType:  "current_state",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		utils.LogGreen(resp, "is the resp")
	})
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
			utils.LogDebug("triggering a Basic workflow")
			result, err := TriggerWorkflow[workflowtype.WorkflowBasicOut](shared.QueueNameBasic, workflows.Basic, workflowType)
			utils.LogGreen(fmt.Sprintf("DB Output:%+v", *result.DBOut))
			utils.LogGreen(fmt.Sprintf("Git Output:%+v", *result.GitOut))
			if err != nil {
				utils.LogRed(err)
			}
			out, _ := json.Marshal(result)
			_, _ = w.Write(out)
		case "async_v1":
			utils.LogDebug("triggering a AsyncWithChild workflow")
			result, err := TriggerWorkflow[workflowtype.WorkflowAsyncV1Out](shared.QueueNameAsyncV1, workflows.AsyncWithChild, data)
			if err != nil {
				utils.LogRed(err)
			}
			utils.LogGreen(fmt.Sprintf("DB Output:%+v", *result.DBOut))
			out, _ := json.Marshal(result)
			_, _ = w.Write(out)
		case "async_v2":
			utils.LogDebug("triggering a AsyncWithQueries workflow")
			result, err := TriggerWorkflowAsync2(shared.QueueNameAsyncV2, workflows.AsyncWithQueries, workflowType, data)
			if err != nil {
				utils.LogRed(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
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
