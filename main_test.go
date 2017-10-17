package main

import (
	"flag"
	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/companieshouse/chs.go/kafka/producer"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUnitProcessMessages(t *testing.T) {
	offsetArraySingle := []int64{1}

	Convey("test successful - message published", t, func() {
		consumerMock := mocks.NewConsumer(t, nil)
		consumerGroup := consumerMock.ExpectConsumePartition("", 0, 1)
		consumerGroup.YieldMessage(&sarama.ConsumerMessage{})

		producerMock := mocks.NewSyncProducer(t, nil)
		producerMock.ExpectSendMessageAndSucceed()
		argu := Arguments{
			OffsetArray: offsetArraySingle,
			Consumer:    consumerMock,
			Producer:    &producer.Producer{producerMock},
		}
		processMessages(argu)
	})
}

func TestUnitValidateFlags(t *testing.T) {
	flag.Set("broker", "broker")
	flag.Set("schema", "schema")
	flag.Set("offset", "1")
	flag.Set("schema-registry", "schema-registry")
	flag.Set("topic", "topic")
	//flag.Set("test.bench", "test.bench")
	//flag.Set("test.blockprofile", "test.blockprofile")
	//flag.Set("test.coverprofile", "test.coverprofile")
	//flag.Set("test.cpu", "test.cpu")
	//flag.Set("test.memprofile", "test.memprofile")
	//flag.Set("test.cpuprofile", "test.cpuprofile")
	//flag.Set("test.mutexprofile", "test.mutexprofile")
	//flag.Set("test.outputdir", "test.outputdir")
	//flag.Set("test.trace", "test.trace")
	//flag.Set("test.run", "test.run")

	Convey("test successful", t, func() {
		So(validateFlags(), ShouldBeNil)
	})
}

func TestUnitCreateOffsetArray(t *testing.T) {
	arraySingle := make([]int64, 0)
	arraySingle = append(arraySingle, 10)

	arrayRange := make([]int64, 0)
	arrayRange = append(arrayRange, 10)
	arrayRange = append(arrayRange, 11)
	arrayRange = append(arrayRange, 12)
	arrayRange = append(arrayRange, 13)
	arrayRange = append(arrayRange, 14)
	arrayRange = append(arrayRange, 15)

	Convey("test successful", t, func() {
		So(createOffsetArray("10"), ShouldResemble, arraySingle)
		So(createOffsetArray("10-15"), ShouldResemble, arrayRange)
	})
}
