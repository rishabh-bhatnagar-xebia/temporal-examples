package workflows

import (
	"async/activities"
	"async/shared"
	"async/utils"
	"go.temporal.io/sdk/workflow"
	"time"
)

func Basic(ctx workflow.Context, data shared.WorkflowIn) (shared.WorkflowBasicOut, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: time.Second * 5})
	utils.LogGreen("inside the basic workflow")

	// Writing to DB
	utils.LogDebug("calling the db activity")
	var dbOut shared.DBOut
	err := workflow.ExecuteActivity(ctx, activities.WriteToDB, data).Get(ctx, &dbOut)
	if err != nil {
		utils.LogRed(err)
	}
	utils.LogDebug("db out", dbOut)

	utils.LogDebug("calling the git activity")
	// Storing To Git
	var gitOut shared.GitOut
	err = workflow.ExecuteActivity(ctx, activities.WriteToGit, data, dbOut).Get(ctx, &gitOut)
	if err != nil {
		utils.LogRed(err)
	}
	utils.LogDebug("git out", gitOut)

	utils.LogGreen("completed the workflow")
	return shared.WorkflowBasicOut{
		DBOut: &dbOut, GitOut: &gitOut,
	}, err
}
