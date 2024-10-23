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

func NewIDStore(existingSchema *metadataclient.Schema) *IDStore {
	existingModelMap := existingSchema.ModelIDsByName()
	byName := make(map[string]clientmodels.PennsieveSchemaID, len(existingModelMap))

	for name, id := range existingModelMap {
		pennsieveID := clientmodels.PennsieveSchemaID(id)
		byName[name] = pennsieveID
	}
	return &IDStore{
		ModelByName:   byName,
		RecordIDbyKey: make(RecordIDLookup),
	}
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
