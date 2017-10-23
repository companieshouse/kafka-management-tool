package main

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/companieshouse/chs.go/kafka/producer"
	. "github.com/smartystreets/goconvey/convey"
)

var lvalue int8

// TestUnitProcessMessages tests the processMessages function in main.go
// that is the main method for the processing of the messages
func TestUnitProcessMessages(t *testing.T) {
	offsetArraySingle := []int64{1}

	Convey("test successful - message published", t, func() {
		consumerMock := mocks.NewConsumer(t, nil)
		consumerGroup := consumerMock.ExpectConsumePartition("", 0, 1)
		consumerGroup.YieldMessage(&sarama.ConsumerMessage{})

		producerMock := mocks.NewSyncProducer(t, nil)
		producerMock.ExpectSendMessageAndSucceed()

		producerMock.ExpectSendMessageWithCheckerFunctionAndSucceed(increment())
		argu := Arguments{
			OffsetArray: offsetArraySingle,
			Consumer:    consumerMock,
			Producer:    &producer.Producer{producerMock},
		}
		processMessages(argu)
		So(lvalue, ShouldEqual, 1)
	})
}

func increment() func([]byte) error {
	return func(val []byte) error {
		lvalue = lvalue + 1
		return nil
	}
}

// TestUnitCreateFlagMap tests the createFlagMap function in main.go
// it checks that it doesn't return an empty map
func TestUnitCreateFlagMap(t *testing.T) {
	Convey("test successful - flags validated", t, func() {
		So(createFlagMap(), ShouldNotBeEmpty)
	})
}

// TestUnitValidateFlags tests the validateFlags function in main.go
// that is the validation method for the flags
func TestUnitValidateFlags(t *testing.T) {
	flagsMap := make(map[string]string)
	flagsMap["broker"] = "broker"
	flagsMap["topic"] = "topic"
	flagsMap["schema"] = "schema"
	flagsMap["schema-registry"] = "schema-registry"
	flagsMap["partition"] = "partition"
	flagsMap["offset"] = "offset"
	flagsMap["json-out"] = "json-out"

	Convey("test successful - flags validated", t, func() {
		So(validateFlags(flagsMap), ShouldBeNil)
	})
}

// TestUnitValidateFlagsWithInvalidMap tests the validateFlags function
// that is the validation method for the flags with an invalid map
func TestUnitValidateFlagsWithInvalidMap(t *testing.T) {
	flagsMap := make(map[string]string)
	flagsMap["offset"] = ""

	Convey("test successful - flags validated", t, func() {
		So(validateFlags(flagsMap), ShouldNotBeNil)
	})
}

// TestUnitCreateOffsetArray tests the createOffsetArray function in main.go
// that create the offset array from the offset passed into the tool as a param arg
func TestUnitCreateOffsetArray(t *testing.T) {
	arraySingle := []int64{10}

	arrayRange := []int64{10, 11, 12, 13, 14, 15}

	Convey("test successful - offsetArray created", t, func() {
		So(createOffsetArray("10"), ShouldResemble, arraySingle)
		So(createOffsetArray("10-15"), ShouldResemble, arrayRange)
	})
}
