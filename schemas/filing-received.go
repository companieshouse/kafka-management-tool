package schemas

// FilingReceived struct for filing-received topic
type FilingReceived struct {
	ApplicationID string                   `avro:"application_id"                  json:"application_id"`
	Attempt       int32                    `avro:"attempt"                         json:"attempt"`
	ChannelID     string                   `avro:"channel_id"                      json:"channel_id"`
	Presenter     PresenterFilingReceived  `avro:"presenter"                       json:"presenter"`
	Submission    SubmissionFilingReceived `avro:"submission"                      json:"submission"`
	Items         []Items                  `avro:"items"                           json:"items"`
}

// PresenterFilingReceived struct represents PresenterFilingReceived in the FilingReceived struct
type PresenterFilingReceived struct {
	Forename string `avro:"forename"                        json:"forename"`
	Language string `avro:"language"                        json:"language"`
	Surname  string `avro:"surname"                         json:"surname"`
	UserID   string `avro:"user_id"                         json:"user_id"`
}

// SubmissionFilingReceived struct represents SubmissionFilingReceived in the FilingReceived struct
type SubmissionFilingReceived struct {
	CompanyNumber string `avro:"company_number"                  json:"company_number"`
	CompanyName   string `avro:"company_name"                    json:"company_name"`
	ReceivedAt    string `avro:"received_at"                     json:"received_at"`
	TransactionID string `avro:"transaction_id"                  json:"transaction_id"`
}

// Items struct represents Items in the FilingReceived struct
type Items struct {
	Data               string `avro:"data"                            json:"data"`
	Kind               string `avro:"kind"                            json:"kind"`
	SubmissionLanguage string `avro:"submission_language"             json:"submission_language"`
	SubmissionID       string `avro:"submission_id"                   json:"submission_id"`
}
