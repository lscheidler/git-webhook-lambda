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

package rules

import (
	"encoding/json"
	"log"

	"github.com/lscheidler/git-webhook-lambda/rules/hooks"
	"github.com/lscheidler/git-webhook-lambda/rules/hooks/ec2"
	"github.com/lscheidler/git-webhook-lambda/rules/hooks/lambda"
	"github.com/lscheidler/git-webhook-lambda/rules/hooks/rest"
	"github.com/lscheidler/git-webhook-lambda/rules/hooks/sqs"
)

type Rules map[string]Type

func (r *Rules) UnmarshalJSON(b []byte) error {
	var types map[string]Type
	if err := json.Unmarshal(b, &types); err != nil {
		return err
	}
	*r = Rules(types)
	return nil
}

type Type map[string][]Hook

func (t *Type) UnmarshalJSON(b []byte) error {
	var hooks map[string][]Hook
	if err := json.Unmarshal(b, &hooks); err != nil {
		return err
	}
	*t = Type(hooks)
	return nil
}

type Hook struct {
	Plugin hooks.Hook
}

func (h *Hook) UnmarshalJSON(b []byte) error {
	var settings map[string]interface{}
	if err := json.Unmarshal(b, &settings); err != nil {
		return err
	}
	var plugin hooks.Hook
	switch settings["type"] {
	case "ec2":
		plugin = ec2.Load(settings)
	case "lambda":
		plugin = lambda.Load(settings)
	case "rest":
		plugin = rest.Load(settings)
	case "sqs":
		plugin = sqs.Load(settings)
	}
	*h = Hook{Plugin: plugin}
	return nil
}

func Load(rules string) *Rules {
	var result Rules
	err := json.Unmarshal([]byte(rules), &result)
	if err != nil {
		log.Println("err:", err)
		return nil
	}
	return &result
}
