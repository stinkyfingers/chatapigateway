package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	chat "github.com/stinkyfingers/chatapigateway"
)

type Event struct {
	Body           string              `json:"body"`
	RequestContext chat.RequestContext `json:"requestContext"`
}

type Data struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

type Response struct {
	IsBase64 bool              `json:"isBase64Encoded"`
	Status   int               `json:"statusCode"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
}

const (
	region   = "us-west-1"
	profile  = "jds"
	endpoint = "https://1u2upi76ni.execute-api.us-west-1.amazonaws.com/test"
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
	connection := &chat.Connection{Channel: d.Channel, ConnectionID: event.RequestContext.ConnectionID, Sess: sess}

	conns, err := chat.GetChannelConnections(sess, d.Channel)
	if err != nil {
		return r, err
	}
	err = connection.PostToConnections(conns, []byte(d.Text), false)
	if err != nil {
		return r, err
	}

	r.Status = 200
	return r, nil
}
