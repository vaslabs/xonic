package main

import (
	"testing"

	evdev "github.com/gvalkov/golang-evdev"
	"github.com/vaslabs/codecs"
)


func Test_capabilities_decoding_inverse_of_encoding(t *testing.T) {
	input := map[int][]int{1:{1,2,3}, 2:{2,3,4}}
	size, encoded := codecs.Encode_Capabilities(input)
	expected_size := 4 + 8 + 2*8 + 8*3 + 2*8 + 8*3
	if (size != uint64(expected_size)) {
		t.Errorf("Expected size %d but got %d", expected_size, size)
	}

	if len(encoded) != expected_size {
		t.Errorf("%d != %d", len(encoded), expected_size)
	}

	decoded := codecs.Decode_Capabilities(encoded)

	if (decoded[1][0] != 1 ) {
		t.Errorf("%v\n!=\n%v", input, decoded)
	}


}

func Test_string_decoding_inverse_of_encoding(t *testing.T) {
	input := "Hello, world"
	magic := uint32(0x1234)
	encoded := codecs.Encode_Identifiable_String(magic, input)
	
	dmagic, decoded := codecs.Decode_Identifiable_String(encoded)

	if (dmagic != magic) {
		t.Errorf("Magic value mismatch: %d != %d", magic, dmagic)
	}

	if (input != decoded) {
		t.Errorf("%v\n", encoded)
		t.Errorf("Codec inverse failure: %s != %s", input, decoded)
	}
}

func Test_uint16_decoding_inverse_of_encoding(t *testing.T) {
	input := uint16(5432)
	magic := uint32(0x65AA)

	encoded := codecs.Encode_Identifiable_uint16(magic, input)

	dmagic, decoded := codecs.Decode_Identifiable_uint16(encoded)

	if (dmagic != magic) {
		t.Errorf("Magic value mismatch: %d != %d", magic, dmagic)
	}

	if (input != decoded) {
		t.Errorf("%v\n", encoded)
		t.Errorf("Codec inverse failure: %d != %d", input, decoded)
	}
}

func Test_Device_decoding_inverse_of_encoding(t *testing.T) {
	capabilities := map[int][]int{1:{1,2,3}, 2:{2,4}}

	input := evdev.InputDevice{Name: "Wireless ps4 controller", Vendor: 1, Product: 2, CapabilitiesFlat: capabilities}

	encoded := codecs.Encode_Device(&input)
	decoded := codecs.Decode_Device(encoded)

	if (decoded.Name != input.Name) {
		t.Errorf("Name mismatch %s != %s", decoded.Name, input.Name)
	}

	if (decoded.Vendor != input.Vendor) {
		t.Errorf("Vendor mismatch %d != %d", decoded.Vendor, input.Vendor)
	}


	if (decoded.Product != input.Product) {
		t.Errorf("Product mismatch %d != %d", decoded.Product, input.Product)
	}

	if (decoded.CapabilitiesFlat[2][1] != 4) {
		t.Errorf("CapabilitiesFlat mismatch %v\n !=\n %v\n", decoded.CapabilitiesFlat, input.CapabilitiesFlat)
	}
}