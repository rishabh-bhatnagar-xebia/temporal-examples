package workflows

import (
	"learn_temporal/activities"
	"learn_temporal/utils"
	"learn_temporal/workflowtype"
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

	encoded := workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
		return activities.GitSideEffect(*data.Data)
	})
	var sideEffectValue string
	encoded.Get(&sideEffectValue)

	utils.LogDebug("WORKFLOW;", "calling the git activity")
	// Storing To Git
	var gitOut workflowtype.GitOutWithSideEffect
	err = workflow.ExecuteActivity(ctx, activities.WriteToGit, data, dbOut, time.Second*3).Get(ctx, &gitOut)
	if err != nil {
		utils.LogRed("WORKFLOW;", err)
	}
	utils.LogDebug("WORKFLOW;", "git out", gitOut)

	utils.LogGreen("WORKFLOW;", "completed the workflow")
	return workflowtype.WorkflowBasicOut{
		DBOut: &dbOut, GitOut: &workflowtype.GitOutWithSideEffect{
			GitOut:        gitOut.GitOut,
			SideEffectOut: &workflowtype.SideEffectOut{Out: &sideEffectValue},
		},
	}, err
}
