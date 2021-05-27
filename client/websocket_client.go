package client

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	// evdev "github.com/gvalkov/golang-evdev"
	"github.com/vaslabs/protocol"
	// uinput "github.com/ynsta/uinput"
)

func Run(address string) {
	u := url.URL{Scheme: "ws", Host: address, Path: "/gamepads"}
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

	var gamepad string
	fmt.Scanf("Select device %s", gamepad)

	c.WriteJSON(protocol.GamepadMessage{Command: protocol.SelectedDevice, Payload: gamepad})


	// for {
	// 	next_message := c.ReadJSON(&message)
	// 	// if (message.Command == protocol.RegisterDevice)
	// }
}

// func register(input_device *evdev.InputDevice) uinput.WriteDevice {
// 	var wd uinput.WriteDevice
// 	var ui uinput.UInput

// 	wd.Open()

// 	ui.Init(&wd, input_device.Name, input_device.Vendor, input_device.Product, 0x8100, 
// 	)
// }
