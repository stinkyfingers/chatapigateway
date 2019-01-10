
##Client:##
go run client/main.go

- subscribe
{"option":"channel","action":"subscribe","channel":"foo"}
- post message
{"option":"message","text":"this is text","channel":"foo"}

##API Gateway Websockets:##
- RouteSelectionExpression: $request.body.option
- Routes:
  - $connect,$disconnect,$default point to connectionmanager lambda
  - channel points to channelmanager
  - message points to messagemanager
- Note: I've hardcoded the aws ws endpoint in the code for simplicity



Wscat
https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-how-to-call-websocket-api-wscat.html

