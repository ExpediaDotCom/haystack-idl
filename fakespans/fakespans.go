package main

import (
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"
	"./proto"
	"github.com/golang/protobuf/proto"
	"github.com/codeskyblue/go-uuid"
	"math/rand"
)

var (
	kafkaTopic = flag.String("topic", "spans", "Kafka Topic")
	kafkaBroker = flag.String("kafka-broker", "192.168.99.100:9092", "kafka TCP address for Span-Proto messages. e.g. localhost:9092")
	spanInterval = flag.Int("interval", 1, "period in seconds between spans")
	totalDuration = flag.Int("total-duration", 120, "total period worth of spans in seconds")
)

func main() {
	flag.Parse()

	if *kafkaBroker != "" {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
		config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
		err := config.Validate()
		if err != nil {
			log.Fatal(4, "invalid kafka producer config. %s", err)
		}

		client, err := sarama.NewSyncProducer([]string{*kafkaBroker}, config)
		if err != nil {
			log.Fatal(4, "failed to create kafka  producer for broker path ", *kafkaBroker, err)
		}
		produceSpansSync(client, *spanInterval, *totalDuration)
	}
}

func produceSpansSync(client sarama.SyncProducer, interval, totalDuration int) {

	spanCount := totalDuration / interval
	timestamp := time.Now().Unix()

	payload := make([]*sarama.ProducerMessage, spanCount)

	for index := 0; index < spanCount;index++ {
		timestamp += int64(interval)
		testSpan := generateSpan(timestamp, index)
		data, err := proto.Marshal(&testSpan)
		if err != nil {
			log.Fatal("marshaling error: ", err)
		}
		payload[index] = &sarama.ProducerMessage{
			Key:   sarama.StringEncoder(testSpan.TraceId),
			Topic: *kafkaTopic,
			Value: sarama.ByteEncoder(data),
		}
	}
	fmt.Print("pushing spans to kafka")
	client.SendMessages(payload)

}
func generateSpan(epochTimeInSecs int64, index int) span.Span {
	operationName := "some-span"
	host := "some-service"
	process := &span.Process{
		ServiceName: host,
	}
	return span.Span{
		TraceId: uuid.NewRandom().String(),
		SpanId: uuid.NewRandom().String(),
		ParentSpanId: uuid.NewRandom().String(),
		OperationName: operationName,
		StartTime: epochTimeInSecs * 1000,
		Duration: int64(rand.Int31()),
		Process: process,
	}
}
