package workflows

import (
	"async/activities"
	"async/shared"
	"async/utils"
	"go.temporal.io/sdk/workflow"
	"time"
)

func AsyncWithChild(ctx workflow.Context, data shared.WorkflowIn) (shared.WorkflowBasicOut, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: time.Second * 5})
	utils.LogGreen("inside the basic workflow")

	// Writing to DB
	utils.LogDebug("calling the db activity")
	var dbOut shared.DBOut
	err := workflow.ExecuteActivity(ctx, activities.WriteToDB, data).Get(ctx, &dbOut)
	if err != nil {
		utils.LogRed(err)
		return shared.WorkflowBasicOut{}, err
	}
	utils.LogDebug("db out", dbOut)

	// Storing To Git
	utils.LogDebug("calling the git child workflow")
	var gitOut shared.GitOut
	ctxForChild := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{WorkflowID: "workflow-id-db"})
	err = workflow.ExecuteChildWorkflow(ctxForChild, GitWorkflow, data, dbOut).Get(ctx, &gitOut)
	if err != nil {
		utils.LogRed(err)
		return shared.WorkflowBasicOut{DBOut: &dbOut, GitOut: nil}, err
	}
	utils.LogDebug("git out", gitOut)

	utils.LogGreen("completed the workflow")
	return shared.WorkflowBasicOut{
		DBOut: &dbOut, GitOut: &gitOut,
	}, err
}

func GitWorkflow(ctx workflow.Context, data shared.WorkflowIn, dbOut shared.DBOut) (*shared.GitOut, error) {
	utils.LogDebug("executing the git-child-workflow")
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: time.Second * 5})
	var gitOut shared.GitOut
	err := workflow.ExecuteActivity(ctx, activities.WriteToGit, data, dbOut).Get(ctx, &gitOut)
	if err != nil {
		utils.LogRed(err)
		return nil, err
	}
	utils.LogDebug("executed the child workflow for git")
	return &gitOut, nil
}
