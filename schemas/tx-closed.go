package schemas

//  'tx-closed' topic struct
type TxClosed struct {
	Attempt                     int32                   `avro:"attempt"                            json:"attempt"`
	TransactionURL              string                  `avro:"transaction_url"                    json:"transaction_url"`
}
