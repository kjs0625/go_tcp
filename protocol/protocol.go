package protocol

const (
	SERVICE_1      = (10001)
	SERVICE_1_RESP = (10002)
	SERVICE_2      = (10003)
	SERVICE_2_RESP = (10004)
	SERVICE_3      = (10005)
	SERVICE_3_RESP = (10006)
	SERVICE_4      = (10007)
	SERVICE_4_RESP = (10008)

	MAX_PACKET_SIZE = 4096
)

type PacketHeader struct {
	UiPacketSize uint32
	UsPacketType uint16
}
