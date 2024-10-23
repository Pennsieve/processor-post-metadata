package models

type Dataset struct {
	Models           []ModelChanges          `json:"models"`
	LinkedProperties []LinkedPropertyChanges `json:"linked_properties"`
	Proxies          *ProxyChanges           `json:"proxies"`
	RecordIDMaps     []RecordIDMap           `json:"record_id_maps"`
}

type RecordIDMap struct {
	ModelName           string                                     `json:"model_name"`
	ExternalToPennsieve map[ExternalInstanceID]PennsieveInstanceID `json:"external_to_pennsieve"`
}

func NewRecordIDMap(modelName string) RecordIDMap {
	return RecordIDMap{
		ModelName:           modelName,
		ExternalToPennsieve: make(map[ExternalInstanceID]PennsieveInstanceID),
	}
}
