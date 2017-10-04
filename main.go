package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"./schemas"
	"github.com/Shopify/sarama"
	"github.com/companieshouse/chs.go/avro"
	"github.com/companieshouse/chs.go/avro/schema"
	"github.com/companieshouse/chs.go/kafka/producer"
	// "encoding/json"
	"encoding/json"
)

// Assigns all flags to variables
var (
	brokerPtr         = flag.String("broker", "", "Broker address")
	topicPtr          = flag.String("topic", "", "Topic name")
	schemaPtr         = flag.String("schema", "", "Schema name")
	schemaRegistryPtr = flag.String("schema-registry", "", "Schema Registry")
	partitionPtr      = flag.Int64("partition", 0, "Partition (default: 0)")
	offsetPtr         = flag.String("offset", "", "Offset number")
	jsonOutPtr        = flag.Int64("json-out", 0, "Print deserialized JSON message (default: 0)")
)

func main() {
	flag.Parse()

	// check if all mandatory flags have been provided
	validateFlags()

	// Print flags and parameters of the tool
	printFlags()

	// create offset array
	offsetArray := createOffsetArray(*offsetPtr)

	// create default config for sarama
	config := sarama.NewConfig()

	//create consumer
	consumer, err := sarama.NewConsumer([]string{*brokerPtr}, config)
	if err != nil {
		panic(err)
	}

	// create messages chan
	messages := make(chan *sarama.ConsumerMessage)
	go consumePartition(consumer, *topicPtr, int32(*partitionPtr), messages, offsetArray)

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

			//create new producer
			p, err := producer.New(&producer.Config{Acks: &producer.WaitForAll, BrokerAddrs: []string{*brokerPtr}})
			if err != nil {
				fmt.Errorf("error creating producer: %v", err)
				os.Exit(1)
			}

			// republish message
			fmt.Printf("Republishing message for offset: %v \n", msg.Offset)
			partition, offset, err := p.Send(producerMessage)
			if err != nil {
				fmt.Errorf("error republishing message with offset: %v, %v \n", msg.Offset, err)
				os.Exit(1)
			}
			fmt.Printf("Successfully republished message with offset %v to topic: %v using partition: %v new offset: %v \n", msg.Offset, *topicPtr, partition, offset)
			fmt.Println("---------------------")

			consumed++
		case <-signals:
			fmt.Println("break")
			break ConsumerLoop
		}
	}
}

func outputJSON(msg *sarama.ConsumerMessage, messageBytes chan []byte) {
	schemaStruct := schemas.IdentifySchema(*topicPtr)
	var err error

	// Get schema from schema registry and use it to create avro consumer
	fmt.Printf("Getting schema for topic: %v \n", *topicPtr)
	schema, err := schema.Get(*schemaRegistryPtr, *schemaPtr)
	if err != nil {
		fmt.Errorf("rror getting schema: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Retrieved schema for topic: %v \n", *topicPtr)

	consumerAvro := &avro.Schema{
		Definition: schema,
	}

	// Unmarshal message value and assign it to schemaStruct
	fmt.Printf("Unmarshalling message for offset: %v \n", msg.Offset)
	if err = consumerAvro.Unmarshal(msg.Value, schemaStruct); err != nil {
		fmt.Errorf("error unmarshalling avro message for offset: %v, %v \n", msg.Offset, err)
		os.Exit(1)
	}
	fmt.Printf("Message successfully unmarshalled for offset: %v \n", msg.Offset)

	// Print the unmarshalled message out to the terminal
	fmt.Println("Unmarshalled message:")
	data, err := json.Marshal(schemaStruct)
	if err != nil {
		fmt.Errorf("error marshalling JSON message for offset: %v, %v \n", msg.Offset, err)
		os.Exit(1)
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
		fmt.Errorf("error marshalling avro message for offset: %v, %v \n", msg.Offset, err)
		os.Exit(1)
	}
	fmt.Printf("Message successfully marshalled for offset: %v \n", msg.Offset)

	// output the marshalled message into the messageBytes chan
	messageBytes <- message
}

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
