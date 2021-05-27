package protocol

type GamepadMessage struct {
	Command string
	Payload string
}


const (
	ListDevices = "ListDevices"
)