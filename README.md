# kafka-management-tool

## Introduction
The kafka-management-tool is a command line `Go` tool that has been created for the sole purpose of being able to republish kafka messages in a given topic. Additionally it can also print the deserialised JSON into the terminal window for the user to read.

## Prerequisites 
To use this tool you will need to:
- Have `Go` installed

## Getting started
To build the service, you must first git clone it into your `$GOPATH` under `src/github.com/companieshouse`, and then run `make build` to build an executable in the current directory. To run the tool you can either do the following:
1) Perform `make build` in the `kafka-management-tool` root to create an executable. Then, assuming you are in the root, you can just run `kafka-management-tool [insert flag arguements here]` and it will run the tool.   
Example:  
`kafka-management-tool -broker=kafka-broker:9092 -topic=topic-one -schema=schema-one  -schema-registry=http://kafka-registry:8081 -partition=0 -offset=693 -json-out=0`
2) Perform `make install` and it will create an executable in the `$GOPATH/bin` directory where you can then execute the tool from anywhere using the following command `$GOPATH/bin/kafka-management-tool [insert flag arguements here]` and it will run the tool.  
Example:  
`$GOPATH/bin/kafka-management-tool -broker=kafka-broker:9092 -topic=topic-one -schema=schema-one  -schema-registry=http://kafka-registry:8081 -partition=0 -offset=693 -json-out=0`

## In-depth
As described in the introduction, this tool will print and republish specific messages into the specified Kafka topics. Given that the mandatory flags are entered into the tool, the tool will go through a number of steps in order to print and republish the chosen message. Below is the basic flow of the tool.

![alt text](https://user-images.githubusercontent.com/29541485/31217447-e5c36546-a9ae-11e7-94b8-89f38f59b273.png)

### Walkthrough
1) One of the first things the tool does is validate the flags that are passed in. If it sees that one or more mandatory flag isn't passed in. It will throw and error explaining this (the mandatory flags can be found below). Once validation is passed, the tool assumes that the information you have passed into the flags is correct. It then prints these flags back out into the terminal for reference sake.

```go
// Validates the flags to make sure all mandatory flags have been supplied
// throw an error if not the case
func validateFlags() {
	flag.VisitAll(func(f *flag.Flag) {
		if string(f.Value.String()) == "" {
			fmt.Printf("Value not supplied for: %v \n", f.Name)
			os.Exit(1)
		}
	})
}

// Prints the flags and parameters of the tool
func printFlags() {
	fmt.Println("---------------------")
	fmt.Println("Parameters:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%v: %v\n", f.Name, f.Value)
	})
	fmt.Println("---------------------")
}
```

2) The second step is where the tool determines the offset(s). There can be two types of offset passed into the tool. A single offset i.e. `20` or a range of offsets i.e. `20-100`. As long as you include the minus (`-`) between the `from` and `to` offsets then the tool will know that it's a range and it will add those and the numbers between, to the offset array. Note: if you have specified just one offset, it will still add this to the array.

```go
// This creates the offset array depending on whether a range or a single value
// is entered as an argument to the tool
func createOffsetArray(offset string) []int64 {
	arraySize := make([]int64, 0)
	if strings.ContainsAny(offset, "-") {
		slice := strings.Split(offset, "-")
		minRange, err := strconv.ParseInt(slice[0], 10, 64)
		maxRange, err := strconv.ParseInt(slice[1], 10, 64)
		if err != nil {
			panic(err)
		}

		index := 0
		for value := minRange; value <= maxRange; value++ {
			arraySize = append(arraySize, value)
			index++
		}
	} else {
		value, err := strconv.ParseInt(offset, 10, 64)
		if err != nil {
			panic(err)
		}
		arraySize = append(arraySize, value)
	}
	return arraySize
}
```

3) This is the step where the tool uses the offset array to retrieve the messages from the specified topic and partition. It does this by iterating over each offset within the offsetArray and retrieving the message based on that offset. Once the messages has been retrieved, assuming nothing has gone wrong in the retrieval. It will then add this message into the message channel.

```go
// The offset slice is passed into method and is iterated over and each message
// related to the offset in the topic/partition is outputted into the 'out' chan
func consumePartition(consumer sarama.Consumer, topic string, partition int32, out chan *sarama.ConsumerMessage, offsetArray []int64) {
	for _, offset := range offsetArray {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, int64(offset))
		if err != nil {
			panic(err)
		}

		select {
		case msg := <-partitionConsumer.Messages():
			partitionConsumer.Close()
			out <- msg
		}
	}
}
```
At this point the tool iterates over all retrieved messages added into the message channel and does the following steps to each message.

