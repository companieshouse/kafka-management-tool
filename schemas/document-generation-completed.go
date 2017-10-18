package schemas

// DocumentGenerationCompleted struct for document-generation-completed topic
type DocumentGenerationCompleted struct {
	RequesterID           string            `avro:"requester_id"                    json:"requester_id"`
	ID                    string            `avro:"id"                              json:"id"`
	Description           string            `avro:"description"                     json:"description"`
	DescriptionIdentifier string            `avro:"description_identifier"          json:"description_identifier"`
	DocumentCreatedAt     string            `avro:"document_created_at"             json:"document_created_at"`
	DescriptionValues     DescriptionValues `avro:"description_values"              json:"description_values"`
	Location              string            `avro:"location"                        json:"location"`
	DocumentSize          string            `avro:"document_size"                   json:"document_size"`
}

//  DescriptionValues struct for document-generation-completed topic
type DescriptionValues struct {
	Date string `avro:"date"                            json:"date"`
}
