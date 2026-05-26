// Package session provides parsing for Photon Protocol session layer packets.
// The session layer is the outermost protocol layer, containing session metadata
// and one or more commands.
package session

import (
	"errors"

	"github.com/AutoDruid/photon-parser/internal/command"
	"github.com/AutoDruid/photon-parser/internal/context"
	photonErrors "github.com/AutoDruid/photon-parser/internal/errors"
	"github.com/AutoDruid/photon-parser/internal/hooks"
	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/types"
)

// ParseInto parses a Photon session packet from a parser.Reader.
// This function reads the session header, then iterates through and parses
// each command as specified by the CommandCount field.
//
// Returns an error if any part of parsing fails.
func ParseInto[P types.ParameterView](ctx *context.Context[P], dest *types.Session[P]) error {
	err := readSessionHeaderInto(ctx.Reader, dest)
	if err != nil {
		return err
	}

	items := ctx.PoolCommand.Get(int(dest.CommandCount))
	dest.Commands = items.Items

	for i := uint8(0); i < dest.CommandCount; i++ {
		err := command.ParseInto(ctx, &dest.Commands[i])

		if errors.Is(err, photonErrors.ErrHeaderSize) {
			break
		}
		if err != nil {
			return err
		}

		if !types.IsKnownCommandType(dest.Commands[i].Type) {
			break
		}
	}

	emit(ctx.Hooks, dest)

	ctx.PoolCommand.Put(items)

	return nil
}

func readSessionHeaderInto[P types.ParameterView](r *reader.Reader, dest *types.Session[P]) error {
	var err error

	dest.PeerID, err = r.ReadUInt16BE()
	if err != nil {
		return err
	}

	dest.CRCEnabled, err = r.ReadUInt8()
	if err != nil {
		return err
	}

	dest.CommandCount, err = r.ReadUInt8()
	if err != nil {
		return err
	}

	dest.Timestamp, err = r.ReadUInt32BE()
	if err != nil {
		return err
	}

	dest.Challenge, err = r.ReadInt32BE()
	if err != nil {
		return err
	}

	return nil
}

func emit[P types.ParameterView](hooks *hooks.Hooks[P], dest *types.Session[P]) {
	if hooks == nil {
		return
	}

	if hooks.SyncHooks.OnSession != nil {
		hooks.SyncHooks.OnSession(*dest)
	}

	ch := hooks.AsyncHooks.OnSession
	if ch == nil || len(ch) == cap(ch) {
		return
	}

	s := *dest
	if n := len(s.Commands); n > 0 {
		cmds := make([]types.Command[P], n)
		copy(cmds, s.Commands)
		for i := range cmds {
			cmds[i] = command.DetachForAsync(cmds[i])
		}
		s.Commands = cmds
	}

	select {
	case ch <- s:
	default:
	}
}
