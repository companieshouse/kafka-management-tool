package main

import (
	"./schemas"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/companieshouse/chs.go/avro"
	"github.com/companieshouse/chs.go/avro/schema"
	"github.com/companieshouse/chs.go/kafka/producer"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

// Assigns all flags to variables
var (
	brokerPtr         = flag.String("broker", "", "Broker")
	topicPtr          = flag.String("topic", "", "Topic name")
	schemaPtr         = flag.String("schema", "", "Schema")
	schemaRegistryPtr = flag.String("schema-registry", "", "Schema Registry")
	partitionPtr      = flag.Int64("partition", 0, "Partition")
	zookeeperPtr      = flag.String("zookeeper", "", "Zookeeper")
	offsetPtr         = flag.String("offset", "-1", "Offset number")
	jsonOutPtr        = flag.Int("json-out", -1, "Does deserialized JSON get printed to terminal?")
)

func main() {
	flag.Parse()

	fmt.Println("broker:", *brokerPtr)
	fmt.Println("topic:", *topicPtr)
	fmt.Println("schema:", *schemaPtr)
	fmt.Println("schema-registry:", *schemaRegistryPtr)
	fmt.Println("partition:", *partitionPtr)
	fmt.Println("zookeeper-registry:", *zookeeperPtr)
	fmt.Println("offset:", *offsetPtr)
	fmt.Println("json-out:", *jsonOutPtr)

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

			schemaStruct := schemas.IdentifySchema(*topicPtr)
			var err error

			// Get schema from schema registry and use it to create avro consumer
			schema, err := schema.Get(*schemaRegistryPtr, *topicPtr)
			consumerAvro := &avro.Schema{
				Definition: schema,
			}

			// Unmarshal message value and assign it to schemaStruct
			if err = consumerAvro.Unmarshal(msg.Value, schemaStruct); err != nil {
				fmt.Println("Error unmarshalling avro:")
			}

			// create avro producer
			producerAvro := &avro.Schema{
				Definition: schema,
			}

			// Marshall schemaStruct that contains message value ready for republishing
			messageBytes, err := producerAvro.Marshal(schemaStruct)
			if err != nil {
				fmt.Println("Error Marshalling avro:")
			}

			// create producer message
			producerMessage := &producer.Message{
				Value: messageBytes,
				Topic: *topicPtr,
			}

			//create new producer
			p, err := producer.New(&producer.Config{Acks: &producer.WaitForAll, BrokerAddrs: []string{*brokerPtr}})
			if err != nil {
				os.Exit(1)
			}

			// republish message
			partition, offset, err := p.Send(producerMessage)
			if err != nil {
				fmt.Println("Error republishing message:")
			}
			fmt.Println("Message republished!", partition, offset)

			consumed++
		case <-signals:
			fmt.Println("break")
			break ConsumerLoop
		}
	}
}

// The offset slice is passed into method and is iterated over and each message
// related to the offset in the topic/partition is outputted into the 'out' chan
func consumePartition(consumer sarama.Consumer, topic string, partition int32, out chan *sarama.ConsumerMessage, offsetArray []int64) {
	for offset := range offsetArray {
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
