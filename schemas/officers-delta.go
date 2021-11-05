package schemas

// ChsDelta struct for officers-delta topic
type ChsDelta struct {
	Data      string `avro:"data"                            json:"data"`
	Attempt   int32  `avro:"attempt"                         json:"attempt"`
	ContextId string `avro:"context_id"                      json:"context_id"`
}
