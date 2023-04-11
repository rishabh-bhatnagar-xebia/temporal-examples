package activities

import (
	"context"
	"crypto/sha1"
	"fmt"
	"learn_temporal/utils"
	"learn_temporal/workflowtype"
	"os"
	"time"

	"github.com/google/uuid"
)

func WriteToDB(ctx context.Context, data workflowtype.WorkflowIn, waitFor time.Duration) (workflowtype.DBOut, error) {
	utils.LogDebug("writing", data, "to the database")
	id, err := uuid.NewUUID()
	time.Sleep(waitFor)
	idS := id.String()
	utils.LogGreen("completed the db activity")
	return workflowtype.DBOut{ID: &idS}, err
}

func WriteToGit(ctx context.Context, data workflowtype.WorkflowIn, dbOut workflowtype.DBOut, waitFor time.Duration) (workflowtype.GitOutWithSideEffect, error) {
	utils.LogDebug("data will be read from", dbOut)
	utils.LogDebug("writing", data.Data, "to git")
	time.Sleep(waitFor)
	hash := sha1.Sum([]byte(*data.Data))
	hashString := fmt.Sprintf("%x", hash)
	utils.LogGreen("completed the git activity")
	sideEffectValue := "value"
	return workflowtype.GitOutWithSideEffect{
		GitOut: &workflowtype.GitOut{
			ID: &hashString,
		},
		SideEffectOut: &workflowtype.SideEffectOut{
			Out: &sideEffectValue,
		},
	}, nil
}

func GitSideEffect(content string) string {
	f, err := os.Create("temp")
	if err != nil {
		return err.Error()
	}

	_, err = f.WriteString(content)
	if err != nil {
		return err.Error()
	}

	f.Sync()
	return "data written to file successfully"
}
