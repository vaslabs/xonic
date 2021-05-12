package main

import (
	"fmt"

	evdev "github.com/gvalkov/golang-evdev"
	"github.com/vaslabs/codecs"
)


func user_output(devices []*evdev.InputDevice) {
	for i := 0; i < len(devices); i++ {
		device_address := devices[i].File.Name()
		fmt.Printf("%d. %s\n", i, device_address)
	}
}

func list_input_devices() []*evdev.InputDevice {
	devices, err := evdev.ListInputDevices("/dev/input/*")

	if err == nil {
		return devices
	} else {
		fmt.Printf("Error getting input devices: %s", err.Error())
		return []*evdev.InputDevice{}
	}
}

func main() {
	devices := list_input_devices()
	user_output(devices)
	if len(devices) == 0 {
		return
	}
	var selected_device int
	_, err := fmt.Scanf("%d\n", &selected_device)

	if err != nil {
		fmt.Printf("Unrecognised input: %s\n", err.Error())
		return
	}
	if selected_device < 0 || selected_device >= len(devices) {
		fmt.Printf("Unrecognised input: %d, accepted numbers are 0-%d\n", selected_device, len(devices)-1)
		return
	}

	fmt.Printf("Selected %d\n", selected_device)
	fmt.Println("======================")
	fmt.Printf("%v\n", codecs.Decode_Device(codecs.Encode_Device(devices[selected_device])))
	device := devices[selected_device]

	event, _ := device.ReadOne()
	
	fmt.Println(event)
	encoded_event := codecs.Encode_Input_Event(event) 
	fmt.Println(encoded_event)

	decoded_event := codecs.Decode_Input_Event(encoded_event)
	fmt.Println(decoded_event)


}
