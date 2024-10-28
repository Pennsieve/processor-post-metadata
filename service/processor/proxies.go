package processor

import (
	"fmt"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"log/slog"
)

func (p *MetadataPostProcessor) ProcessProxyInstanceDeletes(datasetID string, proxyChanges clientmodels.ProxyChanges) error {
	if len(proxyChanges.RecordChanges) == 0 {
		logger.Info("no proxy deletes")
		return nil
	}
	for _, proxyRecordChanges := range proxyChanges.RecordChanges {
		if err := p.ProcessProxyRecordChangesDeletes(datasetID, proxyRecordChanges); err != nil {
			return err
		}
	}
	return nil
}

func (p *MetadataPostProcessor) ProcessProxyRecordChangesDeletes(datasetID string, proxyRecordChanges clientmodels.ProxyRecordChanges) error {
	proxyLogger := logger.With(slog.Any("modelName", proxyRecordChanges.ModelName))
	if len(proxyRecordChanges.InstanceIDDeletes) == 0 {
		proxyLogger.Info("no proxy deletes")
		return nil
	}
	proxyLogger.Info("starting proxy deletes")
	targetRecordID, err := p.lookupTargetID(proxyRecordChanges.ModelName, proxyRecordChanges.RecordExternalID)
	if err != nil {
		return fmt.Errorf("unable to delete package proxies for model %s: %w", proxyRecordChanges.ModelName, err)
	}
	proxyLogger = proxyLogger.With(slog.Any("targetRecordID", targetRecordID))
	body := models.NewDeleteProxyInstancesBody(targetRecordID, proxyRecordChanges.InstanceIDDeletes...)
	if err := p.Pennsieve.DeleteProxyInstances(datasetID, body); err != nil {
		return fmt.Errorf("error deleting proxy instances for model %s record %s: %w",
			proxyRecordChanges.ModelName,
			targetRecordID,
			err)
	}
	proxyLogger.Info("finished proxy deletes", slog.Int("count", len(proxyRecordChanges.InstanceIDDeletes)))
	return nil
}

func (p *MetadataPostProcessor) ProcessProxyChanges(datasetID string, proxyChanges *clientmodels.ProxyChanges) error {
	if proxyChanges == nil {
		logger.Info("no proxy changes")
		return nil
	}
	logger.Info("starting proxy changes")
	if proxyChanges.CreateProxyRelationshipSchema {
		logger.Info("creating proxy relationship schema")
		_, err := p.Pennsieve.CreateProxyRelationshipSchema(datasetID)
		if err != nil {
			return fmt.Errorf("error creating proxy relationship schema: %w", err)
		}
	}
	if err := p.ProcessProxyRecordChanges(datasetID, proxyChanges.RecordChanges); err != nil {
		return err
	}
	logger.Info("finished proxy changes")
	return nil
}

func (p *MetadataPostProcessor) ProcessProxyRecordChanges(datasetID string, proxyRecordChanges []clientmodels.ProxyRecordChanges) error {
	if len(proxyRecordChanges) == 0 {
		logger.Info("no proxy changes")
		return nil
	}
	for _, changes := range proxyRecordChanges {
		if err := p.ProcessProxyInstanceCreates(datasetID, changes.ModelName, changes.RecordExternalID, changes.NodeIDCreates); err != nil {
			return err
		}
	}
	return nil
}

func (p *MetadataPostProcessor) ProcessProxyInstanceCreates(datasetID string, modelName string, recordExternalID clientmodels.ExternalInstanceID, packageNodeIDs []string) error {
	proxyLogger := logger.With(slog.Any("modelName", modelName))
	if len(packageNodeIDs) == 0 {
		proxyLogger.Info("no proxy creates")
		return nil
	}
	proxyLogger.Info("starting proxy creates")
	targetRecordID, err := p.lookupTargetID(modelName, recordExternalID)
	if err != nil {
		return fmt.Errorf("unable to create package proxies for model %s: %w", modelName, err)
	}
	proxyLogger = proxyLogger.With(slog.Any("targetRecordID", targetRecordID))
	for _, packageID := range packageNodeIDs {
		body := models.NewCreateProxyInstanceBody(targetRecordID, packageID)
		if err := p.Pennsieve.CreateProxyInstance(datasetID, body); err != nil {
			return fmt.Errorf("error creating proxy instance for model %s record %s package %s: %w",
				modelName,
				targetRecordID,
				packageID,
				err)
		}
	}
	proxyLogger.Info("finished proxy creates", slog.Int("count", len(packageNodeIDs)))
	return nil
}

func (p *MetadataPostProcessor) lookupTargetID(modelName string, targetExternalID clientmodels.ExternalInstanceID) (clientmodels.PennsieveInstanceID, error) {
	modelID, err := p.IDStore.ModelID(modelName)
	if err != nil {
		return "", err
	}
	return p.IDStore.RecordID(modelID, targetExternalID)
}
