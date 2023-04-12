package activities

import (
	"context"
	"crypto/sha1"
	"fmt"
	"learn_temporal/utils"
	"learn_temporal/workflowtype"
	"time"

	"github.com/google/uuid"
)

func WriteToDB(ctx context.Context, data workflowtype.WorkflowIn, waitFor time.Duration) (workflowtype.DBOut, error) {
	utils.LogDebug(fmt.Sprintf("writing '%s' to the database", *data.Data))
	id, err := uuid.NewUUID()

	// simulate writing data to the database
	time.Sleep(waitFor)
	idS := id.String()

	utils.LogGreen("completed the db activity")
	return workflowtype.DBOut{ID: &idS}, err
}

func WriteToGit(ctx context.Context, data workflowtype.WorkflowIn, dbOut workflowtype.DBOut, waitFor time.Duration) (workflowtype.GitOut, error) {
	utils.LogDebug("data will be read from", dbOut)
	utils.LogDebug("writing", *data.Data, "to git")

	// immitate writing to git and generating the sha1
	time.Sleep(waitFor)
	hash := sha1.Sum([]byte(*data.Data))
	hashString := fmt.Sprintf("%x", hash)

	utils.LogGreen("completed the git activity")
	return workflowtype.GitOut{
		ID: &hashString,
	}, nil
}

func WriteToDBWithSideEffect(ctx context.Context, data workflowtype.WorkflowIn, sideEffectIn workflowtype.SideEffectOut, waitFor time.Duration) (workflowtype.DBOutWithSideEffect, error) {
	utils.LogDebug(fmt.Sprintf("writing '%s' to the database", *data.Data))
	id, err := uuid.NewUUID()

	// simulate writing data to the database
	time.Sleep(waitFor)
	idS := id.String()

	utils.LogGreen("completed the db activity")
	return workflowtype.DBOutWithSideEffect{
		DBOut:         &workflowtype.DBOut{ID: &idS},
		SideEffectOut: &sideEffectIn,
	}, err
}

func ReadDBConfig() (result workflowtype.SideEffectOut) {
	sideEffectValue := "{'db': 'postgres'}"
	return workflowtype.SideEffectOut{
		Message: &sideEffectValue,
	}

	// f, err := os.Open("db.json") // magic fname string
	// errS := err.Error()
	// if err != nil {
	// 	return workflowtype.SideEffectOut{
	// 		Message: &sideEffectValue,
	// 		Error:   &errS,
	// 	}
	// }
	// defer f.Close()

	// content, err := io.ReadAll(f)
	// if err != nil {
	// 	errS = err.Error()
	// 	return workflowtype.SideEffectOut{
	// 		Message: &sideEffectValue,
	// 		Error:   &errS,
	// 	}
	// }
	// contentS := string(content)
	// utils.LogGreen("content from the file", contentS)
	// return workflowtype.SideEffectOut{
	// 	Message: &contentS,
	// 	Error:   &errS,
	// }
}
