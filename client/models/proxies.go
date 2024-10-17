package models

// ProxyChanges contains all the changes to package proxies.
// If CreateProxyRelationshipSchema is true we send POST /models/datasets/<dataset id>/relationships the body
// {"name":"belongs_to","displayName":"Belongs To","description":"","schema":[]} to create the special
// relationship schema used by all package proxies regardless of record or pacakge type.
type ProxyChanges struct {
	CreateProxyRelationshipSchema bool                 `json:"create_proxy_relationship_schema"`
	RecordChanges                 []ProxyRecordChanges `json:"record_changes"`
}

func (pc ProxyChanges) Summary() (createCount int, deleteCount int) {
	for _, recordChanges := range pc.RecordChanges {
		createCount += len(recordChanges.NodeIDCreates)
		deleteCount += len(recordChanges.InstanceIDDeletes)
	}
	return
}

// ProxyRecordChanges holds the changes to the package proxies of one Pennsieve record
// Creates will use POST /models/datasets/<dataset id>/proxy/package/instances with body
//
//	{
//	 "externalId": <package node id>,
//	 "targets": [
//	   {
//	     "direction": "FromTarget",
//	     "linkTarget": {
//	       "ConceptInstance": {
//	         "id": <record id>
//	       }
//	     },
//	     "relationshipType": "belongs_to",
//	     "relationshipData": []
//	   }
//	 ]
//	}
//
// Deletes will use DELETE /models/datasets/<dataset id>/proxy/package/instances/bulk with body
//
//	{
//	   "sourceRecordId": <record id>,
//	   "proxyInstanceIds": [
//	       <proxy id>,
//	       <proxy id>
//	   ]
//	}
type ProxyRecordChanges struct {
	// ModelName is the name of the model for this record
	// Together ModelName and RecordExternalID can be used to look up the ID we really need, the
	// Record's ID in Pennsieve
	ModelName string `json:"model_name"`
	// RecordExternalID is the external ID of the record. Will have to be used along with ModelName to look up the Pennsieve
	// instance ID when executing the changes
	RecordExternalID ExternalInstanceID `json:"record_external_id"`
	// NodeIDCreates The package node ids that should be linked to this record
	NodeIDCreates []string `json:"node_id_creates"`
	// InstanceIDDeletes the proxy instance ids to delete for this record
	InstanceIDDeletes []PennsieveInstanceID `json:"instance_id_deletes"`
}
