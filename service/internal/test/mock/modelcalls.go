package mock

import (
	"fmt"
	"github.com/google/uuid"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"net/http"
)

func NewExpectedModelCreateCall(expectedDatasetID string, modelID clientmodels.PennsieveSchemaID, expectedCreate clientmodels.ModelCreate) *ExpectedAPICall[clientmodels.ModelCreate, models.APIResponse] {
	return &ExpectedAPICall[clientmodels.ModelCreate, models.APIResponse]{
		Method:              http.MethodPost,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts", expectedDatasetID),
		ExpectedRequestBody: &expectedCreate,
		APIResponse: models.APIResponse{
			Name: expectedCreate.Name,
			ID:   modelID.String(),
		},
	}
}

func NewExpectedPropertiesCreateCall(expectedDatasetID string, modelID clientmodels.PennsieveSchemaID, expectedCreate clientmodels.PropertiesCreate) *ExpectedAPICall[clientmodels.PropertiesCreate, []models.APIResponse] {
	var apiResponse []models.APIResponse
	for _, prop := range expectedCreate {
		apiResponse = append(apiResponse, models.APIResponse{
			Name: prop.Name,
			ID:   uuid.NewString(),
		})
	}
	return &ExpectedAPICall[clientmodels.PropertiesCreate, []models.APIResponse]{
		Method:              http.MethodPut,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/properties", expectedDatasetID, modelID),
		ExpectedRequestBody: &expectedCreate,
		APIResponse:         apiResponse,
	}
}

func NewExpectedRecordCreateCall(datasetID string, modelID clientmodels.PennsieveSchemaID, expectedCreate clientmodels.RecordValues) *ExpectedAPICall[clientmodels.RecordValues, models.APIResponse] {
	return &ExpectedAPICall[clientmodels.RecordValues, models.APIResponse]{
		Method:              http.MethodPost,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/instances", datasetID, modelID),
		ExpectedRequestBody: &expectedCreate,
		APIResponse: models.APIResponse{
			Name: uuid.NewString(),
			ID:   uuid.NewString(),
		},
	}
}

func NewExpectedRecordUpdateCall(datasetID string, modelID clientmodels.PennsieveSchemaID, recordID clientmodels.PennsieveInstanceID, expectedUpdate clientmodels.RecordValues) *ExpectedAPICall[clientmodels.RecordValues, models.APIResponse] {
	return &ExpectedAPICall[clientmodels.RecordValues, models.APIResponse]{
		Method:              http.MethodPut,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/instances/%s", datasetID, modelID, recordID),
		ExpectedRequestBody: &expectedUpdate,
		APIResponse: models.APIResponse{
			Name: uuid.NewString(),
			ID:   uuid.NewString(),
		},
	}
}

func NewExpectedRecordDeleteCall(datasetID string, modelID clientmodels.PennsieveSchemaID, expectedDelete []clientmodels.PennsieveInstanceID) *ExpectedAPICall[[]clientmodels.PennsieveInstanceID, models.BulkDeleteRecordsResponse] {
	return &ExpectedAPICall[[]clientmodels.PennsieveInstanceID, models.BulkDeleteRecordsResponse]{
		Method:              http.MethodDelete,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/instances", datasetID, modelID),
		ExpectedRequestBody: &expectedDelete,
		APIResponse: models.BulkDeleteRecordsResponse{
			Success: expectedDelete,
		},
	}
}

func NewExpectedRecordDeleteCallFailure(datasetID string, modelID clientmodels.PennsieveSchemaID, expectedDelete []clientmodels.PennsieveInstanceID) *ExpectedAPICall[[]clientmodels.PennsieveInstanceID, models.BulkDeleteRecordsResponse] {
	var errs [][]string
	for _, recordID := range expectedDelete {
		errs = append(errs, []string{string(recordID), "could not delete record"})
	}
	return &ExpectedAPICall[[]clientmodels.PennsieveInstanceID, models.BulkDeleteRecordsResponse]{
		Method:              http.MethodDelete,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/instances", datasetID, modelID),
		ExpectedRequestBody: &expectedDelete,
		APIResponse: models.BulkDeleteRecordsResponse{
			Errors: errs,
		},
	}
}
