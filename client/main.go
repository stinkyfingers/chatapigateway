package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

//var  uri = "ws://localhost:8080/"
var uri = "wss://1u2upi76ni.execute-api.us-west-1.amazonaws.com/test"

func main() {
	con, res, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		b, _ := ioutil.ReadAll(res.Body)
		log.Print(string(b))
		log.Print(res)
		log.Fatal(err)
	}

	go func() {
		for {
			_, b, err := con.ReadMessage()
			if err != nil {
				log.Print("read err")
				log.Fatal(err)
			}
			log.Print(string(b))
		}
	}()

	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		log.Print("sending ")
		err = con.WriteMessage(1, scan.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}

}
