package server

/*
Hosts the gamepad, a client
registers to it to receive
gamepad commands
*/

import (
	b64 "encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	evdev "github.com/gvalkov/golang-evdev"
	"github.com/vaslabs/cmd"
	"github.com/vaslabs/codecs"
	"github.com/vaslabs/protocol"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	conn                   *websocket.Conn
	not_registered_gamepad *evdev.InputDevice
	registered_gamepad     *evdev.InputDevice
}

func NewClient(conn *websocket.Conn) Client {
	return Client{conn, nil, nil}
}

func (client *Client) Stream_Gamepad(inputDevice *evdev.InputDevice) {
	device_info := codecs.Encode_Device(inputDevice)
	device_info_payload := b64.StdEncoding.EncodeToString(device_info)

	client.send(&protocol.GamepadMessage{Command: "RegisterDevice", Payload: device_info_payload})

	client.await_registration(inputDevice)
}

func (client *Client) await_registration(inputDevice *evdev.InputDevice) {
	client.not_registered_gamepad = inputDevice
}

func (client *Client) process_messages() {
	defer client.Close()
	for {
		gamepad_message := &protocol.GamepadMessage{Command: "", Payload: ""}
		err := client.conn.ReadJSON(gamepad_message)
		if err != nil {
			log.Printf("Closing websocket connection due to %s", err.Error())
			break
		}
		if gamepad_message.Command == "Registered" {
			client.registered_gamepad = client.not_registered_gamepad
			client.stream_gamepad_input()
		} else if gamepad_message.Command == protocol.ListDevices {
			client.send(&protocol.GamepadMessage{Command: "SelectDevice", Payload: display(cmd.List_Input_Devices())})
		} else if gamepad_message.Command == "SelectedDevice" {
			selected_device, err := cmd.Find_Device(gamepad_message.Payload)
			if err != nil {
				client.Stream_Gamepad(selected_device)
			} else {
				fmt.Println(err.Error())
			}
		}
	}
}

func (client *Client) Close() error {
	return client.conn.Close()
}

func (client *Client) send(message *protocol.GamepadMessage) {
	client.conn.WriteJSON(message)
}

func (client *Client) stream_gamepad_input() {
	for {
		events, err := client.registered_gamepad.Read()
		if err == nil {
			for i := 0; i < len(events); i++ {
				payload := b64.StdEncoding.EncodeToString(codecs.Encode_Input_Event(&events[i]))
				client.conn.WriteJSON(protocol.GamepadMessage{Command: "Input", Payload: payload})
			}
		} else {
			log.Printf("Closing device input stream due to %s", err.Error())
			break
		}
	}
}

func configure_connection(ws *websocket.Conn) {
	go ping_writer(ws)
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

func ping_writer(ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for {
		<-pingTicker.C
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			return
		}
	}
}

func display(devices []*evdev.InputDevice) string {
	addresses := make([]string, len(devices))

	for i := 0; i < len(devices); i++ {
		device := devices[i]
		addresses[i] = device.Fn
	}
	return fmt.Sprintf("%v", addresses)
}

func (client *Client) open_web_socket(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)

	log.Print("Connecting to websocket")
	if err == nil {
		configure_connection(ws)
		if client.conn != nil {
			client.conn.Close()
		}
		client.conn = ws
		client.process_messages()
	} else {
		log.Printf("Error upgrading connection for websocket use: %s", err.Error())
	}
}

func Run() {
	client := NewClient(nil)
	ws_handler := client.open_web_socket

	http.HandleFunc("/gamepads", ws_handler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))

}