4) Here, is where the tool can go one of two ways. If the `json-out` flag has been set to true `1`, it will unmarshall each message and print the raw JSON into the terminal, if however the `json-out` flag has been set to `0` or isn't specified at all, it will skip the JSON printing logic and go straight to step 5.

 - a)  - Assuming that the `json-out` flag has been set to true (`1`), the next step is that it will identify and create the avro struct and get the schema from the regsistry. The avro struct will be used to hold the mapped unmarshalled json message. The schema will be retrieved from the registry by using the values that are passed into the `schema-registry` and `schema` flags.

 - b)  - Once the avro struct has been created and the schema has been retrieved, it will now use the avro consumer to unmarshall the message into the avro struct. At this point the values are mapped to the correct keys.

 - c)  -With the unmarshalled message mapped and held in the avro struct, it now marshalls that data into raw JSON using `json.Marshal` then printing the JSON to the terminal

```go
	// Print the unmarshalled message out to the terminal
	fmt.Println("Unmarshalled message:")
	data, err := json.Marshal(schemaStruct)
	if err != nil {
		fmt.Errorf("error marshalling JSON message for offset: %v, %v \n", msg.Offset, err)
		os.Exit(1)
	}
	fmt.Println(string(data))
```

 - d)  - Once the message has been printed to the terminal, the avro struct is marshalled back into the same format it was in initially

5) At this point, it creates the producer and the producerMessage (which contains the message) ready for republishing.

6) Depending on whether the JSON was printed or not, it shouldn't matter at this point. If the `json-out` flag was false, it goes from step 4 to 5 and then 6 and has no requirement to unmarshal the message (no need to unmarshal the message as we aren't printing the JSON). This means that the message that was originally retrieve can just get republished straight away. If the `json-out` flag was true, then it just republishes the message that was returned from step 4d.

Once all steps have completed and the message loop has ended. The tool then finishes and waits for the user to interrupt or kill it. (CTRL + C).

## Flags
The tool does require certain parameters for it to work. These will be provided in the form of `flags`.

| Flag                          | Description                         | Mandatory    | Example                 | Default  |
| ----------------------------- | ----------------------------------- | ------------ | ----------------------- | -------- |
| `broker`                      | Broker address                      | Yes          | `kafka-broker:9092`        |          |
| `topic`                       | Topic name                          | Yes          | `topic-one`             |          |
| `schema`                      | Schema name                         | Yes          | `schema-one`             |          |
| `schema-registry`             | Schema registry                     | Yes          | `http://kafka-registry:8081` |          |
| `partition`                   | Partition                           | No           | 0                       | 0        |
| `offset`                      | Schema name                         | Yes          | Single: 20, Range 10-20 |          |
| `json-out`                    | Print deserialized JSON message     | No           | 1                       | 0        |

### Example
This is an example usage of the tool along with its output.

With no JSON outputt:  
```
$ kafka-management-tool -broker=kafka-broker:9092 -topic=topic-one -schema=schema-one  -schema-registry=http://kafka-registry:8081 -partition=0 -offset=693 -json-out=0
---------------------
Parameters:
broker: kafka-broker:9092
json-out: 0
offset: 693
partition: 0
schema: schema-one
schema-registry: http://kafka-registry:8081
topic: topic-one
---------------------
Republishing message for offset: 693
Successfully republished message with offset 693 to topic: topic-one using partition: 0 new offset: 834
---------------------
^Cbreak
```

With JSON output:  
```
$ kafka-management-tool -broker=kafka-broker:9092 -topic=topic-one -schema=schema-one  -schema-registry=http://kafka-registry:8081 -partition=0 -offset=693 -json-out=1
---------------------
Parameters:
broker: kafka-broker:9092
json-out: 1
offset: 693
partition: 0
schema: schema-one
schema-registry: http://kafka-registry:8081
topic: topic-one
---------------------
Getting schema for topic: topic-one
Retrieved schema for topic: topic-one
Unmarshalling message for offset: 693
Message successfully unmarshalled for offset: 693
Unmarshalled message:
{"key":"value","key_two":"value_two","key_three":"value_three","key_four":"value_four"}
Marshalling message for offset: 693
Message successfully marshalled for offset: 693
Republishing message for offset: 693
Successfully republished message with offset 693 to topic: topic-one using partition: 0 new offset: 836
---------------------
```
