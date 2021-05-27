package cmd

import (
	"fmt"

	evdev "github.com/gvalkov/golang-evdev"
)

func List_Input_Devices() []*evdev.InputDevice {
	devices, err := evdev.ListInputDevices("/dev/input/*")
	fmt.Printf("Devices %v\n", devices)
	for i := 0; i < len(devices); i++ {
		device := devices[i]
		fmt.Printf("%v\n", device.Capabilities)
		fmt.Printf("%v\n", device.CapabilitiesFlat)


	}
	if err == nil {
		return devices
	} else {
		fmt.Printf("Error getting input devices: %s", err.Error())
		return []*evdev.InputDevice{}
	}
}

func Find_Device(address string) (*evdev.InputDevice, error) {
	devices, err := evdev.ListInputDevices(address)
	if (err != nil) {
		return nil, err
	}
	if (len(devices) > 0) {
		return devices[0], nil
	}

	return nil, fmt.Errorf("Could not find any device on path %s", address)
}