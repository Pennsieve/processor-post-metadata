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
	modelID, err := p.CreateModelIfNecessary(datasetID, modelChanges)
	if err != nil {
		return err
	}
	modelLogger := logger.With(slog.Any("modelID", modelID))
	modelLogger.Info("creating records")
	for _, recordCreate := range modelChanges.Records.Create {
		if err := p.CreateRecord(datasetID, modelID, recordCreate); err != nil {
			return err
		}
	}
	modelLogger.Info("created records", slog.Int("count", len(modelChanges.Records.Create)))
	modelLogger.Info("updating records")
	for _, recordUpdate := range modelChanges.Records.Update {
		if err := p.UpdateRecord(datasetID, modelID, recordUpdate); err != nil {
			return err
		}
	}
	modelLogger.Info("updated models", slog.Int("count", len(modelChanges.Records.Update)))
	// TODO handle deletes. I think the link and proxy deletes need to happen first

	return nil
}

func (p *MetadataPostProcessor) CreateModelIfNecessary(datasetID string, modelChange clientmodels.ModelChanges) (clientmodels.PennsieveSchemaID, error) {
	if modelChange.Create == nil {
		logger.Info("model already exists", slog.String("modelID", modelChange.ID))
		return clientmodels.PennsieveSchemaID(modelChange.ID), nil
	}
	modelCreate := modelChange.Create
	modelLogger := logger.With(slog.String("modelName", modelCreate.Model.Name))
	modelLogger.Info("creating model")
	modelID, err := p.Pennsieve.CreateModelAndProps(datasetID, *modelCreate)
	if err != nil {
		return "", fmt.Errorf("error creating model: %w", err)
	}
	p.IDStore.AddModel(modelCreate.Model.Name, modelID)
	modelLogger.Info("model created", slog.Any("modelID", modelID))
	return modelID, nil
}

func (p *MetadataPostProcessor) CreateRecord(datasetID string, modelID clientmodels.PennsieveSchemaID, recordCreate clientmodels.RecordCreate) error {
	recordID, err := p.Pennsieve.CreateRecord(datasetID, modelID, recordCreate.RecordValues)
	if err != nil {
		return err
	}
	p.IDStore.AddRecord(modelID, recordCreate.ExternalID, recordID)
	return nil
}

func (p *MetadataPostProcessor) UpdateRecord(datasetID string, modelID clientmodels.PennsieveSchemaID, recordUpdate clientmodels.RecordUpdate) error {
	_, err := p.Pennsieve.UpdateRecord(datasetID, modelID, recordUpdate.PennsieveID, recordUpdate.RecordValues)
	if err != nil {
		return err
	}
	return nil
}

func (p *MetadataPostProcessor) DeleteRecords(datasetID string, modelID clientmodels.PennsieveSchemaID, recordIDs []clientmodels.PennsieveInstanceID) error {
	if err := p.Pennsieve.DeleteRecords(datasetID, modelID, recordIDs); err != nil {
		return err
	}
	return nil
}
