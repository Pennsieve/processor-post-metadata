package processor

import (
	"fmt"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"log/slog"
)

func (p *MetadataPostProcessor) ProcessRecordModelDeletes(datasetID string, modelUpdates []clientmodels.ModelUpdate, modelDeletes []clientmodels.ModelDelete) error {
	for _, modelChange := range modelUpdates {
		if err := p.ProcessRecordDeletes(datasetID, modelChange.ID, modelChange.Records.Delete); err != nil {
			return err
		}
	}
	for _, modelChange := range modelDeletes {
		if err := p.ProcessRecordDeletes(datasetID, modelChange.ID, modelChange.Records); err != nil {
			return err
		}
		// Now that records are deleted, we can delete the model
		if err := p.ProcessModelDelete(datasetID, modelChange.ID); err != nil {
			return err
		}
	}
	return nil
}

func (p *MetadataPostProcessor) ProcessRecordDeletes(datasetID string, modelID clientmodels.PennsieveSchemaID, recordIDs []clientmodels.PennsieveInstanceID) error {
	modelLogger := logger.With(slog.Any("modelID", modelID))
	if len(recordIDs) == 0 {
		modelLogger.Info("no record deletes")
		return nil
	}
	modelLogger.Info("starting record deletes")
	if err := p.Pennsieve.DeleteRecords(datasetID, modelID, recordIDs); err != nil {
		return err
	}
	modelLogger.Info("finished record deletes", slog.Int("count", len(recordIDs)))
	return nil
}

func (p *MetadataPostProcessor) ProcessModelDelete(datasetID string, modelID clientmodels.PennsieveSchemaID) error {
	modelLogger := logger.With(slog.Any("modelID", modelID))
	modelLogger.Info("deleting model")
	if err := p.Pennsieve.DeleteModel(datasetID, modelID); err != nil {
		return err
	}
	modelLogger.Info("deleted model")
	return nil
}

func (p *MetadataPostProcessor) ProcessModelCreatesUpdates(datasetID string, modelCreates []clientmodels.ModelCreate, modelUpdates []clientmodels.ModelUpdate) error {
	if len(modelCreates)+len(modelUpdates) == 0 {
		logger.Info("no model changes")
		return nil
	}
	logger.Info("starting model changes")
	for _, modelChange := range modelCreates {
		if err := p.ProcessModelCreate(datasetID, modelChange); err != nil {
			return err
		}
	}
	for _, modelChange := range modelUpdates {
		if err := p.ProcessModelUpdate(datasetID, modelChange); err != nil {
			return err
		}
	}
	logger.Info("finished model changes")
	return nil
}

func (p *MetadataPostProcessor) ProcessModelCreate(datasetID string, modelCreate clientmodels.ModelCreate) error {
	modelID, err := p.CreateModel(datasetID, modelCreate.Create)
	if err != nil {
		return err
	}
	modelLogger := logger.With(slog.Any("modelID", modelID))
	modelLogger.Info("creating records")
	for _, recordCreate := range modelCreate.Records {
		if err := p.CreateRecord(datasetID, modelID, recordCreate); err != nil {
			return err
		}
	}
	modelLogger.Info("created records", slog.Int("count", len(modelCreate.Records)))
	return nil
}

func (p *MetadataPostProcessor) ProcessModelUpdate(datasetID string, modelUpdate clientmodels.ModelUpdate) error {
	modelID := modelUpdate.ID
	modelLogger := logger.With(slog.Any("modelID", modelID))
	modelLogger.Info("creating records")
	for _, recordCreate := range modelUpdate.Records.Create {
		if err := p.CreateRecord(datasetID, modelID, recordCreate); err != nil {
			return err
		}
	}
	modelLogger.Info("created records", slog.Int("count", len(modelUpdate.Records.Create)))
	modelLogger.Info("updating records")
	for _, recordUpdate := range modelUpdate.Records.Update {
		if err := p.UpdateRecord(datasetID, modelID, recordUpdate); err != nil {
			return err
		}
	}
	modelLogger.Info("updated models", slog.Int("count", len(modelUpdate.Records.Update)))

	return nil
}

func (p *MetadataPostProcessor) CreateModel(datasetID string, modelCreate clientmodels.ModelPropsCreate) (clientmodels.PennsieveSchemaID, error) {
	modelLogger := logger.With(slog.String("modelName", modelCreate.Model.Name))
	modelLogger.Info("creating model")
	modelID, err := p.Pennsieve.CreateModelAndProps(datasetID, modelCreate)
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
