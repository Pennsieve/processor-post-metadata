package processortest

import (
	"fmt"
	"github.com/google/uuid"
)

func NewDatasetID() string {
	return NewDatasetIDWithUUID(uuid.NewString())
}

func NewDatasetIDWithUUID(datasetUUID string) string {
	return fmt.Sprintf("N:dataset:%s", datasetUUID)
}
