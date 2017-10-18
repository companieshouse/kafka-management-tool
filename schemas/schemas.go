package schemas

// IdentifySchema takes the topic name as a parameter
// and determines the struct to return as a schema.
func IdentifySchema(topicName string) interface{} {
	switch topicName {
	case "tx-closed":
		return &TxClosed{}
	case "document-generation-started":
		return &DocumentGenerationStarted{}
	case "document-generation-completed":
		return &DocumentGenerationCompleted{}
	case "email-send":
		return &EmailSend{}
	case "render-submitted-data-document":
		return &RenderSubmittedDataDocument{}
	case "filing-received":
		return &FilingReceived{}
	case "filing-processed":
		return &FilingProcessed{}
	}
	return nil
}
