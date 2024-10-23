package processor

import (
	"fmt"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-post-metadata/service/models"
	"log/slog"
)

func (p *MetadataPostProcessor) ProcessLinks(datasetID string, linkChanges []clientmodels.LinkedPropertyChanges) error {
	if len(linkChanges) == 0 {
		logger.Info("no link changes")
		return nil
	}
	logger.Info("starting link changes")
	for _, linkChange := range linkChanges {
		if err := p.ProcessLinkChanges(datasetID, linkChange); err != nil {
			return err
		}
	}
	logger.Info("finished link changes")
	return nil
}

func (p *MetadataPostProcessor) ProcessLinkChanges(datasetID string, linkChange clientmodels.LinkedPropertyChanges) error {
	schemaIDs, err := p.CreateLinkIfNecessary(datasetID, linkChange)
	if err != nil {
		return err
	}
	linkLogger := logger.With(slog.Any("linkSchemaID", schemaIDs.Link))

	linkLogger.Info("creating link instances")
	for _, instanceCreate := range linkChange.Instances.Create {
		if err := p.CreateLinkInstance(datasetID, schemaIDs, instanceCreate); err != nil {
			return err
		}
	}
	linkLogger.Info("created link instances", slog.Int("count", len(linkChange.Instances.Create)))

	return nil
}

func (p *MetadataPostProcessor) CreateLinkInstance(datasetID string, schemaIDs SchemaID, instanceCreate clientmodels.InstanceLinkedPropertyCreate) error {
	fromRecordID, err := p.IDStore.RecordID(schemaIDs.FromModel, instanceCreate.FromExternalID)
	if err != nil {
		return fmt.Errorf("'from' record id not found: %w", err)
	}
	toRecordID, err := p.IDStore.RecordID(schemaIDs.ToModel, instanceCreate.ToExternalID)
	if err != nil {
		return fmt.Errorf("'to' record id not found: %w", err)
	}
	body := models.CreateLinkInstanceBody{
		SchemaLinkedPropertyId: schemaIDs.Link,
		To:                     toRecordID,
	}
	if err := p.Pennsieve.CreateLinkedPropertyInstance(datasetID, schemaIDs.FromModel, fromRecordID, body); err != nil {
		return fmt.Errorf("error creating linked property instance: %w", err)
	}
	return nil
}

type SchemaID struct {
	// Link is the schema id of the linked property
	Link      clientmodels.PennsieveSchemaID
	FromModel clientmodels.PennsieveSchemaID
	ToModel   clientmodels.PennsieveSchemaID
}

func (p *MetadataPostProcessor) CreateLinkIfNecessary(datasetID string, linkChange clientmodels.LinkedPropertyChanges) (SchemaID, error) {
	fromModelID, foundFrom := p.IDStore.ModelByName[linkChange.FromModelName]
	if !foundFrom {
		return SchemaID{}, fmt.Errorf("from model id for name %s not found", linkChange.FromModelName)
	}
	toModelID, toFound := p.IDStore.ModelByName[linkChange.ToModelName]
	if !toFound {
		return SchemaID{}, fmt.Errorf("to model id for name %s not found", linkChange.ToModelName)
	}
	if linkChange.Create == nil {
		linkSchemaID := clientmodels.PennsieveSchemaID(linkChange.ID)
		logger.Info("linked property already exists", slog.Any("linkSchemaID", linkSchemaID))
		return SchemaID{
			FromModel: fromModelID,
			ToModel:   toModelID,
			Link:      linkSchemaID,
		}, nil
	}

	linkCreate := linkChange.Create
	linkLogger := logger.With(slog.String("linkName", linkCreate.Name))
	linkLogger.Info("creating link schema")

	body := models.CreateLinkSchemaBody{
		Name:        linkCreate.Name,
		DisplayName: linkCreate.DisplayName,
		To:          toModelID,
		Position:    linkCreate.Position,
	}
	linkID, err := p.Pennsieve.CreateLinkedPropertySchema(datasetID, fromModelID, body)
	if err != nil {
		return SchemaID{}, fmt.Errorf("error creating link schema: %w", err)
	}
	linkLogger.Info("link schema created", slog.Any("linkID", linkID))
	return SchemaID{
		FromModel: fromModelID,
		Link:      linkID,
		ToModel:   toModelID,
	}, nil
}
