package sql

import "encoding/binary"

type PGHeader struct {
	PacketType   int
	PacketLength int
}

func CreateHeaderFromBytes(bytes []byte) PGHeader {
	packetLength := int(binary.BigEndian.Uint32(bytes[1:5]))
	return PGHeader{
		PacketType:   int(bytes[0]),
		PacketLength: packetLength,
	}
}
