package schemas

//  DocumentGenerationFailed struct for document-generation-failed topic
type DocumentGenerationFailed struct {
	RequesterID           string            `avro:"requester_id"                    json:"requester_id"`
	Description           string            `avro:"description"                     json:"description"`
	DescriptionIdentifier string            `avro:"description_identifier"          json:"description_identifier"`
	ID                    string            `avro:"id"                              json:"id"`
	DescriptionValues     DescriptionValues `avro:"description_values"              json:"description_values"`
}
