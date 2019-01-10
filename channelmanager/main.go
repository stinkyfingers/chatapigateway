package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	chat "github.com/stinkyfingers/chatapigateway"
)

type Event struct {
	Body           string              `json:"body"`
	RequestContext chat.RequestContext `json:"requestContext"`
	Records        interface{}         `json:"records"`
}

type Data struct {
	Action  string `json:"action"`
	Channel string `json:"channel"`
}

type Output struct {
	Message string `json:"message"`
}

type Response struct {
	IsBase64 bool              `json:"isBase64Encoded"`
	Status   int               `json:"statusCode"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
}

const (
	region  = "us-west-1"
	profile = "jds"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, event Event) (Response, error) {
	r := Response{Status: 500}
	sess, err := chat.NewAwsSession(region, profile)
	if err != nil {
		return r, err
	}

	var d Data
	err = json.Unmarshal([]byte(event.Body), &d)
	if err != nil {
		return r, err
	}
	switch d.Action {
	case "subscribe":
		connection := chat.Connection{Channel: d.Channel, ConnectionID: event.RequestContext.ConnectionID}
		err = chat.Put(sess, connection)
		if err != nil {
			log.Print("subscribe error", err)
			return r, err
		}

	case "unsubscribe":
		fallthrough
	default:
		err = chat.Delete(sess, d.Channel, event.RequestContext.ConnectionID)
		if err != nil {
			log.Print("unsubscribe error", err)
			return r, err
		}
	}

	r.Status = 200
	return r, nil
}
