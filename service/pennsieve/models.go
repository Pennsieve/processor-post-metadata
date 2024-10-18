package pennsieve

import (
	"encoding/json"
	"fmt"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"github.com/pennsieve/processor-post-metadata/service/util"
	"net/http"
)

func (s *Session) CreateModelAndProps(datasetID string, modelPropsCreate clientmodels.ModelPropsCreate) (clientmodels.PennsieveInstanceID, error) {
	modelID, err := s.CreateModel(datasetID, modelPropsCreate.Model)
	if err != nil {
		return "", err
	}
	if err := s.CreateModelProperties(datasetID, modelID, modelPropsCreate.Properties); err != nil {
		return "", fmt.Errorf("model %s created; error creating properties: %w", modelPropsCreate.Model.Name, err)
	}
	return modelID, nil
}

func (s *Session) CreateModel(datasetID string, modelCreate clientmodels.ModelCreate) (clientmodels.PennsieveInstanceID, error) {
	url := fmt.Sprintf("%s/models/datasets/%s/concepts", s.APIHost, datasetID)
	response, err := s.InvokePennsieve(http.MethodPost, url, modelCreate)
	if err != nil {
		return "", fmt.Errorf("error creating model %s: %w", modelCreate.Name, err)
	}
	defer util.CloseAndWarn(response)

	var apiResponse models.APIResponse
	if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
		return "", fmt.Errorf("error decoding create model response for %s: %w", modelCreate.Name, err)
	}
	return apiResponse.ID, nil
}

func (s *Session) CreateModelProperties(datasetID string, modelID clientmodels.PennsieveInstanceID, propsCreate clientmodels.PropertiesCreate) error {
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
