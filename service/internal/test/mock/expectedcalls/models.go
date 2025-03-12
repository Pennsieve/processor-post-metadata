package expectedcalls

import (
	"fmt"
	"github.com/google/uuid"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/internal/test/mock"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"net/http"
)

func ModelCreate(expectedDatasetID string, modelID clientmodels.PennsieveSchemaID, expectedCreate clientmodels.ModelCreateParams) *mock.ExpectedAPICall[clientmodels.ModelCreateParams, models.APIResponse] {
	return &mock.ExpectedAPICall[clientmodels.ModelCreateParams, models.APIResponse]{
		Method:              http.MethodPost,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts", expectedDatasetID),
		ExpectedRequestBody: &expectedCreate,
		APIResponse: models.APIResponse{
			Name: expectedCreate.Name,
			ID:   modelID.String(),
		},
	}
}

func PropertiesCreate(expectedDatasetID string, modelID clientmodels.PennsieveSchemaID, expectedCreate clientmodels.PropertiesCreateParams) *mock.ExpectedAPICall[clientmodels.PropertiesCreateParams, []models.APIResponse] {
	var apiResponse []models.APIResponse
	for _, prop := range expectedCreate {
		apiResponse = append(apiResponse, models.APIResponse{
			Name: prop.Name,
			ID:   uuid.NewString(),
		})
	}
	return &mock.ExpectedAPICall[clientmodels.PropertiesCreateParams, []models.APIResponse]{
		Method:              http.MethodPut,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/properties", expectedDatasetID, modelID),
		ExpectedRequestBody: &expectedCreate,
		APIResponse:         apiResponse,
	}
}

func RecordCreate(datasetID string, modelID clientmodels.PennsieveSchemaID, expectedCreate clientmodels.RecordValues) *mock.ExpectedAPICall[clientmodels.RecordValues, models.APIResponse] {
	return &mock.ExpectedAPICall[clientmodels.RecordValues, models.APIResponse]{
		Method:              http.MethodPost,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/instances", datasetID, modelID),
		ExpectedRequestBody: &expectedCreate,
		APIResponse: models.APIResponse{
			Name: uuid.NewString(),
			ID:   uuid.NewString(),
		},
	}
}

func RecordUpdate(datasetID string, modelID clientmodels.PennsieveSchemaID, recordID clientmodels.PennsieveInstanceID, expectedUpdate clientmodels.RecordValues) *mock.ExpectedAPICall[clientmodels.RecordValues, models.APIResponse] {
	return &mock.ExpectedAPICall[clientmodels.RecordValues, models.APIResponse]{
		Method:              http.MethodPut,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/instances/%s", datasetID, modelID, recordID),
		ExpectedRequestBody: &expectedUpdate,
		APIResponse: models.APIResponse{
			Name: uuid.NewString(),
			ID:   uuid.NewString(),
		},
	}
}

func RecordDelete(datasetID string, modelID clientmodels.PennsieveSchemaID, expectedDelete []clientmodels.PennsieveInstanceID) *mock.ExpectedAPICall[[]clientmodels.PennsieveInstanceID, models.BulkDeleteRecordsResponse] {
	return &mock.ExpectedAPICall[[]clientmodels.PennsieveInstanceID, models.BulkDeleteRecordsResponse]{
		Method:              http.MethodDelete,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/instances", datasetID, modelID),
		ExpectedRequestBody: &expectedDelete,
		APIResponse: models.BulkDeleteRecordsResponse{
			Success: expectedDelete,
		},
	}
}

func RecordDeleteFailure(datasetID string, modelID clientmodels.PennsieveSchemaID, expectedDelete []clientmodels.PennsieveInstanceID) *mock.ExpectedAPICall[[]clientmodels.PennsieveInstanceID, models.BulkDeleteRecordsResponse] {
	var errs [][]string
	for _, recordID := range expectedDelete {
		errs = append(errs, []string{string(recordID), "could not delete record"})
	}
	return &mock.ExpectedAPICall[[]clientmodels.PennsieveInstanceID, models.BulkDeleteRecordsResponse]{
		Method:              http.MethodDelete,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/instances", datasetID, modelID),
		ExpectedRequestBody: &expectedDelete,
		APIResponse: models.BulkDeleteRecordsResponse{
			Errors: errs,
		},
	}
}
