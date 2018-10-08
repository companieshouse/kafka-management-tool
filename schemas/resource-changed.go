package schemas

// ResourceChanged resource changed
type ResourceChanged struct {
	ResourceKind string `json:"resource_kind" avro:"resource_kind"`
	ChangeType   string `json:"change_type" avro:"change_type"`
	PublishedAt  string `json:"published_at" avro:"published_at"`
	ResourceURI  string `json:"resource_uri" avro:"resource_uri"`
}
