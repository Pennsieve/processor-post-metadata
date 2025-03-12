package models

import (
	"encoding/json"
	"fmt"
)

// PennsieveSchemaID is the internal ID of a model or other schema object in Pennsieve.
// Not usually seen by user, but needed for API calls
type PennsieveSchemaID string

func (i PennsieveSchemaID) String() string {
	return string(i)
}

// PennsieveInstanceID is the internal ID of the record or other instance in Pennsieve.
// Not usually seen by user, but needed for API calls
type PennsieveInstanceID string

// ExternalInstanceID is an ID that the user supplies, or we can calculate from the values a user supplies. It identifies
// an instance, so that two instances with the same ID should be considered the same.
type ExternalInstanceID string

type ModelChanges struct {
	Creates []ModelCreate `json:"creates"`
	Updates []ModelUpdate `json:"updates"`
	Deletes []ModelDelete `json:"deletes"`
}

type ModelCreate struct {
	// Create contains the params necessary to create the model and its (non-link) properties
	Create ModelPropsCreate `json:"create"`
	// Records are records that should be created
	Records []RecordCreate `json:"records"`
}

type ModelUpdate struct {
	// The ID of the model in Pennsieve.
	ID PennsieveSchemaID `json:"id"`
	// Records describes the changes to the records of this model type
	Records RecordChanges `json:"records"`
}

type ModelDelete struct {
	// The ID of the model in Pennsieve.
	ID PennsieveSchemaID `json:"id"`
	// A list of RecordIDs to delete. All records must be deleted in order to delete a model
	Records []PennsieveInstanceID `json:"records"`
}

type ModelPropsCreate struct {
	Model      ModelCreateParams      `json:"model"`
	Properties PropertiesCreateParams `json:"properties"`
}

// ModelCreateParams can be used as a payload for POST /models/datasets/<dataset id>/concepts to create a model
type ModelCreateParams struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Locked      bool   `json:"locked"`
}

// PropertiesCreateParams can be uses as a payload for PUT /models/datasets/<dataset id>/concepts/<model id>/properties to add properties to a model
type PropertiesCreateParams []PropertyCreateParams

type PropertyCreateParams struct {
	DisplayName  string          `json:"displayName"`
	Name         string          `json:"name"`
	DataType     json.RawMessage `json:"dataType"`
	ConceptTitle bool            `json:"conceptTitle"`
	Default      bool            `json:"default"`
	Required     bool            `json:"required"`
	IsEnum       bool            `json:"isEnum"`
	IsMultiValue bool            `json:"isMultiValue"`
	Value        string          `json:"value"`
	Locked       bool            `json:"locked"`
	Description  string          `json:"description"`
}

func (pc *PropertyCreateParams) SetDataType(dataType any) error {
	bytes, err := json.Marshal(dataType)
	if err != nil {
		return fmt.Errorf("error marshalling data type: %w", err)
	}
	pc.DataType = bytes
	return nil
}

type RecordChanges struct {
	// A list of RecordIDs to delete
	Delete []PennsieveInstanceID `json:"delete"`
	// Create are records that should be created
	Create []RecordCreate `json:"create"`
	// Update are records that should be updated
	Update []RecordUpdate `json:"update"`
}

// RecordCreate wraps a RecordValues that can be used as a payload for
// POST /models/datasets/<dataset id>/concepts/<model id>/instances to create a new record.
// The ExternalID is not part of the payload, but is a non-pennsieve identifier for the record that
// can be used to map this record to the eventual PennsieveInstanceID needed for links or package proxies
type RecordCreate struct {
	ExternalID ExternalInstanceID `json:"external_id"`
	RecordValues
}

type RecordValue struct {
	Value any    `json:"value"`
	Name  string `json:"name"`
}

type RecordValues struct {
	Values []RecordValue `json:"values"`
}

// RecordUpdate wraps a RecordValues that can be used as a payload for PUT /models/datasets/<dataset id>/concepts/<model id>/instances/<record id> to update values in record
// Include both changed and unchanged values
// The PennsieveID is not part of the payload, but is the record id needed as a request path parameter
type RecordUpdate struct {
	PennsieveID PennsieveInstanceID `json:"pennsieve_id"`
	RecordValues
}
