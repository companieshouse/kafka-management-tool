package schemas

// RenderSubmittedDataDocument struct for render-submitted-data-document topic
type RenderSubmittedDataDocument struct {
	ID           string `avro:"id"                              json:"id"`
	Resource     string `avro:"resource"                        json:"resource"`
	ResourceID   string `avro:"resource_id"                     json:"resource_id"`
	ContentType  string `avro:"content_type"                    json:"content_type"`
	DocumentType string `avro:"document_type"                   json:"document_type"`
	UserID       string `avro:"user_id"                         json:"user_id"`
}
