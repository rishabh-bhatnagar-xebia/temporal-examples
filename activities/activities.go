package activities

import (
	"async/shared"
	"async/utils"
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func WriteToDB(ctx context.Context, data shared.WorkflowIn, waitFor time.Duration) (shared.DBOut, error) {
	utils.LogDebug("writing", data, "to the database")
	id, err := uuid.NewUUID()
	time.Sleep(waitFor)
	return shared.DBOut{ID: "dbId: " + id.String()}, err
}

func WriteToGit(ctx context.Context, data shared.WorkflowIn, dbOut shared.DBOut, waitFor time.Duration) (shared.GitOut, error) {
	utils.LogDebug("data will be read from", dbOut)
	utils.LogDebug("writing", data.Data, "to git")
	time.Sleep(waitFor)
	hash := sha1.Sum([]byte(data.Data))
	return shared.GitOut{ID: fmt.Sprintf("%x", hash)}, nil
}
