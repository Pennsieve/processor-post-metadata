package models

// LinkedPropertyChanges represents a changeset for a particular LinkedProperty, optionally creating the link
// in the metadata schema.
// Executing the changeset will depend on other models and records existing, so it should only be executed once
// those model changesets have been executed.
type LinkedPropertyChanges struct {
	// The ID of the LinkedProperty in the schema. Can be empty or missing if the linkedProperty does not exist.
	// In this case, Create below should be non-nil
	ID string `json:"id,omitempty"`

	// If Create is non-nil, the link schema should be created in the model schema
	Create *SchemaLinkedPropertyCreate `json:"create,omitempty"`

	// Instances contains the create and delete info for link instances
	Instances InstanceChanges `json:"instances"`
}

// SchemaLinkedPropertyCreate will have to be translated to a POST /models/datasets/<dataset id>/concepts/<from model id>/linked/bulk
// request with body [{"name": Name,"displayName": DisplayName,"to": <to model id>,"position": Position}] once to and from model id
// values are known. (These models may not exist when instances of this struct are created by clients.)
type SchemaLinkedPropertyCreate struct {
	// Name is the name of the linked property in the schema of the parent model
	Name string `json:"name"`

	// DisplayName is the display name of the linked property in the schema of the parent model
	DisplayName string `json:"display_name"`

	// FromModelName is the name of the parent model of the linked proeprty, i.e., the "from" model
	FromModelName string `json:"from_model_name"`

	// ToModelName is the name of the model of the child model of the linked property, i.e., the "to" model
	ToModelName string `json:"to_model_name"`

	// Position is the position of the linked property in the schema of the parent model. (?)
	Position int `json:"position"`
}

type InstanceChanges struct {
	Create []InstanceLinkedPropertyCreate `json:"create"`
	Delete []InstanceLinkedPropertyDelete `json:"delete"`
}

// InstanceLinkedPropertyCreate will have to be translated to a POST /models/datasets/<dataset id>/concepts/<model id>/instances/<from record id>/linked
// request with body {"schemaLinkedPropertyId": <linked property schema id>, "to": <to record id>} once those id values are known
// (The LinkedProperty schema and/or from and to records may not yet exist when instances of this struct are created
type InstanceLinkedPropertyCreate struct {
	FromExternalID ExternalInstanceID `json:"from_external_id"`
	ToExternalID   ExternalInstanceID `json:"to_external_id"`
}

// InstanceLinkedPropertyDelete contains the additional info needed to call
// DELETE /models/datasets/{dataset id}/concepts/{model Id}/instances/{record id}/linked/{link instance id}
// No payload required
type InstanceLinkedPropertyDelete struct {
	FromRecordID             PennsieveInstanceID `json:"from_record_id"`
	InstanceLinkedPropertyID PennsieveInstanceID `json:"instance_linked_property_id"`
}
