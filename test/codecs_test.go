package main

import (
	"testing"

	"github.com/vaslabs/codecs"
)


func Test_decoding_inverse_of_encoding(t *testing.T) {
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