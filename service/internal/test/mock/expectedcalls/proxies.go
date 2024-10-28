package expectedcalls

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/service/internal/test/mock"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"net/http"
)

func CreateProxyRelationshipSchema(datasetID string) *mock.ExpectedAPICall[models.CreateProxyRelationshipSchemaBody, models.APIResponse] {
	expectedBody := models.NewCreateProxyRelationshipSchemaBody()
	return &mock.ExpectedAPICall[models.CreateProxyRelationshipSchemaBody, models.APIResponse]{
		Method:              http.MethodPost,
		APIPath:             fmt.Sprintf("/models/datasets/%s/relationships", datasetID),
		ExpectedRequestBody: &expectedBody,
		APIResponse: models.APIResponse{
			Name: models.ProxyRelationshipSchemaName,
			ID:   uuid.NewString(),
		},
	}
}

func DeleteProxyInstances(datasetID string, expectedBody ...models.DeleteProxyInstancesBody) *mock.ExpectedAPICallMulti[models.DeleteProxyInstancesBody, any] {
	var calls []mock.ExpectedAPICallData[models.DeleteProxyInstancesBody, any]
	for i := range expectedBody {
		calls = append(calls, mock.ExpectedAPICallData[models.DeleteProxyInstancesBody, any]{
			Method:              http.MethodDelete,
			ExpectedRequestBody: &expectedBody[i],
		})
	}
	return &mock.ExpectedAPICallMulti[models.DeleteProxyInstancesBody, any]{
		APIPath: fmt.Sprintf("/models/datasets/%s/proxy/package/instances/bulk", datasetID),
		Calls:   calls,
	}
}

func CreateProxyInstance(datasetID string, expectedBody ...models.CreateProxyInstanceBody) *mock.ExpectedAPICallMulti[models.CreateProxyInstanceBody, any] {
	var calls []mock.ExpectedAPICallData[models.CreateProxyInstanceBody, any]
	for i := range expectedBody {
		calls = append(calls, mock.ExpectedAPICallData[models.CreateProxyInstanceBody, any]{
			Method:              http.MethodPost,
			ExpectedRequestBody: &expectedBody[i],
		})
	}
	return &mock.ExpectedAPICallMulti[models.CreateProxyInstanceBody, any]{
		APIPath: fmt.Sprintf("/models/datasets/%s/proxy/package/instances", datasetID),
		Calls:   calls,
	}
}
