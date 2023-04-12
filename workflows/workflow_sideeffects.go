package workflows

import (
	"learn_temporal/activities"
	"learn_temporal/utils"
	"learn_temporal/workflowtype"
	"time"

	"go.temporal.io/sdk/workflow"
)

func WithSideEffects(ctx workflow.Context, data workflowtype.WorkflowIn) (workflowtype.WorkflowSideEffectOut, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: time.Second * 25})
	utils.LogGreen("WORKFLOW;", "inside the side-effects workflow")

	// calling the side-effect to load the config
	encoded := workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
		return activities.ReadDBConfig()
	})
	utils.LogDebug("done the side effect")
	var sideEffectData workflowtype.SideEffectOut
	encoded.Get(&sideEffectData)
	utils.LogDebug("sideEffectDta", *sideEffectData.Message)
	utils.LogDebug("sideEffectData *****", sideEffectData)

	// Writing to DB
	utils.LogDebug("WORKFLOW;", "calling the db activity")
	var dbAndSideEffectOut workflowtype.DBOutWithSideEffect
	err := workflow.ExecuteActivity(ctx, activities.WriteToDBWithSideEffect, data, sideEffectData, time.Second*0).Get(ctx, &dbAndSideEffectOut)
	if err != nil {
		utils.LogRed("WORKFLOW;", err)
	}
	utils.LogDebug("WORKFLOW;", "db out", dbAndSideEffectOut)

	// Storing To Git
	utils.LogDebug("WORKFLOW;", "calling the git activity")
	var gitOut workflowtype.GitOut
	err = workflow.ExecuteActivity(ctx, activities.WriteToGit, data, dbAndSideEffectOut.DBOut, time.Second*3).Get(ctx, &gitOut)
	if err != nil {
		utils.LogRed("WORKFLOW;", err)
	}
	utils.LogDebug("WORKFLOW;", "git out", gitOut)

	utils.LogGreen("WORKFLOW;", "completed the workflow")
	return workflowtype.WorkflowSideEffectOut{
		GitOut:        &gitOut,
		DBOut:         dbAndSideEffectOut.DBOut,
		SideEffectOut: dbAndSideEffectOut.SideEffectOut,
	}, err
}
