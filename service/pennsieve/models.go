package pennsieve

import (
	"encoding/json"
	"errors"
	"fmt"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"github.com/pennsieve/processor-post-metadata/service/util"
	"net/http"
)

func (s *Session) CreateModelAndProps(datasetID string, modelPropsCreate clientmodels.ModelPropsCreate) (clientmodels.PennsieveSchemaID, error) {
	modelID, err := s.CreateModel(datasetID, modelPropsCreate.Model)
	if err != nil {
		return "", err
	}
	if err := s.CreateModelProperties(datasetID, modelID, modelPropsCreate.Properties); err != nil {
		return "", fmt.Errorf("model %s created; error creating properties: %w", modelPropsCreate.Model.Name, err)
	}
	return modelID, nil
}

func (s *Session) CreateModel(datasetID string, modelCreate clientmodels.ModelCreate) (clientmodels.PennsieveSchemaID, error) {
	url := fmt.Sprintf("%s/models/datasets/%s/concepts", s.APIHost, datasetID)
	response, err := s.InvokePennsieve(http.MethodPost, url, modelCreate)
	if err != nil {
		return "", fmt.Errorf("error creating model %s: %w", modelCreate.Name, err)
	}
	apiResponse, err := handleResponseBody(response)
	if err != nil {
		return "", fmt.Errorf("error decoding create model response for %s: %w", modelCreate.Name, err)
	}
	return clientmodels.PennsieveSchemaID(apiResponse.ID), nil
}

func (s *Session) CreateModelProperties(datasetID string, modelID clientmodels.PennsieveSchemaID, propsCreate clientmodels.PropertiesCreate) error {
	if len(propsCreate) == 0 {
		return nil
	}
	url := fmt.Sprintf("%s/models/datasets/%s/concepts/%s/properties", s.APIHost, datasetID, modelID)
	_, err := s.InvokePennsieve(http.MethodPut, url, propsCreate)
	if err != nil {
		return fmt.Errorf("error creating properties for modelID %s: %w", modelID, err)
	}
	return nil
}

func (s *Session) CreateRecord(datasetID string, modelID clientmodels.PennsieveSchemaID, values clientmodels.RecordValues) (clientmodels.PennsieveInstanceID, error) {
	url := fmt.Sprintf("%s/models/datasets/%s/concepts/%s/instances", s.APIHost, datasetID, modelID)
	response, err := s.InvokePennsieve(http.MethodPost, url, values)
	if err != nil {
		return "", fmt.Errorf("error creating record for model %s: %w", modelID, err)
	}
	apiResponse, err := handleResponseBody(response)
	if err != nil {
		return "", fmt.Errorf("error decoding create record response for model %s: %w", modelID, err)
	}
	return clientmodels.PennsieveInstanceID(apiResponse.ID), nil
}

func (s *Session) UpdateRecord(datasetID string, modelID clientmodels.PennsieveSchemaID, recordID clientmodels.PennsieveInstanceID, values clientmodels.RecordValues) (clientmodels.PennsieveInstanceID, error) {
	url := fmt.Sprintf("%s/models/datasets/%s/concepts/%s/instances/%s",
		s.APIHost,
		datasetID,
		modelID,
		recordID)
	response, err := s.InvokePennsieve(http.MethodPut, url, values)
	if err != nil {
		return "", fmt.Errorf("error updating record %s for model %s: %w", recordID, modelID, err)
	}
	apiResponse, err := handleResponseBody(response)
	if err != nil {
		return "", fmt.Errorf("error decoding update record %s response for model %s: %w", recordID, modelID, err)
	}
	return clientmodels.PennsieveInstanceID(apiResponse.Name), nil
}

type bulkDeleteResponse struct {
	Success []clientmodels.PennsieveInstanceID `json:"success"`
	// Errors is a slice of slices. Each slice in the outer slice should be of the form [instance-id, error-message]
	Errors [][]string `json:"errors"`
}

func (s *Session) DeleteRecords(datasetID string, modelID clientmodels.PennsieveSchemaID, recordIDs []clientmodels.PennsieveInstanceID) error {
	url := fmt.Sprintf("%s/models/datasets/%s/concepts/%s/instances", s.APIHost, datasetID, modelID)
	response, err := s.InvokePennsieve(http.MethodDelete, url, recordIDs)
	if err != nil {
		return fmt.Errorf("error deleting %d records for model %s: %w", len(recordIDs), modelID, err)
	}

	defer util.CloseAndWarn(response)

	var bulkResponse bulkDeleteResponse
	if err := json.NewDecoder(response.Body).Decode(&bulkResponse); err != nil {
		return fmt.Errorf("error decoding response from deleting %d records for model %s: %w", len(recordIDs), modelID, err)
	}

	if len(bulkResponse.Errors) == 0 {
		return nil
	}

	var errs []error
	errs = append(errs, fmt.Errorf("errors deleting %d of %d records for model %s",
		len(bulkResponse.Errors),
		len(recordIDs),
		modelID))

	for _, errResp := range bulkResponse.Errors {
		errs = append(errs, fmt.Errorf("error deleting record %s: %s",
			errResp[0],
			errResp[1],
		))
	}
	return errors.Join(errs...)
}

func handleResponseBody(response *http.Response) (models.APIResponse, error) {
	defer util.CloseAndWarn(response)

	var apiResponse models.APIResponse
	if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
		return models.APIResponse{}, err
	}
	return apiResponse, nil
}
