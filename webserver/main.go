package main

import (
	"async/shared"
	"async/utils"
	"async/workflows"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.temporal.io/sdk/client"
)

func main() {
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
			result, err := TriggerWorkflow[shared.WorkflowBasicOut](shared.QueueNameBasic, workflows.Basic, workflowType)
			utils.LogGreen(fmt.Sprintf("DB Output:%+v", *result.DBOut))
			utils.LogGreen(fmt.Sprintf("Git Output:%+v", *result.GitOut))
			if err != nil {
				utils.LogRed(err)
			}
			out, _ := json.Marshal(result)
			_, _ = w.Write(out)
		case "async_v1":
			utils.LogDebug("triggering a AsyncWithChild workflow")
			result, err := TriggerWorkflow[shared.WorkflowAsyncV1Out](shared.QueueNameAsyncV1, workflows.AsyncWithChild, workflowType)
			if err != nil {
				utils.LogRed(err)
			}
			utils.LogGreen(fmt.Sprintf("DB Output:%+v", *result.DBOut))
			out, _ := json.Marshal(result)
			_, _ = w.Write(out)
		case "async_v2":
			utils.LogDebug("triggering a AsyncWithQueries workflow")
			result, err := TriggerWorkflow[shared.WorkflowAsyncV1Out](shared.QueueNameAsyncV2, workflows.AsyncWithQueries, workflowType)
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
