package codecs

import (
	"bytes"
	"encoding/binary"

)

const (
	CAPABILITIES_FLAG uint32 = 0xCABE
	NAME_STRING uint32 = 0x5A4E
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

/*
magic for string (4 bytes), size of string (4 bytes), string bytes
*/
func Encode_Identifieable_String(magic uint32, value string) []byte {
	total_size := 4 + 4 + len(value)
	store := bytes.NewBuffer(make([]byte, total_size))

	uint32_store := make([]byte, 4)

	write32(store, uint32_store, magic)
	write32(store, uint32_store, uint32(len(value)))

	store.WriteString(value)
	byte_buffer := store.Bytes()
	return byte_buffer[len(byte_buffer) - total_size:]
}

func Decode_Identifieable_String(encoded []byte) (uint32, string) {
	flag := binary.BigEndian.Uint32(encoded[0:4])

	size_of_string := binary.BigEndian.Uint32(encoded[4:8])

	return flag, string(encoded[8:8+size_of_string])
}

func read_key(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}

func write32(buffer *bytes.Buffer, store []byte, value uint32) {
	binary.BigEndian.PutUint32(store, value)
	buffer.Write(store)
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