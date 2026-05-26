package photon

import "github.com/AutoDruid/photon-parser/internal/types"

// Option is a functional option that configures a [Parser] at construction time.
// Pass Options to [NewParserV16] or [NewParserV18].
type Option func(*types.Config)

// defaultConfig returns a Config with all parsing enabled and no commands or
// event codes skipped. It is applied automatically by [NewParserV16] and
// [NewParserV18] before any caller-supplied Options are evaluated.
func defaultConfig() types.Config {
	return types.Config{}
}

// SkipUnknownPayloads controls whether the parser silently skips command
// payloads it does not recognise (for example encrypted or future command
// types) instead of storing their raw bytes.
//
// Enabling this option reduces allocations when unknown payloads are frequent
// and their raw content is not needed.
//
//	parser := photon.NewParserV18(photon.SkipUnknownPayloads(true))
func SkipUnknownPayloads(skip bool) Option {
	return func(c *types.Config) {
		c.SkipUnknownPayloads = skip
	}
}

// SkipParameterParsing controls whether the parser skips decoding parameters
// inside reliable and unreliable payloads.
//
// Enable this option when you only need command-level metadata (type, sequence
// number, event code) and do not need to inspect individual parameters. It
// avoids all per-parameter allocations on the hot path.
//
//	parser := photon.NewParserV18(photon.SkipParameterParsing(true))
func SkipParameterParsing(skip bool) Option {
	return func(c *types.Config) {
		c.SkipParameterParsing = skip
	}
}

// SkipCommands registers one or more command types that the parser should skip
// entirely. The payload bytes of a skipped command are consumed but not decoded,
// and no hooks are fired for it.
//
// This is useful when a capture contains high-frequency commands (such as
// [AcknowledgeCommand] or [PingCommand]) that are irrelevant to your analysis.
//
//	parser := photon.NewParserV18(
//	    photon.SkipCommands(photon.AcknowledgeCommand, photon.PingCommand),
//	)
func SkipCommands(commands ...types.CommandType) Option {
	return func(c *types.Config) {
		if len(commands) == 0 {
			return
		}
		if c.SkipCommands == nil {
			c.SkipCommands = make(map[types.CommandType]bool, len(commands))
		}
		for _, t := range commands {
			c.SkipCommands[t] = true
		}
	}
}

// SkipTargetEventCodes registers one or more message types whose parameters
// the parser should skip. The reliable header is still decoded (giving you
// the event code and parameter count) but the parameter bytes are consumed
// without being parsed, and no parameter hooks are fired.
//
// Use this option to ignore high-volume event types that are not relevant to
// your workload, reducing CPU and allocation overhead.
//
//	parser := photon.NewParserV18(
//	    photon.SkipTargetEventCodes(photon.OperationRequest, photon.OperationResponse),
//	)
func SkipTargetEventCodes(codes ...types.MessageType) Option {
	return func(c *types.Config) {
		if len(codes) == 0 {
			return
		}
		if c.SkipTargetEventCodes == nil {
			c.SkipTargetEventCodes = make(map[types.MessageType]bool, len(codes))
		}
		for _, code := range codes {
			c.SkipTargetEventCodes[code] = true
		}
	}
}
