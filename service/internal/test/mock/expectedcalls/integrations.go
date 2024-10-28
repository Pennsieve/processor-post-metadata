package expectedcalls

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/processor-post-metadata/service/internal/test/mock"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"net/http"
)

func GetIntegration(integrationID, datasetID string) *mock.ExpectedAPICall[any, models.Integration] {
	return &mock.ExpectedAPICall[any, models.Integration]{
		Method:              http.MethodGet,
		APIPath:             fmt.Sprintf("/integrations/%s", integrationID),
		ExpectedRequestBody: nil,
		APIResponse: models.Integration{
			Uuid:          uuid.NewString(),
			DatasetNodeID: datasetID,
		},
	}
}
