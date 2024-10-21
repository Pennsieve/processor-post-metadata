package processor

import clientmodels "github.com/pennsieve/processor-post-metadata/client/models"

type InstanceIDLookup map[clientmodels.ExternalInstanceID]clientmodels.PennsieveInstanceID
type RecordIDLookup InstanceIDLookup

// IDStore will hold maps to the Pennsieve IDs of the metadata objects the processor creates
type IDStore struct {
	ModelByName          map[string]clientmodels.PennsieveSchemaID
	RecordIDByExternalID InstanceIDLookup
}

func NewIDStore() *IDStore {
	return &IDStore{
		ModelByName:          make(map[string]clientmodels.PennsieveSchemaID),
		RecordIDByExternalID: make(InstanceIDLookup),
	}
}

func (s *IDStore) AddModel(name string, id clientmodels.PennsieveSchemaID) {
	s.ModelByName[name] = id
}

func (s *IDStore) AddRecord(externalID clientmodels.ExternalInstanceID, id clientmodels.PennsieveInstanceID) {
	s.RecordIDByExternalID[externalID] = id
}
