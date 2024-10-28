package mock

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"net/http"
)

func NewExpectedGetIntegrationCall(integrationID, datasetID string) *ExpectedAPICall[any, models.Integration] {
	return &ExpectedAPICall[any, models.Integration]{
		Method:              http.MethodGet,
		APIPath:             fmt.Sprintf("/integrations/%s", integrationID),
		ExpectedRequestBody: nil,
		APIResponse: models.Integration{
			Uuid:          uuid.NewString(),
			DatasetNodeID: datasetID,
		},
	}
}
