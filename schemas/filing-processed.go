package schemas

// FilingProcessed struct for filing-processed topic
type FilingProcessed struct {
	ApplicationID string     `avro:"application_id"                  json:"application_id"`
	ChannelID     string     `avro:"channel_id"                      json:"channel_id"`
	Presenter     Presenter  `avro:"presenter"                       json:"presenter"`
	Submission    Submission `avro:"submission"                      json:"submission"`
	Response      Response   `avro:"response"                        json:"response"`
	Attempt       int32      `avro:"attempt"                         json:"-"`
}

// Presenter struct represents Presenter within FilingProcessed struct
type Presenter struct {
	Language string `avro:"language"                        json:"language"`
	UserID   string `avro:"user_id"                         json:"user_id"`
}

// Submission represents Submission within filing-processed struct
type Submission struct {
	TransactionID string `avro:"transaction_id"                  json:"transaction_id"`
}

// Response represents Response within FilingProcessed struct
type Response struct {
	CompanyName    string `avro:"company_name"                    json:"company_name,omitempty"`
	CompanyNumber  string `avro:"company_number"                  json:"company_number,omitempty"`
	DateOfCreation string `avro:"date_of_creation"                json:"date_of_creation,omitempty"`
	Status         string `avro:"status"                          json:"status"`
	SubmissionID   string `avro:"submission_id"                   json:"submission_id"`
	Reject         Reject `avro:"reject"                          json:"reject,omitempty"`
	ProcessedAt    string `avro:"processed_at"                    json:"processed_at"`
}

// Reject represents Reject within Response struct
type Reject struct {
	ReasonsEnglish []string `avro:"reasons_english"                 json:"reasons_english"`
	ReasonsWelsh   []string `avro:"reasons_welsh"                   json:"reasons_welsh"`
}
