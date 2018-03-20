# haystack-idl
Span and other data models used by Haystack are defined as [Protocol Buffer](https://developers.google.com/protocol-buffers/) files in [proto](./proto) folder

## Generating Java source for Haystack Spans
A simple maven pom file is available in [java](./java) folder to compile Haystack proto files in to a jar

## Creating test data in kafka 
Simple utility in Go to generate and send sample Spans to Kakfa is in [faksespans](./fakespans) folder


