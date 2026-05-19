package types

// CommandType represents a Photon Protocol command type.
type CommandType uint8

const (
	AcknowledgeCommand          CommandType = 0x01
	ConnectCommand              CommandType = 0x02
	VerifyConnectCommand        CommandType = 0x03
	DisconnectCommand           CommandType = 0x04
	PingCommand                 CommandType = 0x05
	SendReliableCommand         CommandType = 0x06
	SendUnreliableCommand       CommandType = 0x07
	SendReliableFragmentCommand CommandType = 0x08
	SendUnreliableUnsequenced   CommandType = 0x0B
	FetchServerTimestampCommand CommandType = 0x0C
)

// COMMAND_HEADER_SIZE is the size in bytes of a command header (12 bytes).
const COMMAND_HEADER_SIZE = 12

type CommandHeader struct {
	Type                   CommandType `json:"type"`
	ChannelID              uint8       `json:"channel_id"`
	Flags                  uint8       `json:"flags"`
	ReservedByte           uint8       `json:"reserved_byte"`
	Length                 uint32      `json:"length"`
	ReliableSequenceNumber uint32      `json:"reliable_sequence_number"`
}

// Command represents a Photon command with its header and payload data.
// To gain in performance we use a union of single struct for all command types.
type Command[P ParameterView] struct {
	CommandHeader `json:"header"`

	UnreliablePayload       Reliable[P]    `json:"unreliable_payload"`
	ReliablePayload         Reliable[P]    `json:"reliable_payload"`
	ReliableFragmentPayload Fragment       `json:"reliable_fragment_payload"`
	AcknowledgePayload      Acknowledge    `json:"acknowledge_payload"`
	ConnectPayload          Connect        `json:"connect_payload"`
	UnknownPayload          UnknownPayload `json:"unknown_payload"`
	PingPayload             struct{}       `json:"ping_payload"`
	FetchTimestampPayload   struct{}       `json:"fetch_timestamp_payload"`
	DisconnectPayload       struct{}       `json:"disconnect_payload"`
}

func IsKnownCommandType(t CommandType) bool {
	switch t {
	case AcknowledgeCommand,
		ConnectCommand,
		VerifyConnectCommand,
		DisconnectCommand,
		PingCommand,
		SendReliableCommand,
		SendUnreliableCommand,
		SendReliableFragmentCommand,
		SendUnreliableUnsequenced,
		FetchServerTimestampCommand:
		return true
	default:
		return false
	}
}

type Connect struct {
	Mtu                        uint32 `json:"mtu"`
	WindowSize                 uint32 `json:"window_size"`
	ChannelCount               uint32 `json:"channel_count"`
	IncomingBandwidth          uint32 `json:"incoming_bandwidth"`
	OutgoingBandwidth          uint32 `json:"outgoing_bandwidth"`
	DisconnectThrottle         uint32 `json:"disconnect_throttle"`
	PacketThrottleAcceleration uint32 `json:"packet_throttle_acceleration"`
	PacketThrottleDeceleration uint32 `json:"packet_throttle_deceleration"`
}

type Acknowledge struct {
	AckReliableSequenceNumber uint32 `json:"ack_reliable_sequence_number"`
	AckSentTime               uint32 `json:"ack_sent_time"`
}

type UnknownPayload struct {
	Raw  []byte      `json:"raw"`
	Kind CommandType `json:"kind"`
}

// Some commands data are fragmented, this struct represents the fragment metadata.
type Fragment struct {
	ID     uint32 `json:"id"`
	Count  uint32 `json:"count"`
	Index  uint32 `json:"index"`
	Size   uint32 `json:"size"`
	Offset uint32 `json:"offset"`
	Data   []byte `json:"data"`
}

type MessageType uint8

const (
	OperationRequest       MessageType = 0x02 // Client requests an operation
	OperationResponse      MessageType = 0x07 // Server responds to an operation
	OtherOperationResponse MessageType = 0x03 // Alternative response format
	EventDataType          MessageType = 0x04 // Server sends an event to client
	ExchangeKeys           MessageType = 0x06 // Key exchange for encryption
)

type ReliableHeader struct {
	Signature      uint8       `json:"signature"`       // Message signature (typically 0xF3)
	Type           MessageType `json:"type"`            // Message type (operation, event, etc.)
	EventCode      uint8       `json:"event_code"`      // Operation/event code (application-specific)
	ParameterCount int         `json:"parameter_count"` // Number of parameters following this header
}

type Reliable[P ParameterView] struct {
	ReliableHeader `json:"reliable_header"`
	Parameters     []P `json:"parameters"`
}
