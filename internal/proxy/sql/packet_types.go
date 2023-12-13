package sql

var packetTypes = map[int]string{
	0x52: "AUTHENTICATION_REQUEST",
	0x31: "PARSE_COMPLETION",
	0x32: "BIND_COMPLETION",
	0x42: "BIND",
	0x43: "COMMAND_COMPLETION",
	0x44: "DESCRIBE",
	0x45: "EXECUTE",
	0x4B: "BACKEND_KEY_DATA",
	0x50: "PARSE",
	0x51: "QUERY",
	0x53: "PARAMETER_STATUS",
	0x54: "ROW_DESCRIPTION",
	0x58: "TERMINATION",
	0x5A: "READY_FOR_QUERY",
	0x70: "SASL_RESPONSE_MESSAGE",
}

func GetPacketType(typeFlag int) string {
	return packetTypes[typeFlag]
}
