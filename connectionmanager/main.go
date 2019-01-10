package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	chat "github.com/stinkyfingers/chatapigateway"
)

type Event struct {
	Body           string              `json:"body"`
	RequestContext chat.RequestContext `json:"requestContext"`
	Records        interface{}         `json:"records"`
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
	region   = "us-west-1"
	profile  = "jds"
	endpoint = "https://1u2upi76ni.execute-api.us-west-1.amazonaws.com/test"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, event Event) (Response, error) {
	var r Response

	sess, err := chat.NewAwsSession(region, profile)
	if err != nil {
		return r, err
	}
	connection := &chat.Connection{Channel: "testchannel", ConnectionID: event.RequestContext.ConnectionID, Sess: sess}

	switch event.RequestContext.EventType {
	case "CONNECT":
		err = chat.Put(sess, *connection)
		if err != nil {
			log.Print("connect error", err)
			return r, err
		}

	case "DISCONNECT":
		err = connection.PostToConnections([]chat.Connection{*connection}, []byte("disconnecting"), true)
		if err != nil {
			return r, err
		}
		err = chat.Delete(sess, "testchannel", event.RequestContext.ConnectionID)
		if err != nil {
			log.Print("disconnect error", err)
			return r, err
		}

	default:
		err = connection.PostToConnections([]chat.Connection{*connection}, []byte("I don't know what to do. Specify an {\"option\":}"), true)
		if err != nil {
			return r, err
		}
	}

	r.Status = 200
	r.Headers = map[string]string{"testheader": "I'm a header"}

	return r, nil
}
