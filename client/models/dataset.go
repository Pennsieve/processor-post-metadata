package models

type Dataset struct {
	Models           []ModelChanges          `json:"models"`
	LinkedProperties []LinkedPropertyChanges `json:"linked_properties"`
	Proxies          *ProxyChanges           `json:"proxies"`
}
