package chatapigateway

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var tableName = "chatter"

func Put(sess *session.Session, c Connection) error {
	item, err := dynamodbattribute.MarshalMap(c)
	if err != nil {
		return err
	}

	_, err = dynamodb.New(sess).PutItem(&dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})
	return err
}

func Get(sess *session.Session, channel, connectionId string) (Connection, error) {
	var rc Connection

	ch, err := dynamodbattribute.Marshal(channel)
	if err != nil {
		return rc, err
	}
	co, err := dynamodbattribute.Marshal(connectionId)
	if err != nil {
		return rc, err
	}

	resp, err := dynamodb.New(sess).GetItem(&dynamodb.GetItemInput{
		TableName: &tableName,
		Key:       map[string]*dynamodb.AttributeValue{"channel": ch, "connectionId": co},
	})
	if err != nil {
		return rc, err
	}
	err = dynamodbattribute.UnmarshalMap(resp.Item, &rc)
	return rc, err
}

func GetChannelConnections(sess *session.Session, channel string) ([]Connection, error) {
	var connections []Connection

	resp, err := dynamodb.New(sess).Query(&dynamodb.QueryInput{
		TableName: &tableName,
		KeyConditions: map[string]*dynamodb.Condition{
			"channel": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(channel),
					},
				},
			},
		},
	})
	if err != nil {
		return connections, err
	}

	for _, item := range resp.Items {
		var rc Connection
		err = dynamodbattribute.UnmarshalMap(item, &rc)
		if err != nil {
			return connections, err
		}
		connections = append(connections, rc)
	}
	return connections, err
}

func Delete(sess *session.Session, channel, connectionId string) error {
	ch, err := dynamodbattribute.Marshal(channel)
	if err != nil {
		return err
	}
	co, err := dynamodbattribute.Marshal(connectionId)
	if err != nil {
		return err
	}
	_, err = dynamodb.New(sess).DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key:       map[string]*dynamodb.AttributeValue{"channel": ch, "connectionId": co},
	})
	return err
}
