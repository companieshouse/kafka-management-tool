package schemas

// EmailSend struct for email-send topic
type EmailSend struct {
	AppID        string `avro:"app_id"                          json:"app_id"`
	MessageID    string `avro:"message_id"                      json:"message_id"`
	MessageType  string `avro:"message_type"                    json:"message_type"`
	Data         string `avro:"data"                            json:"data"`
	EmailAddress string `avro:"email_address"                   json:"email_address"`
	CreatedAt    string `avro:"created_at"                      json:"created_at"`
}
