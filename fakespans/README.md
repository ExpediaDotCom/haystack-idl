
## Creating test data in Kafka 

`fakespans` is a simple app written in the Go language which can generate random spans and push them to Kafka. 

## Using fakespans

Run the following commands on your terminal to start using `fakespans`. You will need to have the Go language installed in order to run `fakespans`.

 ```shell
export GOPATH= location where you want your go binaries (should end in /bin)
export GOBIN=  your GOPATH + `/bin`
cd fakespans/
go get
go install
cd $GOBIN
./fakespans
```

## fakespans command line options

```
./fakespans -h
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
```

