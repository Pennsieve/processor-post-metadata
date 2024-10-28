package mock

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"net/http"
)

func NewExpectedCreateProxyRelationshipSchemaCall(datasetID string) *ExpectedAPICall[models.CreateProxyRelationshipSchemaBody, models.APIResponse] {
	expectedBody := models.NewCreateProxyRelationshipSchemaBody()
	return &ExpectedAPICall[models.CreateProxyRelationshipSchemaBody, models.APIResponse]{
		Method:              http.MethodPost,
		APIPath:             fmt.Sprintf("/models/datasets/%s/relationships", datasetID),
		ExpectedRequestBody: &expectedBody,
		APIResponse: models.APIResponse{
			Name: models.ProxyRelationshipSchemaName,
			ID:   uuid.NewString(),
		},
	}
}

func NewExpectedDeleteProxyInstancesCall(datasetID string, expectedBody ...models.DeleteProxyInstancesBody) *ExpectedAPICallMulti[models.DeleteProxyInstancesBody, any] {
	var calls []ExpectedAPICallData[models.DeleteProxyInstancesBody, any]
	for i := range expectedBody {
		calls = append(calls, ExpectedAPICallData[models.DeleteProxyInstancesBody, any]{
			Method:              http.MethodDelete,
			ExpectedRequestBody: &expectedBody[i],
		})
	}
	return &ExpectedAPICallMulti[models.DeleteProxyInstancesBody, any]{
		APIPath: fmt.Sprintf("/models/datasets/%s/proxy/package/instances/bulk", datasetID),
		Calls:   calls,
	}
}

func NewExpectedCreateProxyInstanceCall(datasetID string, expectedBody ...models.CreateProxyInstanceBody) *ExpectedAPICallMulti[models.CreateProxyInstanceBody, any] {
	var calls []ExpectedAPICallData[models.CreateProxyInstanceBody, any]
	for i := range expectedBody {
		calls = append(calls, ExpectedAPICallData[models.CreateProxyInstanceBody, any]{
			Method:              http.MethodPost,
			ExpectedRequestBody: &expectedBody[i],
		})
	}
	return &ExpectedAPICallMulti[models.CreateProxyInstanceBody, any]{
		APIPath: fmt.Sprintf("/models/datasets/%s/proxy/package/instances", datasetID),
		Calls:   calls,
	}
}
