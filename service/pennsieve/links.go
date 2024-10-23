package pennsieve

import (
	"fmt"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"net/http"
)

func (s *Session) CreateLinkedPropertySchema(datasetID string, fromModelID clientmodels.PennsieveSchemaID, body models.CreateLinkSchemaBody) (clientmodels.PennsieveSchemaID, error) {
	url := fmt.Sprintf("%s/models/datasets/%s/concepts/%s/linked", s.APIHost, datasetID, fromModelID)
	response, err := s.InvokePennsieve(http.MethodPost, url, body)
	if err != nil {
		return "", fmt.Errorf("error creating linked property schema %s: %w", body.Name, err)
	}
	apiResponse, err := handleResponseBody(response)
	if err != nil {
		return "", fmt.Errorf("error decoding create linked propert schema response for %s: %w", body.Name, err)
	}
	return clientmodels.PennsieveSchemaID(apiResponse.ID), nil
}

func (s *Session) CreateLinkedPropertyInstance(datasetID string, fromModelID clientmodels.PennsieveSchemaID, fromRecordID clientmodels.PennsieveInstanceID, body models.CreateLinkInstanceBody) error {
	url := fmt.Sprintf("%s/models/datasets/%s/concepts/%s/instances/%s/linked", s.APIHost, datasetID, fromModelID, fromRecordID)
	_, err := s.InvokePennsieve(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("error creating linked property %s instance from record %s to record %s: %w",
			body.SchemaLinkedPropertyId, fromRecordID, body.To, err)
	}
	return nil
}
