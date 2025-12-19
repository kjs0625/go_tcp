package protocol

const (
	SERVICE_1 = (10001)
	SERVICE_2 = (10002)
	SERVICE_3 = (10003)
	SERVICE_4 = (10004)

	MAX_PACKET_SIZE = 4096
)

type PacketHeader struct {
	UiPacketSize uint32
	UsPacketType uint16
}
