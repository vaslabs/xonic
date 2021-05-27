package main

import (
	"time"
	"github.com/vaslabs/codecs"
	"github.com/gorilla/websocket"
	evdev "github.com/gvalkov/golang-evdev"
	b64 "encoding/base64"
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
	conn *websocket.Conn
	not_registered_gamepad *evdev.InputDevice
	registered_gamepad *evdev.InputDevice
}

func NewClient(conn *websocket.Conn) Client {
	return Client{conn, nil, nil}
}



func (client *Client) Stream_Gamepad(inputDevice *evdev.InputDevice) {
	device_info := codecs.Encode_Device(inputDevice)
	device_info_payload := b64.StdEncoding.EncodeToString(device_info)

	client.send(&GamepadMessage{"RegisterDevice", device_info_payload})

	client.awaitRegistration(inputDevice)
}


func (client *Client) awaitRegistration(inputDevice *evdev.InputDevice) {
	client.not_registered_gamepad = inputDevice
}

func (client *Client) Close() error {
	return client.conn.Close()
}

func (client *Client) send(message *GamepadMessage) {
	client.conn.WriteJSON(message)
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

type GamepadMessage struct {
	command string
	payload string
}
