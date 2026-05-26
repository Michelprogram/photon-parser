// Public type aliases and wire-level constants for Photon decoding.
package photon

import (
	v16 "github.com/AutoDruid/photon-parser/internal/parameters/v16"
	v18 "github.com/AutoDruid/photon-parser/internal/parameters/v18"
	"github.com/AutoDruid/photon-parser/internal/types"
)

type Session[P types.ParameterView] = types.Session[P]
type Command[P types.ParameterView] = types.Command[P]

// SessionV16 is a decoded Photon session using protocol 16 parameters (see NewV16).
type SessionV16 = types.Session[v16.Parameter]

// SessionV18 is a decoded Photon session using protocol 18 parameters (see NewV18).
type SessionV18 = types.Session[v18.Parameter]

// CommandV18 is a single parsed Photon command with v18 parameter payloads.
type CommandV18 = types.Command[v18.Parameter]

// CommandV16 is a single parsed Photon command with v16 parameter payloads.
type CommandV16 = types.Command[v16.Parameter]

// Reliable is a parsed reliable/unreliable message body: header plus a parameter list.
// It appears inside Command payloads (for example SendReliable / SendUnreliable).
type Reliable[P types.ParameterView] = types.Reliable[P]

// HookOptions configures asynchronous hooks; Size is the channel buffer capacity.
type HookOptions = types.HookOptions

// ParameterV16 is the decoded v16 parameter value type used with Parser v16 APIs.
type ParameterV16 = v16.Parameter

// ParameterV18 is the decoded v18 parameter value type used with Parser v18 APIs.
type ParameterV18 = v18.Parameter

// ParameterV18Type is the v18 on-wire type discriminator for one parameter slot.
type ParameterV18Type = v18.ParameterType

// ParameterV16Type is the v16 on-wire type discriminator for one parameter slot.
type ParameterV16Type = v16.ParameterType

// ReliableV18 is Reliable[ParameterV18], the shape of reliable payloads for NewV18 parsers.
type ReliableV18 = Reliable[ParameterV18]

// ReliableV16 is Reliable[ParameterV16], the shape of reliable payloads for NewV16 parsers.
type ReliableV16 = Reliable[ParameterV16]

type CommandType = types.CommandType
type MessageType = types.MessageType

// Photon command type bytes (CommandHeader.Type).
const (
	AcknowledgeCommand          = types.AcknowledgeCommand
	ConnectCommand              = types.ConnectCommand
	VerifyConnectCommand        = types.VerifyConnectCommand
	DisconnectCommand           = types.DisconnectCommand
	PingCommand                 = types.PingCommand
	SendReliableCommand         = types.SendReliableCommand
	SendUnreliableCommand       = types.SendUnreliableCommand
	SendReliableFragmentCommand = types.SendReliableFragmentCommand
	SendUnreliableUnsequenced   = types.SendUnreliableUnsequenced
	FetchServerTimestampCommand = types.FetchServerTimestampCommand
)

// Reliable payload message kinds (ReliableHeader.Type).
const (
	OperationRequest       = types.OperationRequest
	OperationResponse      = types.OperationResponse
	OtherOperationResponse = types.OtherOperationResponse
	EventDataType          = types.EventDataType
	ExchangeKeys           = types.ExchangeKeys
)
