package workflows

import (
	"async/activities"
	"async/shared"
	"async/utils"
	"time"

	"go.temporal.io/sdk/workflow"
)

func AsyncWithQueries(ctx workflow.Context, data shared.WorkflowIn) (shared.WorkflowAsyncV2Out, error) {
	take2 := func() { time.Sleep(time.Second * 2) }
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: time.Second * 250})
	status := "inside the basic workflow"
	take2()
	utils.LogGreen("WORKFLOW;", status)

	err := workflow.SetQueryHandler(ctx, "current_state", func() (string, error) {
		return status, nil
	})
	if err != nil {
		utils.LogRed(err)
		status = "failed to set a query handler: " + err.Error()
		take2()
	}

	// Writing to DB
	status = "calling the db activity"
	take2()
	utils.LogDebug("WORKFLOW;", status)
	var dbOut shared.DBOut
	err = workflow.ExecuteActivity(ctx, activities.WriteToDB, data).Get(ctx, &dbOut)
	if err != nil {
		utils.LogRed("WORKFLOW;", err)
	}
	status = "db activity completed executing"
	take2()
	utils.LogDebug("WORKFLOW;", "db out", dbOut)

	status = "calling the git activity"
	take2()
	utils.LogDebug("WORKFLOW;", status)
	// Storing To Git
	var gitOut shared.GitOut
	err = workflow.ExecuteActivity(ctx, activities.WriteToGit, data, dbOut, time.Second*5).Get(ctx, &gitOut)
	if err != nil {
		utils.LogRed("WORKFLOW;", err)
	}
	status = "completed running the git-activity"
	take2()
	utils.LogDebug("WORKFLOW;", "git out", gitOut)

	utils.LogGreen("WORKFLOW;", "completed the workflow")
	return shared.WorkflowAsyncV2Out{
		DBOut: &dbOut, GitOut: &gitOut,
	}, err
}
