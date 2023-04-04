package activities

import (
	workflowtype "async/protoc_types"
	"async/utils"
	"context"
	"crypto/sha1"
	"fmt"
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

func WriteToGit(ctx context.Context, data workflowtype.WorkflowIn, dbOut workflowtype.DBOut, waitFor time.Duration) (workflowtype.GitOut, error) {
	utils.LogDebug("data will be read from", dbOut)
	utils.LogDebug("writing", data.Data, "to git")
	time.Sleep(waitFor)
	hash := sha1.Sum([]byte(*data.Data))
	hashString := fmt.Sprintf("%x", hash)
	utils.LogGreen("completed the git activity")
	return workflowtype.GitOut{ID: &hashString}, nil
}
