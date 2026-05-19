package command

import (
	"fmt"

	"github.com/AutoDruid/photon-parser/internal/command/acknowledge"
	"github.com/AutoDruid/photon-parser/internal/command/connect"
	"github.com/AutoDruid/photon-parser/internal/command/reliable"
	"github.com/AutoDruid/photon-parser/internal/context"
	"github.com/AutoDruid/photon-parser/internal/errors"
	"github.com/AutoDruid/photon-parser/internal/hooks"
	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/types"
)

// ParseFromReader parses a Photon command from a parser.Reader.
// It first reads the 12-byte command header, validates the length field,
// then reads the remaining payload data.
//
// Returns an error if:
//   - The header cannot be read
//   - The length field is smaller than the header size (invalid)
//   - The payload data cannot be fully read
//
// The returned Command struct contains all header fields and the raw payload
// in the Data field. For SendReliable commands, the Data can be further parsed
// using the reliable package.
func ParseInto[P types.ParameterView](ctx *context.Context[P], dest *types.Command[P]) error {
	err := readCommandHeaderInto(ctx.Reader, dest)

	if !types.IsKnownCommandType(dest.Type) {
		remaining := ctx.Reader.Max - ctx.Reader.Cursor - 1

		if ctx.Config.SkipUnknownPayloads {
			return ctx.Reader.Skip(remaining)
		}

		rest, err := ctx.Reader.ReadBytes(remaining)
		if err != nil {
			return err
		}
		dest.UnknownPayload = types.UnknownPayload{Raw: rest, Kind: dest.Type}

		emit(ctx.Hooks, dest)
		return nil
	}

	if err != nil {
		return err
	}

	if dest.Length < types.COMMAND_HEADER_SIZE {
		return errors.ErrHeaderSize
	}

	if ctx.Config.SkipCommands[dest.Type] {
		remaining := int(dest.Length - types.COMMAND_HEADER_SIZE)
		return ctx.Reader.Skip(remaining)
	}

	err = readCommandPayloadInto(ctx, dest)
	if err != nil {
		remaining := int(dest.Length - types.COMMAND_HEADER_SIZE)

		if ctx.Config.SkipUnknownPayloads {
			return ctx.Reader.Skip(remaining)
		}

		rest, _ := ctx.Reader.ReadBytes(remaining)
		// don't fatal — just store raw for encrypted packets
		dest.UnknownPayload = types.UnknownPayload{Raw: rest, Kind: dest.Type}
	}

	emit(ctx.Hooks, dest)

	return nil
}

func readCommandHeaderInto[P types.ParameterView](r *reader.Reader, dest *types.Command[P]) error {
	var err error

	b, err := r.ReadUInt8()
	if err != nil {
		return err
	}

	dest.Type = types.CommandType(b)

	if !types.IsKnownCommandType(dest.Type) {
		return nil
	}

	dest.ChannelID, err = r.ReadUInt8()
	if err != nil {
		return err
	}

	dest.Flags, err = r.ReadUInt8()
	if err != nil {
		return err
	}

	dest.ReservedByte, err = r.ReadUInt8()
	if err != nil {
		return err
	}

	dest.Length, err = r.ReadUInt32BE()
	if err != nil {
		return err
	}

	dest.ReliableSequenceNumber, err = r.ReadUInt32BE()
	if err != nil {
		return err
	}

	return nil
}

func readCommandPayloadInto[P types.ParameterView](ctx *context.Context[P], dest *types.Command[P]) error {
	switch dest.Type {
	case types.SendUnreliableCommand, types.SendUnreliableUnsequenced:

		_, err := ctx.Reader.ReadBytes(4)
		if err != nil {
			return err
		}

		relPayload := int64(dest.Length) - types.COMMAND_HEADER_SIZE - 4
		if relPayload < 0 {
			return fmt.Errorf("command length %d too small for unreliable reliable payload", dest.Length)
		}

		err = reliable.Parse(ctx, &dest.UnreliablePayload, uint32(relPayload))
		if err != nil {
			return err
		}
	case types.SendReliableCommand:
		relPayload := int64(dest.Length) - types.COMMAND_HEADER_SIZE
		if relPayload < 0 {
			return fmt.Errorf("command length %d too small for reliable payload", dest.Length)
		}
		err := reliable.Parse(ctx, &dest.ReliablePayload, uint32(relPayload))
		if err != nil {
			return err
		}
	case types.AcknowledgeCommand:
		err := acknowledge.ParseInto(ctx.Reader, &dest.AcknowledgePayload)
		if err != nil {
			return err
		}
	case types.ConnectCommand, types.VerifyConnectCommand:
		err := connect.ParseInto(ctx.Reader, &dest.ConnectPayload)
		if err != nil {
			return err
		}
	case types.SendReliableFragmentCommand:
		err := reliable.ParseIntoFragment(ctx, dest.Length, &dest.ReliableFragmentPayload, &dest.ReliablePayload)
		if err != nil {
			return err
		}
	case types.PingCommand:
		dest.PingPayload = struct{}{}
	case types.FetchServerTimestampCommand:
		dest.FetchTimestampPayload = struct{}{}
	case types.DisconnectCommand:
		dest.DisconnectPayload = struct{}{}
	default:
		return fmt.Errorf("unsupported command type %d", dest.Type)
	}

	return nil
}

func emit[P types.ParameterView](hooks *hooks.Hooks[P], dest *types.Command[P]) {
	if hooks == nil {
		return
	}

	if hooks.SyncHooks.OnCommand != nil {
		hooks.SyncHooks.OnCommand(*dest)
	}

	if hooks.AsyncHooks.OnCommand == nil {
		return
	}

	s := DetachForAsync(*dest)

	select {
	case hooks.AsyncHooks.OnCommand <- s:
	default: // don't block parser
	}

}

func DetachForAsync[P types.ParameterView](cmd types.Command[P]) types.Command[P] {
	s := cmd
	switch s.Type {
	case types.SendReliableCommand:
		if n := len(s.ReliablePayload.Parameters); n > 0 {
			p := make([]P, n)
			copy(p, s.ReliablePayload.Parameters)
			s.ReliablePayload.Parameters = p
		}
	case types.SendUnreliableCommand, types.SendUnreliableUnsequenced:
		if n := len(s.UnreliablePayload.Parameters); n > 0 {
			p := make([]P, n)
			copy(p, s.UnreliablePayload.Parameters)
			s.UnreliablePayload.Parameters = p
		}
	}
	return s
}
