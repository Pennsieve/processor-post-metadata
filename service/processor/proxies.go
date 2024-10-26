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
	modelID, err := p.IDStore.ModelID(proxyRecordChanges.ModelName)
	if err != nil {
		return fmt.Errorf("unable to delete package proxies for model %s: %w", proxyRecordChanges.ModelName, err)
	}
	targetRecordID, err := p.IDStore.RecordID(modelID, proxyRecordChanges.RecordExternalID)
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
