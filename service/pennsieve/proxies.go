package pennsieve

import (
	"fmt"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"net/http"
)

func (s *Session) CreateProxyRelationshipSchema(datasetID string) (clientmodels.PennsieveSchemaID, error) {
	url := fmt.Sprintf("%s/models/datasets/%s/relationships", s.APIHost, datasetID)
	response, err := s.InvokePennsieve(http.MethodPost, url, models.NewCreateProxyRelationshipSchemaBody())
	if err != nil {
		return "", fmt.Errorf("error creating proxy relationship schema: %w", err)
	}
	apiResponse, err := handleResponseBody(response)
	if err != nil {
		return "", fmt.Errorf("error decoding response from create proxy relationship schema: %w", err)

	}
	return clientmodels.PennsieveSchemaID(apiResponse.ID), nil
}

func (s *Session) CreateProxyInstance(datasetID string, body models.CreateProxyInstanceBody) error {
	url := fmt.Sprintf("%s/models/datasets/%s/proxy/package/instances", s.APIHost, datasetID)
	_, err := s.InvokePennsieve(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("error creating proxy instance for package %s: %w",
			body.ExternalID,
			err)
	}
	return nil
}

func (s *Session) DeleteProxyInstances(datasetID string, body models.DeleteProxyInstancesBody) error {
	url := fmt.Sprintf("%s/models/datasets/%s/proxy/package/instances/bulk", s.APIHost, datasetID)
	_, err := s.InvokePennsieve(http.MethodDelete, url, body)
	if err != nil {
		return fmt.Errorf("error deleting proxy instances for record %s: %w", body.SourceRecordID, err)
	}
	return nil
}
