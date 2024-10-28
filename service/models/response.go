package models

import clientmodels "github.com/pennsieve/processor-post-metadata/client/models"

// APIResponse is a minimal type that will only capture the name and id of
// model service responses. Create new types that build on this one if more
// info from a particular response is required.
type APIResponse struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type BulkDeleteRecordsResponse struct {
	Success []clientmodels.PennsieveInstanceID `json:"success"`
	// Errors is a slice of slices. Each slice in the outer slice should be of the form [instance-id, error-message]
	Errors [][]string `json:"errors"`
}
