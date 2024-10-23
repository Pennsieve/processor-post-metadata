package mock

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/client/models"
	models2 "github.com/pennsieve/processor-post-metadata/service/models"
	"net/http"
)

func NewExpectedModelCreateCall(expectedDatasetID string, modelID string, expectedCreate models.ModelCreate) *ExpectedAPICall[models.ModelCreate, models2.APIResponse] {
	return &ExpectedAPICall[models.ModelCreate, models2.APIResponse]{
		Method:              http.MethodPost,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts", expectedDatasetID),
		ExpectedRequestBody: &expectedCreate,
		APIResponse: models2.APIResponse{
			Name: expectedCreate.Name,
			ID:   modelID,
		},
	}
}

func NewExpectedPropertiesCreateCall(expectedDatasetID string, modelID string, expectedCreate models.PropertiesCreate) *ExpectedAPICall[models.PropertiesCreate, []models2.APIResponse] {
	var apiResponse []models2.APIResponse
	for _, prop := range expectedCreate {
		apiResponse = append(apiResponse, models2.APIResponse{
			Name: prop.Name,
			ID:   uuid.NewString(),
		})
	}
	return &ExpectedAPICall[models.PropertiesCreate, []models2.APIResponse]{
		Method:              http.MethodPut,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/properties", expectedDatasetID, modelID),
		ExpectedRequestBody: &expectedCreate,
		APIResponse:         apiResponse,
	}
}

func NewExpectedRecordCreateCall(datasetID string, modelID string, expectedCreate models.RecordValues) *ExpectedAPICall[models.RecordValues, models2.APIResponse] {
	return &ExpectedAPICall[models.RecordValues, models2.APIResponse]{
		Method:              http.MethodPost,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/instances", datasetID, modelID),
		ExpectedRequestBody: &expectedCreate,
		APIResponse: models2.APIResponse{
			Name: uuid.NewString(),
			ID:   uuid.NewString(),
		},
	}
}

func NewExpectedRecordUpdateCall(datasetID string, modelID string, recordID models.PennsieveInstanceID, expectedUpdate models.RecordValues) *ExpectedAPICall[models.RecordValues, models2.APIResponse] {
	return &ExpectedAPICall[models.RecordValues, models2.APIResponse]{
		Method:              http.MethodPut,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/instances/%s", datasetID, modelID, recordID),
		ExpectedRequestBody: &expectedUpdate,
		APIResponse: models2.APIResponse{
			Name: uuid.NewString(),
			ID:   uuid.NewString(),
		},
	}
}
