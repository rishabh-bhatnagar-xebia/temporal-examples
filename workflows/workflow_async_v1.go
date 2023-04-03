package workflows

import (
	"async/activities"
	"async/shared"
	"async/utils"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

func AsyncWithChild(ctx workflow.Context, data shared.WorkflowIn) (shared.WorkflowAsyncV1Out, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 25,
	})
	utils.LogGreen("WORKFLOW:", "inside the basic workflow")

	// Writing to DB
	utils.LogDebug("WORKFLOW:", "calling the db activity")
	var dbOut shared.DBOut
	err := workflow.ExecuteActivity(ctx, activities.WriteToDB, data, time.Second*0).Get(ctx, &dbOut)
	if err != nil {
		utils.LogRed("WORKFLOW:", err)
		return shared.WorkflowAsyncV1Out{}, err
	}
	utils.LogDebug("WORKFLOW:", "db out", dbOut)

	// Storing To Git
	utils.LogDebug("WORKFLOW:", "calling the git child workflow")
	cwo := workflow.ChildWorkflowOptions{
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
		TaskQueue:         shared.QueueNameAsyncV1,
	}
	ctxForChild := workflow.WithChildOptions(ctx, cwo)
	err = workflow.ExecuteChildWorkflow(ctxForChild, GitWorkflow, data, dbOut).GetChildWorkflowExecution().Get(ctx, nil)
	if err != nil {
		utils.LogRed("WORKFLOW:", err)
	}

	utils.LogGreen("WORKFLOW:", "completed the workflow")
	return shared.WorkflowAsyncV1Out{
		DBOut: &dbOut,
	}, err
}

func GitWorkflow(ctx workflow.Context, data shared.WorkflowIn, dbOut shared.DBOut) (*shared.GitOut, error) {
	utils.LogDebug("GIT-WORKFLOW:", "executing the git-child-workflow")
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: time.Second * 15, TaskQueue: shared.QueueNameAsyncV1})
	var gitOut shared.GitOut
	err := workflow.ExecuteActivity(ctx, activities.WriteToGit, data, dbOut, time.Second*5).Get(ctx, &gitOut)
	if err != nil {
		utils.LogRed(err)
		return nil, err
	}
	utils.LogDebug("GIT-WORKFLOW", "executed the child workflow for git")
	return &gitOut, nil
}
