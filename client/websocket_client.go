package client

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/vaslabs/protocol"
)

func Run() {
	u := url.URL{Scheme: "ws", Host: "0.0.0.0:8080", Path: "/gamepads"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	c.WriteJSON(protocol.GamepadMessage{Command: protocol.ListDevices, Payload: ""})

	message := protocol.GamepadMessage{}
	rerr := c.ReadJSON(&message)

	if (rerr != nil) {
		fmt.Println(rerr.Error())
	}
	fmt.Printf("Devices are %v", message)
}
