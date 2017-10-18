package schemas

// DocumentGenerationStarted struct for document-generation-started topic
type DocumentGenerationStarted struct {
	RequesterID string `avro:"requester_id"                    json:"requester_id"`
	ID          string `avro:"id"                              json:"id"`
}
