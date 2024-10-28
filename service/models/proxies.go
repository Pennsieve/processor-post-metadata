package models

import clientmodels "github.com/pennsieve/processor-post-metadata/client/models"

const ProxyRelationshipSchemaName = "belongs_to"

type CreateProxyRelationshipSchemaBody struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Description string   `json:"description"`
	Schema      []string `json:"schema"`
}

func NewCreateProxyRelationshipSchemaBody() CreateProxyRelationshipSchemaBody {
	return CreateProxyRelationshipSchemaBody{
		Name:        ProxyRelationshipSchemaName,
		DisplayName: "Belongs To",
		Schema:      make([]string, 0),
	}
}

type CreateProxyInstanceBody struct {
	// ExternalID is the package node id
	ExternalID string                      `json:"externalId"`
	Targets    []CreateProxyInstanceTarget `json:"targets"`
}

func NewCreateProxyInstanceBody(recordID clientmodels.PennsieveInstanceID, packageNodeID string) CreateProxyInstanceBody {
	recordLinkTarget := linkTarget{ConceptInstance: conceptInstance{ID: recordID}}
	return CreateProxyInstanceBody{
		ExternalID: packageNodeID,
		Targets: []CreateProxyInstanceTarget{
			{
				Direction:        "FromTarget",
				LinkTarget:       recordLinkTarget,
				RelationshipType: ProxyRelationshipSchemaName,
				RelationshipData: make([]any, 0),
			},
		},
	}
}

type CreateProxyInstanceTarget struct {
	Direction        string     `json:"direction"`
	LinkTarget       linkTarget `json:"linkTarget"`
	RelationshipType string     `json:"relationshipType"`
	RelationshipData []any      `json:"relationshipData"`
}

type linkTarget struct {
	ConceptInstance conceptInstance `json:"ConceptInstance"`
}

type conceptInstance struct {
	// ID is the Pennsieve record id of the record to which the package is being linked
	ID clientmodels.PennsieveInstanceID `json:"id"`
}

type DeleteProxyInstancesBody struct {
	SourceRecordID   clientmodels.PennsieveInstanceID   `json:"sourceRecordId"`
	ProxyInstanceIDs []clientmodels.PennsieveInstanceID `json:"proxyInstanceIds"`
}

func NewDeleteProxyInstancesBody(recordID clientmodels.PennsieveInstanceID, proxyIDs ...clientmodels.PennsieveInstanceID) DeleteProxyInstancesBody {
	return DeleteProxyInstancesBody{
		SourceRecordID:   recordID,
		ProxyInstanceIDs: proxyIDs,
	}
}
