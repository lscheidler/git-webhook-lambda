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

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/lscheidler/git-webhook-lambda/rules"
	"github.com/lscheidler/git-webhook-plugin-bitbucket"
)

func main() {
	local := flag.Bool("local", false, "run lambda function localy")
	flag.Parse()

	if *local {
		httpServer()
	} else {
		lambda.Start(Handler)
	}
}

func httpServer() {
	http.HandleFunc("/", HttpHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Handler(ctx context.Context, request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	//fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	//fmt.Printf("Body size = %d.\n", len(request.Body))
	//fmt.Printf("Body %s.\n", request.Body)

	if os.Getenv("debug") == "true" {
		log.Println("Headers:")
		for key, value := range request.Headers {
			log.Printf("    %s: %s\n", key, value)
		}
		log.Println("MultiValueHeaders:")
		for key, value := range request.MultiValueHeaders {
			log.Printf("    %s: %v\n", key, value)
		}

		log.Println("request.RequestContext:", request.RequestContext)
		log.Println("ctx:", ctx)
		log.Println("os.Environ:", os.Environ())
		log.Println("body:", request.Body)
	}

	rules := rules.Load(os.Getenv("rules"))

	bitbkt := bitbucket.Init()
	if bitbkt.IsBitbucketWebhookRequest(request.Headers) {
		if bitbkt.Valid(request) {
			log.Println("received bitbucket webhook", request.Headers["x-event-key"])
			// rules
			for _, hook := range (*rules)["bitbucket"][request.Headers["x-event-key"]] {
				hook.Plugin.Run(bitbkt.Event().Data(), bitbkt.Attributes())
			}
		}
	}

	resultStr := "{\"result\": \"ok\"}"
	log.Println("starting")

	j, err := json.Marshal(resultStr)
	if err == nil {
		resultStr = string(j)
	} else if err != nil {
		log.Println(err)
		return events.ALBTargetGroupResponse{StatusCode: 500, Body: "JSON encode failed"}, err
	}
	headers := map[string]string{}
	headers["Content-Type"] = "application/json; charset=utf-8"
	headers["Cache-Control"] = "max-age=0"

	response := events.ALBTargetGroupResponse{Body: resultStr, StatusCode: 200, StatusDescription: "HTTP OK", IsBase64Encoded: false, Headers: headers}
	return response, nil
}

// net/http handler
func HttpHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(w)
	//fmt.Println(r)

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	headers := map[string]string{}
	multiValueHeaders := map[string][]string{}

	for name, values := range r.Header {
		if len(values) == 1 {
			headers[strings.ToLower(name)] = values[0]
		} else {
			multiValueHeaders[strings.ToLower(name)] = values
		}
	}

	request := events.ALBTargetGroupRequest{Headers: headers, MultiValueHeaders: multiValueHeaders, Body: string(body)}

	response, _ := Handler(context.TODO(), request)
	if response.StatusCode != 200 {
		w.WriteHeader(response.StatusCode)
	}
	fmt.Println(response)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", response.Body)
}
