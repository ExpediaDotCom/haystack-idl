/*
 *
 *  Copyright 2017 Expedia, Inc.
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 *
 */
 package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"./proto"
)

type Log struct {
	Timestamp string	
	Fields []map[string]string
}

type SpanRecord struct {
	TraceId string
	SpanId string
	ParentSpanId string
	ServiceName string
	OperationName string
	StartTime string
	Duration string
	Tags map[string]interface{}
	Logs []Log
}

func atoi(value string, defaultValue int64) int64 {
	i64, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		if defaultValue >= 0 {
			i64 = 0
		} else {
			log.Fatalf("Failed to convert [%s] to int64", value)
			panic(err)
		}
	}
	return i64
}

func ReadFile(fileName string, f func(reader io.Reader)) {
	fileHandle, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file [%s] : %s\n", fileName, err)
		panic(err)
	}
	defer fileHandle.Close()
	f(fileHandle)
}

func ReadSpans(fileName string, callBack func(SpanRecord)) {
	wrapper := func(reader io.Reader) {
		dec := json.NewDecoder(reader)
		for {
			m := SpanRecord{}
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			callBack(m)
		}
	}
	ReadFile(fileName, wrapper)
}


func SpanFromSpanRecord(spanRecord SpanRecord) span.Span {
	spanMessage := span.Span{
		TraceId: spanRecord.TraceId,
		SpanId: spanRecord.SpanId,
		ParentSpanId: spanRecord.ParentSpanId,
		ServiceName: spanRecord.ServiceName,
		OperationName: spanRecord.OperationName,
		StartTime: atoi(spanRecord.StartTime, -1),
		Duration: atoi(spanRecord.Duration, 0),
	}

	spanMessage.Tags = make([]*span.Tag, 0)
	for k, v := range spanRecord.Tags {
		tag := span.Tag {
			Key: k,
		}
		switch value := v.(type) {
		case int64:
			tag.Type = span.Tag_LONG
			tag.Myvalue = &span.Tag_VLong{value}
		case string:
			tag.Type = span.Tag_STRING
			tag.Myvalue = &span.Tag_VStr{value}
		case bool:
			tag.Type = span.Tag_BOOL
			tag.Myvalue = &span.Tag_VBool{value}
		}
		spanMessage.Tags = append(spanMessage.Tags, &tag)
	}

	spanMessage.Logs = make([]*span.Log, 0)
	for _, v := range spanRecord.Logs {
		logEntry := span.Log{
			Timestamp: atoi(v.Timestamp, -1),
		}
		logEntry.Fields = make([]*span.Tag, 0)
		for _,logFields := range v.Fields {
			eventName, ok := logFields["vStr"]
			if ok {
				logEntry.Fields = append(logEntry.Fields, &span.Tag{
					Key: "event",
					Type: span.Tag_STRING,
					Myvalue: &span.Tag_VStr{eventName},
				})
			}
		}

		spanMessage.Logs = append(spanMessage.Logs, &logEntry)
	}

	return spanMessage
}


