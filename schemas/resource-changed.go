package schemas

// EventRecord - describes the event the resource was changed for
type EventRecord struct {
	PublishedAt   string   `json:"published_at" avro:"published_at"`
	Type          string   `json:"type" avro:"type"`
	FieldsChanged []string `json:"fields_changed" avro:"fields_changed"`
}

// ResourceChanged resource changed (json is for output in this utility, nothing to do with input to kafka api)
type ResourceChanged struct {
	ResourceKind string      `json:"resource_kind" avro:"resource_kind"`
	ResourceURI  string      `json:"resource_uri" avro:"resource_uri"`
	ContextID    string      `json:"context_id" avro:"context_id"`
	DeletedData  string      `json:"deleted_data" avro:"deleted_data"`
	Event        EventRecord `json:"event" avro:"event"`
}
