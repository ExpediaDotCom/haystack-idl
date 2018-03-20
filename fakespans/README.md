## Creating test data in kafka 

fakespans is a simple go app which can generate random spans and push to kafka 

## Using fakespans

Run the following commands on your terminal to start using fake spans you should have golang installed on your box

1. export $GOPATH=`location where you want your go binaries`
2. export $GOBIN=$GOPATH/bin
3. cd fakespans
4. go install
5. $GOPATH/bin/fakespans


##fakespans options

`fakespans` utility can either generate random spans or upload spans from a file. If `from-file` option is specified it takes precedence over `trace-count` and `span-count` options. 

File should contain spans serialized in JSON format with one span per line. A [sample JSON](./test-spans.json) file is provided for reference. 

If `span-count and trace-count` options are used, the utility will generate `trace-count` number of traces with a total of `span-count` number of spans. 
 
```
./fake_metrics -h
Usage of fakespans:
  -interval int
        period in seconds between spans (default 1)
  -kafka-broker string
        kafka TCP address for Span-Proto messages. e.g. localhost:9092 (default "localhost:9092")
  -span-count int
        total number of unique spans you want to generate (default 120)
  -topic string
        Kafka Topic (default "spans")
  -trace-count int
        total number of unique traces you want to generate (default 20)
  -from-file File with Spans in JSON format. One span per line
  
```