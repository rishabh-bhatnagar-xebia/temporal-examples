package workflows

import (
	"async/activities"
	workflowtype "async/protoc_types"
	"async/utils"
	"time"

	"go.temporal.io/sdk/workflow"
)

func Basic(ctx workflow.Context, data workflowtype.WorkflowIn) (workflowtype.WorkflowBasicOut, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: time.Second * 25})
	utils.LogGreen("WORKFLOW;", "inside the basic workflow")

	// Writing to DB
	utils.LogDebug("WORKFLOW;", "calling the db activity")
	var dbOut workflowtype.DBOut
	err := workflow.ExecuteActivity(ctx, activities.WriteToDB, data, time.Second*0).Get(ctx, &dbOut)
	if err != nil {
		utils.LogRed("WORKFLOW;", err)
	}
	utils.LogDebug("WORKFLOW;", "db out", dbOut)

	utils.LogDebug("WORKFLOW;", "calling the git activity")
	// Storing To Git
	var gitOut workflowtype.GitOut
	err = workflow.ExecuteActivity(ctx, activities.WriteToGit, data, dbOut, time.Second*3).Get(ctx, &gitOut)
	if err != nil {
		utils.LogRed("WORKFLOW;", err)
	}
	utils.LogDebug("WORKFLOW;", "git out", gitOut)

	utils.LogGreen("WORKFLOW;", "completed the workflow")
	return workflowtype.WorkflowBasicOut{
		DBOut: &dbOut, GitOut: &gitOut,
	}, err
}
