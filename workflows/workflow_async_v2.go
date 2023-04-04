package workflows

import (
	"async/activities"
	workflowtype "async/protoc_types"
	"async/utils"
	"time"

	"go.temporal.io/sdk/workflow"
)

func AsyncWithQueries(ctx workflow.Context, data workflowtype.WorkflowIn) (workflowtype.WorkflowAsyncV2Out, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 250,
	})

	state := newState()
	err := workflow.SetQueryHandler(ctx, "current_state", func() (workflowtype.WorkflowAsyncV2Status, error) {
		return state, nil
	})
	if err != nil {
		utils.LogRed(err)
		return workflowtype.WorkflowAsyncV2Out{}, nil
	}

	*state.Status = "inside the basic workflow"
	utils.LogGreen("WORKFLOW;", state)

	// Writing to DB
	*state.Status = "calling the db activity"
	utils.LogDebug("WORKFLOW;", state)
	var dbOut workflowtype.DBOut
	err = workflow.ExecuteActivity(ctx, activities.WriteToDB, data, time.Second*10).Get(ctx, &dbOut)
	if err != nil {
		utils.LogRed("WORKFLOW;", err)
	}
	*state.Status = "db activity completed executing"
	state.Result.DBOut = &dbOut
	utils.LogDebug("WORKFLOW;", "db out", state.Result.DBOut)

	*state.Status = "calling the git activity"
	utils.LogDebug("WORKFLOW;", state)
	// Storing To Git
	var gitOut workflowtype.GitOut
	err = workflow.ExecuteActivity(ctx, activities.WriteToGit, data, dbOut, time.Second*10).Get(ctx, &gitOut)
	if err != nil {
		utils.LogRed("WORKFLOW;", err)
	}
	*state.Status = "completed running the git-activity"
	utils.LogDebug("WORKFLOW;", "git out", gitOut)

	utils.LogGreen("WORKFLOW;", "completed the workflow")
	*state.Completed = true
	state.Result = &workflowtype.WorkflowAsyncV2Out{
		DBOut: &dbOut,
	}
	return *state.Result, err
}

func newState() workflowtype.WorkflowAsyncV2Status {
	var status = "starting workflow"
	var completed = false
	var state = workflowtype.WorkflowAsyncV2Status{
		Status:    &status,
		Completed: &completed,
		Result: &workflowtype.WorkflowAsyncV2Out{
			DBOut: nil,
		},
	}
	return state
}
