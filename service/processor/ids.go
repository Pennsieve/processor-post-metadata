package processor

import clientmodels "github.com/pennsieve/processor-post-metadata/client/models"

type IDLookup map[clientmodels.ExternalInstanceID]clientmodels.PennsieveInstanceID

type RecordIDLookup IDLookup

// IDStore will hold maps to the Pennsieve IDs of the metadata objects the processor creates
type IDStore struct {
	ModelByName map[string]clientmodels.PennsieveInstanceID
}

func NewIDStore() *IDStore {
	return &IDStore{
		ModelByName: make(map[string]clientmodels.PennsieveInstanceID),
	}
}

func (s *IDStore) AddModel(name string, id clientmodels.PennsieveInstanceID) {
	s.ModelByName[name] = id
}
