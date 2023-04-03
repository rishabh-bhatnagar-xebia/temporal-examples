package workflows

import (
	"async/activities"
	"async/shared"
	"async/utils"
	"time"

	"go.temporal.io/sdk/workflow"
)

func AsyncWithQueries(ctx workflow.Context, data shared.WorkflowIn) (shared.WorkflowAsyncV2Out, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 250,
	})

	var state shared.WorkflowAsyncV2Status
	err := workflow.SetQueryHandler(ctx, "current_state", func() (shared.WorkflowAsyncV2Status, error) {
		return state, nil
	})
	if err != nil {
		utils.LogRed(err)
		return shared.WorkflowAsyncV2Out{}, nil
	}

	state.Status = "inside the basic workflow"
	utils.LogGreen("WORKFLOW;", state)

	// Writing to DB
	state.Status = "calling the db activity"
	utils.LogDebug("WORKFLOW;", state)
	var dbOut shared.DBOut
	err = workflow.ExecuteActivity(ctx, activities.WriteToDB, data, time.Second*10).Get(ctx, &dbOut)
	if err != nil {
		utils.LogRed("WORKFLOW;", err)
	}
	state.Status = "db activity completed executing"
	state.Result.DBOut = &dbOut
	utils.LogDebug("WORKFLOW;", "db out", state.Result.DBOut)

	state.Status = "calling the git activity"
	utils.LogDebug("WORKFLOW;", state)
	// Storing To Git
	var gitOut shared.GitOut
	err = workflow.ExecuteActivity(ctx, activities.WriteToGit, data, dbOut, time.Second*10).Get(ctx, &gitOut)
	if err != nil {
		utils.LogRed("WORKFLOW;", err)
	}
	state.Status = "completed running the git-activity"
	utils.LogDebug("WORKFLOW;", "git out", gitOut)

	utils.LogGreen("WORKFLOW;", "completed the workflow")
	state.Completed = true
	state.Result = shared.WorkflowAsyncV2Out{
		DBOut:  &dbOut,
		GitOut: &gitOut,
	}
	return state.Result, err
}
