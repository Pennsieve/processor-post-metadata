package processor

import (
	"fmt"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"log/slog"
)

func (p *MetadataPostProcessor) ProcessModels(datasetID string, modelChanges []clientmodels.ModelChanges) error {
	if len(modelChanges) == 0 {
		logger.Info("no model changes")
		return nil
	}
	logger.Info("starting model changes")
	for _, modelChange := range modelChanges {
		if err := p.ProcessModelChanges(datasetID, modelChange); err != nil {
			return err
		}
	}
	logger.Info("finished model changes")
	return nil
}

func (p *MetadataPostProcessor) ProcessModelChanges(datasetID string, modelChanges clientmodels.ModelChanges) error {
	if err := p.CreateModelIfNecessary(datasetID, modelChanges); err != nil {
		return err
	}
	return nil
}

func (p *MetadataPostProcessor) CreateModelIfNecessary(datasetID string, modelChange clientmodels.ModelChanges) error {
	if modelChange.Create == nil {
		logger.Info("model already exists", slog.String("modelID", modelChange.ID))
		return nil
	}
	modelCreate := modelChange.Create
	modelLogger := logger.With(slog.String("modelName", modelCreate.Model.Name))
	modelLogger.Info("creating model")
	modelID, err := p.Pennsieve.CreateModelAndProps(datasetID, *modelCreate)
	if err != nil {
		return fmt.Errorf("error creating model: %w", err)
	}
	p.IDStore.AddModel(modelCreate.Model.Name, modelID)
	modelLogger.Info("mmodel created", slog.Any("modelID", modelID))
	return nil
}
