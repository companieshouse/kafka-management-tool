# Maintenance

## Introduction
At some point in time there will be an increase an topics or even a change to the existing topic structure. If at any point one of these happen, this will require a change to the Avro schemas in further requiring a change in the structs that represent the Avro schemas. These structs are used by the kafka-management-tool in order to map the unmarshalled JSON message from with the specified topic.

## Walkthrough
If a change is made to an existing topic structure/schema or if a new topic and Avro schema has been added this will need to be reflected in the kafka-management-tool schema struct for the corresponding Avro schema.

### Addition of new topic and Avro schema
If a new topic is to be added, there will need to be a new Avro schema file that represents the message structure. Additionally, the services that produce and consume messages from the new topic will need to create the mapping object/struct that represents the topic schema. In `Java` services this is a `Pojo` and in `Go`, a `struct`. For information on how the struct will can be created for the kafka-management-tool look below.
Additionally, once the new struct has been created, another condition will need to be added to the `schemas.go` file to make sure that the correct struct is used when the new topic is specified in the tool.

### Changes to an existing topic and Avro schema
If a change is to be made to an existing topic and its message structure, this means that the Avro schema file (.avsc) will be changed to reflect those changes. A hypothetical example of this could be the `tx-closed` topic. If, there was a decision to add a new field called `description` into the topic message structure, within the [tx-closed](https://github.com/companieshouse/chs-kafka-schemas/blob/master/schemas/tx-closed.avsc) Avro schema file, there will need to be an additional field added to the `fields` array. Whenever this topic is then used in any implementations, the Java Pojo model and/or the Go struct that represents the Avro schema will need to be updated. The same applies to the kafka-management-tool. The `TxClosed` schema will need to be amended to include the `description` field.

## Example
Here, some examples on how to amend the kafka-management-tool schema structs will be provided for both new topics and changes to existing ones.

### Addition of new topic and Avro schema
If a new topic called `tx-open` was added with the Avro schema `tx-open.avsc`, with the following structure:
```
{
  \"type\": \"record\",
  \"name\": \"transaction_open\",
  \"namespace\": \"namespaceName\",
  \"fields\": [
    {
      \"name\": \"timedate\",
      \"type\": \"string\"
    },
    {
      \"name\": \"transaction_details\",
      \"type\": \"string\"
    }
  ]
}
```

The `Go` struct will look like the following

```go
type TxClosed struct {
	Timedate                string  `avro:"timedate"                            json:"timedate"`
	TransactionDescription  string  `avro:"transaction_details"                 json:"transaction_details"`
}
```

A basic rule of thumb to go by is:
- struct property = "name" of the field within `.avsc` file in camelcase with no underscore to separate words
- property type  = "type" stated in the `avsc` file
- avro/json name = will be the same as the name in the `avsc` file
