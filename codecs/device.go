package codecs

import (
	"bytes"
	"encoding/binary"

)

const (
	CAPABILITIES_FLAG uint32 = 0xCABE
)

/*
 F L A G 
|_|_|_|_|....

If capabilities
|C|A|B|A|_|_|_|_ |_|_|_|_ |_|_|_|_|          |_|_|_|_|     
		Map size KEY val  size of array    array value    ...repeat next key
*/

func Decode_Capabilities(encoded []byte) map[int][]int  {
	flag := binary.BigEndian.Uint32(encoded[0:4])
	if (flag != CAPABILITIES_FLAG) {
		return map[int][]int{}
	} else {
		return decode_capabilities(encoded[4:])
	}
}

func read_key(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}

func decode_capabilities(encoded_raw []byte) map[int][]int {
	keys := binary.BigEndian.Uint64(encoded_raw[0:8])
	capabilities := make(map[int][]int, keys)
	pointer := 8
	for keys_read := uint64(0); keys_read < keys; keys_read++ {
		key_value := int(read_key(encoded_raw[pointer:pointer+8]))
		pointer += 8
		array_size := read_key(encoded_raw[pointer:pointer+8])
		array := make([]int, array_size)
		capabilities[key_value] = array

		pointer+=8
		for value_index := 0; value_index < int(array_size); value_index++ {
			next_value := int(read_key(encoded_raw[pointer:pointer+8]))
			pointer+=8
			array[value_index] = next_value
		}
	}
	return capabilities
}

func Encode_Capabilities(capabilities map[int][]int) (uint64, []byte) {
	
	key_size, size := calculate_size(capabilities)
	header_size := uint64(4)

	store := bytes.NewBuffer(make([]byte, size + header_size))
	uint32_store := make([]byte, 4)
	binary.BigEndian.PutUint32(uint32_store, uint32(CAPABILITIES_FLAG))
	store.Write(uint32_store)

	ephimeral_store := make([]byte, 8)


	binary.BigEndian.PutUint64(ephimeral_store, uint64(key_size))
	store.Write(ephimeral_store)
	for key, elements := range capabilities {
		binary.BigEndian.PutUint64(ephimeral_store, uint64(key))
		store.Write(ephimeral_store)
		binary.BigEndian.PutUint64(ephimeral_store, uint64(len(elements)))
		store.Write(ephimeral_store)
		for i := 0; i < len(elements); i++ {
			binary.BigEndian.PutUint64(ephimeral_store, uint64(elements[i]))
			store.Write(ephimeral_store)
		}
	}
	total_size := size + header_size
	buffer_bytes := store.Bytes()
	return (size + header_size), buffer_bytes[len(buffer_bytes) - int(total_size):]
}

func calculate_size(multi_map map[int][]int) (uint64, uint64) {
	elements_size := uint64(0)
	arrays := uint64(0)
	for _, element := range multi_map {
		elements_size += uint64(len(element))
		arrays+=1
	}
	keys := arrays
	total_size := 8 + arrays*8 + uint64(elements_size)*8 + arrays*8
	//each key has a value and then the size of the array it points to + the size of all integers in the array
	return keys, total_size
}