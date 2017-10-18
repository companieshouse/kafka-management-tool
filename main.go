package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/companieshouse/chs.go/avro"
	"github.com/companieshouse/chs.go/avro/schema"
	"github.com/companieshouse/chs.go/kafka/producer"
	"github.com/companieshouse/kafka-management-tool/schemas"
)

// Assigns all flags to variables
var (
	brokerPtr         = flag.String("broker", "", "Broker address")
	topicPtr          = flag.String("topic", "", "Topic name")
	schemaPtr         = flag.String("schema", "", "Schema name")
	schemaRegistryPtr = flag.String("schema-registry", "", "Schema Registry")
	partitionPtr      = flag.Int64("partition", 0, "Partition (default: 0)")
	offsetPtr         = flag.String("offset", "", "Offset number, can be single offset number i.e. 10 or a range separated by the - operator i.e. 10-20")
	jsonOutPtr        = flag.Int64("json-out", 0, "Print deserialized JSON message (default: 0)")
)

// Arguments struct will hold the offset, consumer and producer ready for
// passing into the main process method, this allows setup to be done
// before any processing
type Arguments struct {
	OffsetArray []int64
	Consumer    sarama.Consumer
	Producer    *producer.Producer
}

func main() {
	flag.Parse()

	// check if all mandatory flags have been provided
	if err := validateFlags(); err != nil {
		panic(err)
	}

	// Print flags and parameters of the tool
	printFlags()

	//create consumer
	consumer, err := sarama.NewConsumer([]string{*brokerPtr}, nil)
	if err != nil {
		fmt.Println("error creating consumer")
		panic(err)
	}

	//create new producer
	p, err := producer.New(&producer.Config{Acks: &producer.WaitForAll, BrokerAddrs: []string{*brokerPtr}})
	if err != nil {
		fmt.Println("error creating producer")
		panic(err)
	}

	// pass arguments to struct
	arguments := Arguments{
		OffsetArray: createOffsetArray(*offsetPtr),
		Consumer:    consumer,
		Producer:    p,
	}

	// process the messages on the topic
	processMessages(arguments)
}

// processMessages iterates through the topic messages and
// republishes them and outputting them if the json-out flag
// is set to true
func processMessages(argu Arguments) {

	// create messages chan
	messages := make(chan *sarama.ConsumerMessage)
	go consumePartition(argu.Consumer, *topicPtr, int32(*partitionPtr), messages, argu.OffsetArray)

	// create signals chan
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-messages:

			// make messageBytes channel
			messageBytes := make(chan []byte)

			// if jsonOutPtr is set to 1, then output JSON and republish
			// otherwise just republish message
			if *jsonOutPtr == 1 {
				go outputJSON(msg, messageBytes)
			} else {
				go func() {
					messageBytes <- msg.Value
				}()
			}

			// create producer message
			producerMessage := &producer.Message{
				Value: <-messageBytes,
				Topic: *topicPtr,
			}

			// republish message
			fmt.Printf("Republishing message for offset: %v \n", msg.Offset)
			partition, offset, err := argu.Producer.Send(producerMessage)
			if err != nil {
				fmt.Println("error republishing message with offset: %d, \n", msg.Offset)
				panic(err)
			}
			fmt.Printf("Successfully republished message with offset %v to topic: %v using partition: %v new offset: %v \n", msg.Offset, *topicPtr, partition, offset)
			fmt.Println("---------------------")

			consumed++

			// check if we have hit the end of the array of messages
			// if so, exit the tool by breaking the ConsumerLoop
			if consumed == len(argu.OffsetArray) {
				break ConsumerLoop
			}
		case <-signals:
			fmt.Println("break")
			break ConsumerLoop
		}
	}
}

