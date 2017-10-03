package schemas

//  This method takes the topic name as a parameter
//  and determines the struct to return as a schema.
func IdentifySchema(topicName string) interface{} {
    switch topicName {
    case "tx-closed":
        return &TxClosed{}
    case "document-generation-started":
        return &DocumentGenerationStarted{}
    case "document-generation-completed":
        return &DocumentGenerationCompleted{}
    }
    return nil
}