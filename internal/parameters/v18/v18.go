package v18

import (
	"encoding/json"

	"github.com/AutoDruid/photon-parser/internal/types"
)

// Header represents the parameter header containing the parameter ID and type.
// This appears at the beginning of each serialized parameter.
type Header struct {
	ID   uint8         `json:"id"`   // Parameter identifier (application-specific)
	Type ParameterType `json:"type"` // Protocol16 type code indicating how to decode the value
}

// Parameters represents a complete Photon Protocol parameter with its header and decoded value.
// The Value field contains the decoded data according to the Type specified in the Header.
type Parameter struct {
	Header `json:"header"`
	Value  `json:"value"`
}

var _ types.ParameterView = (*Parameter)(nil)

type Value struct {
	Kind    ParameterType `json:"kind"`
	KeyType ParameterType `json:"key_type"`
	ValType ParameterType `json:"val_type"`
	_       [5]byte       `json:"-"`
	Num     uint64        `json:"num"`
	Blob    []byte        `json:"blob,omitempty"`
}

func (p Parameter) ID() uint8 {
	return p.Header.ID
}

func (p Parameter) MarshalJSON() ([]byte, error) {
	type Alias Parameter

	out := struct {
		Alias
		Decoded any `json:"decoded,omitempty"`
	}{
		Alias: Alias(p),
	}

	switch p.Kind {
	case Int8Type, Int8Positive, Int8Negative,
		Int16Type, Int16Positive, Int16Negative,
		CompressedInt32Type,
		Long8Positive, Long8Negative,
		Long16Positive, Long16Negative,
		CompressedInt64Type,
		IntZeroType, ShortZeroType, LongZeroType, ByteZeroType:
		out.Decoded, _ = p.IntValue()
	case StringType:
		out.Decoded, _ = p.StringValue()
	case Float32Type, FloatZeroType:
		out.Decoded, _ = p.Float32Value()
	case Float64Type, DoubleZeroType:
		out.Decoded, _ = p.Float64Value()
	case BooleanType:
		out.Decoded, _ = p.BooleanValue()
	case Float32ArrayType:
		out.Decoded = collect(p.Float32ArrayValue(), p.Num)
	case CompressedIntArrayType:
		out.Decoded = collect(p.Int32ArrayValue(), p.Num)
	case CompressedLongArrayType:
		out.Decoded = collect(p.Int64ArrayValue(), p.Num)
	case ByteArrayType:
		out.Decoded = collect(p.ByteArrayValue(), p.Num)
	case ShortArrayType:
		out.Decoded = collect(p.Int16ArrayValue(), p.Num)
	case StringArrayType:
		out.Decoded = collect(p.StringArrayValue(), p.Num)
	case ArrayType:
		out.Decoded = collect(p.ArrayValue(), p.Num)
	case BooleanArrayType:
		out.Decoded = collect(p.BooleanArrayValue(), p.Num)
	default:
		out.Decoded = p.Num
	}

	return json.Marshal(out)
}
