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

package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type REST struct {
	method   string
	endpoint string
	data     interface{}
}

func Load(config map[string]interface{}) *REST {
	result := REST{}

	if method, ok := config["method"].(string); ok {
		result.method = method
	}

	if endpoint, ok := config["endpoint"].(string); ok {
		result.endpoint = endpoint
	}

	if data, ok := config["data"]; ok {
		result.data = data
	}

	return &result
}

func (r *REST) Run(data []byte, attributes map[string]string) {
	var err error
	var resp *http.Response

	switch r.method {
	case "POST":
		var j []byte
		j, err = json.Marshal(r.data)
		if err != nil {
			log.Println("rest:", err)
			return
		}
		resp, err = http.Post(r.endpoint, "application/json", strings.NewReader(string(j)))
	default:
		resp, err = http.Get(r.endpoint)
	}
	if err != nil {
		log.Println("rest:", err)
	} else {
		log.Println("rest:", resp)
	}
}
