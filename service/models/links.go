package models

import clientmodels "github.com/pennsieve/processor-post-metadata/client/models"

type CreateLinkSchemaBody struct {
	Name        string                         `json:"name"`
	DisplayName string                         `json:"displayName"`
	To          clientmodels.PennsieveSchemaID `json:"to"`
	Position    int                            `json:"position"`
}

type CreateLinkInstanceBody struct {
	SchemaLinkedPropertyId clientmodels.PennsieveSchemaID   `json:"schemaLinkedPropertyId"`
	To                     clientmodels.PennsieveInstanceID `json:"to"`
}
