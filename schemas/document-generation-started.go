package schemas

//  'document-generation-started' topic struct
type DocumentGenerationStarted struct {
    RequesterID string `avro:"requester_id"     json:"requester_id"`
    ID          string `avro:"id"               json:"id"`
}
