package chatapigateway

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type RequestContext struct {
	EventType    string `json:"eventType"`
	RequestID    string `json:"requestId"`
	ConnectionID string `json:"connectionId"`
}

type Connection struct {
	ConnectionID string           `json:"connectionId"`
	Channel      string           `json:"channel"`
	Sess         *session.Session `json:"-"`
}

const (
	endpoint = "https://1u2upi76ni.execute-api.us-west-1.amazonaws.com/test"
)

// PostToConnections posts data to conns, optionally including self.
func (c *Connection) PostToConnections(conns []Connection, data []byte, includeSelf bool) error {
	gatewaySess := apigatewaymanagementapi.New(c.Sess, &aws.Config{
		Endpoint: aws.String(endpoint),
	})
	for _, conn := range conns {
		if !includeSelf && conn.ConnectionID == c.ConnectionID {
			continue
		}
		req, _ := gatewaySess.PostToConnectionRequest(&apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(conn.ConnectionID),
			Data:         data,
		})
		err := req.Send()
		if err != nil {
			if strings.Contains(err.Error(), request.ErrCodeSerialization) {
				err = Delete(c.Sess, c.Channel, conn.ConnectionID)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}
