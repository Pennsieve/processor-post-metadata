package mock

import (
	"fmt"
	"github.com/google/uuid"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"net/http"
)

func NewExpectedCreateLinkSchemaCall(datasetID string, fromModelID clientmodels.PennsieveSchemaID, expectedRequestBody models.CreateLinkSchemaBody) *ExpectedAPICall[models.CreateLinkSchemaBody, models.APIResponse] {
	return &ExpectedAPICall[models.CreateLinkSchemaBody, models.APIResponse]{
		Method:              http.MethodPost,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/linked", datasetID, fromModelID),
		ExpectedRequestBody: &expectedRequestBody,
		APIResponse: models.APIResponse{
			Name: expectedRequestBody.Name,
			ID:   uuid.NewString(),
		},
	}
}

func NewExpectedCreateLinkInstanceCall(datasetID string, fromModelID clientmodels.PennsieveSchemaID, fromRecordID clientmodels.PennsieveInstanceID, expectedRequestBody models.CreateLinkInstanceBody) *ExpectedAPICall[models.CreateLinkInstanceBody, models.APIResponse] {
	return &ExpectedAPICall[models.CreateLinkInstanceBody, models.APIResponse]{
		Method:              http.MethodPost,
		APIPath:             fmt.Sprintf("/models/datasets/%s/concepts/%s/instances/%s/linked", datasetID, fromModelID, fromRecordID),
		ExpectedRequestBody: &expectedRequestBody,
		APIResponse: models.APIResponse{
			ID: uuid.NewString(),
		},
	}
}

func NewExpectedDeleteLinkInstanceCall(datasetID string, modelID clientmodels.PennsieveSchemaID, fromRecordID clientmodels.PennsieveInstanceID, linkInstanceID clientmodels.PennsieveInstanceID) *ExpectedAPICall[any, any] {
	return &ExpectedAPICall[any, any]{
		Method: http.MethodDelete,
		APIPath: fmt.Sprintf("/models/datasets/%s/concepts/%s/instances/%s/linked/%s",
			datasetID, modelID, fromRecordID, linkInstanceID),
	}
}
