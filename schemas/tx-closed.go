package schemas

//  TxClosed struct for tx-closed topic
type TxClosed struct {
	Attempt                         int32                           `avro:"attempt"                         json:"attempt"`
	TransactionURL                  string                          `avro:"transaction_url"                 json:"transaction_url"`
}