// outputJSON prints the deserialised json message to the console for
// easier reading if the json-out flag is set to 1
func outputJSON(msg *sarama.ConsumerMessage, messageBytes chan []byte) {
	schemaStruct := schemas.IdentifySchema(*topicPtr)

	if schemaStruct == nil {
		fmt.Println("no schema identified")
		os.Exit(1)
	}

	var err error

	// Get schema from schema registry and use it to create avro consumer
	fmt.Printf("Getting schema for topic: %v \n", *topicPtr)
	schema, err := schema.Get(*schemaRegistryPtr, *schemaPtr)
	if err != nil {
		fmt.Println("error getting schema")
		panic(err)
	}
	fmt.Printf("Retrieved schema for topic: %v \n", *topicPtr)

	consumerAvro := &avro.Schema{
		Definition: schema,
	}

	// Unmarshal message value and assign it to schemaStruct
	fmt.Printf("Unmarshalling message for offset: %v \n", msg.Offset)
	if err = consumerAvro.Unmarshal(msg.Value, schemaStruct); err != nil {
		fmt.Printf("error unmarshalling avro message for offset: %v \n", msg.Offset)
		panic(err)
	}
	fmt.Printf("Message successfully unmarshalled for offset: %v \n", msg.Offset)

	// Print the unmarshalled message out to the terminal
	fmt.Println("Unmarshalled message:")
	data, err := json.Marshal(schemaStruct)
	if err != nil {
		fmt.Printf("error marshalling JSON message for offset: %v \n", msg.Offset)
		panic(err)
	}
	fmt.Println(string(data))

	// create avro producer
	producerAvro := &avro.Schema{
		Definition: schema,
	}

	// Marshall schemaStruct that contains message value ready for republishing
	fmt.Printf("Marshalling message for offset: %v \n", msg.Offset)
	message, err := producerAvro.Marshal(schemaStruct)
	if err != nil {
		fmt.Printf("error marshalling avro message for offset: %v \n", msg.Offset)
		panic(err)
	}
	fmt.Printf("Message successfully marshalled for offset: %v \n", msg.Offset)

	// output the marshalled message into the messageBytes chan
	messageBytes <- message
}

// Validates the flags to make sure all mandatory flags have been supplied
// throw an error if not the case
func validateFlags() error {
	//args := os.Args
	var error int8
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			fmt.Printf("Value not supplied for: %v \n", f.Name)
			error = 1
		}
	})

	if error == 1 {
		return errors.New("flag validation failed")
	}
	return nil
}

// checks if s is within the e string array.
// this will be used to check if the flag matches the ones passed into
// program. This is for the unit tests as the flags package add extra flags
func contains(s string, e []string) bool {
	for _, a := range e {
		if a == s {
			return true
		}
	}
	return false
}

// printFlags prints the flags and parameters of the tool
// to the console
func printFlags() {
	fmt.Println("---------------------")
	fmt.Println("Parameters:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%v: %v\n", f.Name, f.Value)
	})
	fmt.Println("---------------------")
}

// consumePartition iterates over and each message in the topic related to the offset number
// in the offsetArray and outputs the message into the 'out' chan
func consumePartition(consumer sarama.Consumer, topic string, partition int32, out chan *sarama.ConsumerMessage, offsetArray []int64) {
	for _, offset := range offsetArray {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
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

// createOffsetArray creates the offset array depending on whether a
// range or a single value is entered as an argument to the tool
func createOffsetArray(offset string) []int64 {
	arraySize := make([]int64, 0)
	if strings.ContainsAny(offset, "-") {
		slice := strings.Split(offset, "-")
		minRange, err := strconv.ParseInt(slice[0], 10, 64)
		if err != nil {
			panic(err)
		}

		maxRange, err := strconv.ParseInt(slice[1], 10, 64)
		if err != nil {
			panic(err)
		}

		if minRange > maxRange {
			fmt.Printf("min range cannot be greater than max range")
			os.Exit(1)
		}

		for value := minRange; value <= maxRange; value++ {
			arraySize = append(arraySize, value)
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
