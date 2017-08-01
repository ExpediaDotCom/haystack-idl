# haystack-idl
Span and other data model definitions used by Haystack

## Creating test data in kafka (for Mac)
If homebrew is installed, and you've brought up your local environment with the mechanism provided by the 
[haystack-deployment](https://github.com/ExpediaDotCom/haystack-deployment) package, then you can issue the 
following from a command line prompt:
1. ```brew install protobuf``` to put protoc in /usr/local/bin.
2. ```brew install kafkacat``` to put kafkacat in /usr/local/bin.
3. ```cat ../samples/span.decoded | protoc --encode=Span span.proto | kafkacat -P -b $(minikube ip):9092 -t test```
to push the contents of the ```span.decoded``` file (in the /proto directory that is a **child** of the 
directory in which this README.md is found) into the local Kafka.
