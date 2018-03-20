package main

import (
	"flag"
	"github.com/Shopify/sarama"
	"log"
	"time"
	"./proto"
	"github.com/golang/protobuf/proto"
	"github.com/codeskyblue/go-uuid"
	"math/rand"
	"fmt"
)

var (
	kafkaTopic = flag.String("topic", "proto-spans", "Kafka Topic")
	kafkaBroker = flag.String("kafka-broker", "192.168.99.100:9092", "kafka TCP address for Span-Proto messages")
	spanFile = flag.String("from-file", "", "File with Spans in JSON format. One span per line")
	spanInterval = flag.Int("interval", 1, "period in seconds between spans")
	traceCount = flag.Int("trace-count", 20, "total number of unique traces you want to generate")
	spanCount = flag.Int("span-count", 120, "total number of unique spans you want to generate")
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
		} else {
			//close the client after sending spans
			defer client.Close()

			if len(*spanFile) > 0 {
				produceSpansFromFile(&client, *spanFile, *spanInterval)
			} else {
				produceSpansSync(&client, *spanInterval, *spanCount, *traceCount)
			}
		}
	}
}
func produceSpansFromFile(producer *sarama.SyncProducer, fileName string, interval int) {
	line := 0
	ReadSpans(fileName, func(spanRecord SpanRecord) {
		line = line + 1
		fmt.Printf("%4v: %v\n", line, spanRecord)

		spanMessage := SpanFromSpanRecord(spanRecord)
		fmt.Printf("\n%4v: Transformed record\n %v\n", line, spanMessage)
		data, err := proto.Marshal(&spanMessage)
		if err != nil {
			log.Fatalf("Marshaling error [SpanId : %s] : %v\n", spanRecord.SpanId, err)
		} else {
			message := &sarama.ProducerMessage{
				Key:   sarama.StringEncoder(spanRecord.TraceId),
				Topic: *kafkaTopic,
				Value: sarama.ByteEncoder(data),
				Timestamp:time.Unix(atoi(spanRecord.StartTime, -1), 0),
			}

			_,_,err := (*producer).SendMessage(message)
			if err != nil {
				log.Fatalf( "Failed to produce data in kafka [SpanId: %s] with error %v\n",
					spanRecord.SpanId, err)
			} else {
				log.Printf("Successfully pushed span [%s] to kafka", spanRecord.SpanId)
			}
		}
	})
}

func produceSpansSync(clientPointer *sarama.SyncProducer, interval, spanCount int, traceCount int) {

	timestamp := time.Now().Unix() - int64(spanCount * interval)
	client := *clientPointer
	payload := make([]*sarama.ProducerMessage, spanCount)
	rootSpans := make([]*span.Span, traceCount)

	for index := 1; index <= spanCount; index++ {
		timestamp += int64(interval)
		traceIndex := index % traceCount
		var testSpan span.Span
		if rootSpans[traceIndex] != nil {
			rootSpan := rootSpans[traceIndex]
			testSpan = generateSpan(timestamp, rootSpan.TraceId, rootSpan.SpanId)
		} else {
			testSpan = generateSpan(timestamp, uuid.NewRandom().String(), "")
			rootSpans[traceIndex] = &testSpan
		}

		data, err := proto.Marshal(&testSpan)
		if err != nil {
			log.Fatal("marshaling error: ", err)
		}
		payload[index - 1] = &sarama.ProducerMessage{
			Key:   sarama.StringEncoder(testSpan.TraceId),
			Topic: *kafkaTopic,
			Value: sarama.ByteEncoder(data),
			Timestamp:time.Unix(timestamp, 0),
		}
	}
	log.Println("pushing spans to kafka")
	err := client.SendMessages(payload)
	if err != nil {
		log.Fatal(4, "failed to produce data in kafka with error", err)
	} else {
		log.Println("successfully pushed spans to kafka")
	}

}
func generateSpan(epochTimeInSecs int64, traceid string, parentid string) span.Span {
	operationName := "some-span"
	serviceName := "some-service"
	return span.Span{
		TraceId: traceid,
		SpanId: uuid.NewRandom().String(),
		ParentSpanId: parentid,
		OperationName: operationName,
		StartTime: epochTimeInSecs * 1000,
		Duration: int64(rand.Int31()),
		ServiceName: serviceName,
	}
}
