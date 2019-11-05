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

package lambda

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type Lambda struct {
	function  string
	region    string
	data      interface{}
	qualifier *string
}

func Load(config map[string]interface{}) *Lambda {
	result := Lambda{}

	if function, ok := config["function"].(string); ok {
		result.function = function
	}

	if region, ok := config["region"].(string); ok {
		result.region = region
	}

	if data, ok := config["data"]; ok {
		result.data = data
	}

	if qualifier, ok := config["qualifier"].(string); ok {
		result.qualifier = &qualifier
	}

	return &result
}

func (l *Lambda) Run(data []byte, attributes map[string]string) {
	log.Println(l.data)
	j, err := json.Marshal(l.data)
	if err != nil {
		log.Println("lambda:", err)
		return
	}

	input := lambda.InvokeInput{
		FunctionName: aws.String(l.function),
		Payload:      j,
	}

	if l.qualifier != nil {
		input.Qualifier = l.qualifier
	}

	log.Println(input)

	svc := lambda.New(session.New(), &aws.Config{Region: aws.String(l.region)})
	_, err = svc.Invoke(&input)
	if err != nil {
		log.Println("lambda:", err)
	}
}
