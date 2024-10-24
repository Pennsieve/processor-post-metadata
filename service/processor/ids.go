package processor

import (
	"fmt"
	clientmodels "github.com/pennsieve/processor-post-metadata/client/models"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
)

type RecordIDKey struct {
	ModelID    clientmodels.PennsieveSchemaID
	ExternalID clientmodels.ExternalInstanceID
}

type RecordIDLookup map[RecordIDKey]clientmodels.PennsieveInstanceID

// IDStore will hold maps to the Pennsieve IDs of metadata objects
type IDStore struct {
	ModelByName   map[string]clientmodels.PennsieveSchemaID
	RecordIDbyKey RecordIDLookup
}

func (s *IDStore) AddModel(name string, id clientmodels.PennsieveSchemaID) {
	s.ModelByName[name] = id
}

func (s *IDStore) AddRecord(modelID clientmodels.PennsieveSchemaID, externalID clientmodels.ExternalInstanceID, id clientmodels.PennsieveInstanceID) {
	s.RecordIDbyKey[RecordIDKey{
		ModelID:    modelID,
		ExternalID: externalID,
	}] = id
}

func (s *IDStore) AddRecordIDMaps(recordIDMaps []clientmodels.RecordIDMap) error {
	for _, recordIDMap := range recordIDMaps {
		modelID, err := s.ModelID(recordIDMap.ModelName)
		if err != nil {
			return fmt.Errorf("unable to add recordIDMap for model: %s: %w", recordIDMap.ModelName, err)
		}
		for external, pennsieve := range recordIDMap.ExternalToPennsieve {
			s.AddRecord(modelID, external, pennsieve)
		}

	}
	return nil
}

func (s *IDStore) RecordID(modelID clientmodels.PennsieveSchemaID, externalID clientmodels.ExternalInstanceID) (clientmodels.PennsieveInstanceID, error) {
	recordID, found := s.RecordIDbyKey[RecordIDKey{
		ModelID:    modelID,
		ExternalID: externalID,
	}]
	if !found {
		return "", fmt.Errorf("no record in model %s with external id %s", modelID, externalID)
	}
	return recordID, nil
}

func (s *IDStore) ModelID(modelName string) (clientmodels.PennsieveSchemaID, error) {
	modelID, found := s.ModelByName[modelName]
	if !found {
		return "", fmt.Errorf("id for model %s not found", modelName)
	}
	return modelID, nil
}

type IDStoreBuilder struct {
	store *IDStore
}

func NewIDStoreBuilder() *IDStoreBuilder {
	return &IDStoreBuilder{store: &IDStore{
		ModelByName:   make(map[string]clientmodels.PennsieveSchemaID),
		RecordIDbyKey: make(RecordIDLookup),
	}}
}

func (b *IDStoreBuilder) WithSchema(schema *metadataclient.Schema) *IDStoreBuilder {
	for name, id := range schema.ModelIDsByName() {
		b.store.ModelByName[name] = clientmodels.PennsieveSchemaID(id)
	}
	return b
}

func (b *IDStoreBuilder) WithModel(name string, id clientmodels.PennsieveSchemaID) *IDStoreBuilder {
	b.store.AddModel(name, id)
	return b
}

func (b *IDStoreBuilder) WithRecord(modelID clientmodels.PennsieveSchemaID, externalID clientmodels.ExternalInstanceID, recordID clientmodels.PennsieveInstanceID) *IDStoreBuilder {
	b.store.AddRecord(modelID, externalID, recordID)
	return b
}

func (b *IDStoreBuilder) Build() *IDStore {
	return b.store
}
