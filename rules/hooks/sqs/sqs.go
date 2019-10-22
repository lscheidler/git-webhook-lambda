/*
Copyright 2019 Lars Eric Scheidler

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sqs

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQS struct {
	Method            string
	QueueUrl          string
	Region            string
	MessageAttributes map[string]string
}

func Load(config map[string]interface{}) *SQS {
	method, _ := config["method"].(string)
	queueUrl, _ := config["queueUrl"].(string)
	result := SQS{Method: method, QueueUrl: queueUrl}

	if region, ok := config["region"]; ok {
		result.Region = region.(string)
	} else {
		result.Region = "eu-central-1"
	}
	return &result
}

func (s *SQS) Run(data []byte, attributes map[string]string) {
	switch s.Method {
	case "SendMessage":
		log.Println("sqs:SendMessage")

		messageAttributes := make(map[string]*sqs.MessageAttributeValue)
		for key, value := range attributes {
			messageAttributes[key] = (&sqs.MessageAttributeValue{DataType: aws.String("String")}).SetStringValue(value)
		}

		svc := sqs.New(session.New(), &aws.Config{Region: aws.String(s.Region)})
		_, err := svc.SendMessage(&sqs.SendMessageInput{MessageBody: aws.String(string(data)), QueueUrl: aws.String(s.QueueUrl), MessageAttributes: messageAttributes})
		if err != nil {
			log.Println("Failed sending message:", err)
			return
		}
	}
}
