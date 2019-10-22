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

package ec2

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2 struct {
	Method      string
	InstanceIds []string
	Region      string
}

func Load(config map[string]interface{}) *EC2 {
	result := EC2{}

	if method, ok := config["method"].(string); ok {
		result.Method = method
	}
	if instanceIds, ok := config["instanceIds"].([]interface{}); ok {
		for _, instanceId := range instanceIds {
			if instanceIdStr, ok := instanceId.(string); ok {
				result.InstanceIds = append(result.InstanceIds, instanceIdStr)
			}
		}
	}

	if region, ok := config["region"]; ok {
		result.Region = region.(string)
	} else {
		result.Region = "eu-central-1"
	}
	return &result
}

func (e *EC2) Run(data []byte, attributes map[string]string) {
	if len(e.InstanceIds) == 0 {
		// if list of instance ids is empty, don't try to start/stop instances
		return
	}

	switch e.Method {
	case "StartInstances":
		log.Println("ec2:StartInstances")

		// TODO start instance between 7am and 7pm only
		svc := ec2.New(session.New(), &aws.Config{Region: aws.String(e.Region)})
		_, err := svc.StartInstances(&ec2.StartInstancesInput{InstanceIds: aws.StringSlice(e.InstanceIds)})
		if err != nil {
			log.Println("Failed sending message:", err)
			return
		}
	}
}
