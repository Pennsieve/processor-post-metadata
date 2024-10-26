package mock

import (
	"fmt"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"net/http"
)

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
