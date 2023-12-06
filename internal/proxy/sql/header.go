package sql

import "encoding/binary"

type Header struct {
	PacketType   int
	PacketLength int
}

func CreateHeaderFromBytes(bytes []byte) Header {
	packetLength := int(binary.BigEndian.Uint32(bytes[1:5]))
	return Header{
		PacketType:   int(bytes[0]),
		PacketLength: packetLength,
	}
}
