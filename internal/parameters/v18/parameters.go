package v18

import (
	"fmt"

	"github.com/AutoDruid/photon-parser/internal/context"
	"github.com/AutoDruid/photon-parser/internal/hooks"
	"github.com/AutoDruid/photon-parser/internal/reader"
)

var _ context.ParameterParser[Parameter] = (*Parameter)(nil)

// Parse reads a complete parameter from the reader.
// Format: Header (1 byte ID + 1 byte Type), followed by the typed value.
//
// The function first reads the parameter header to determine the parameter ID
// and type code, then decodes the value according to that type using the
// Protocol16 decoder.
//
// Returns a Parameters struct containing the ID, Type, and decoded Value,
// or an error if parsing fails.
//
// Example usage:
//
//	param, err := parameters.Parse(reader)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Parameter %d has value: %v\n", param.ID, param.Value)
func (p *Parameter) ParseInto(reader *reader.Reader, hooks *hooks.Hooks[Parameter], dest *Parameter) error {

	header, err := p.parseHeader(reader)
	if err != nil {
		return err
	}

	value := Value{Kind: header.Type}

	err = scanPayload(reader, &value)

	if err != nil {
		return err
	}

	dest.Header = header
	dest.Value = value

	p.emit(hooks, dest)

	return nil
}

func (p Parameter) emit(hooks *hooks.Hooks[Parameter], out *Parameter) {
	if hooks == nil {
		return
	}

	if hooks.SyncHooks.OnParameter != nil {
		hooks.SyncHooks.OnParameter(*out)
	}

	if hooks.AsyncHooks.OnParameter == nil {
		return
	}

	select {
	case hooks.AsyncHooks.OnParameter <- *out:
	default:
	}
}

func (p *Parameter) parseHeader(r *reader.Reader) (Header, error) {
	var err error
	var header Header

	header.ID, err = r.ReadUInt8()
	if err != nil {
		return Header{}, err
	}

	b, err := r.ReadUInt8()
	if err != nil {
		return Header{}, err
	}

	header.Type = ParameterType(b)

	return header, nil
}

func scanPayload(reader *reader.Reader, dest *Value) error {
	var err error

	switch dest.Kind {
	case Int8Type:
		err = scanInt8(reader, dest)
		if err != nil {
			return err
		}
	case Int8Positive:
		err = scanInt8Positive(reader, dest)
		if err != nil {
			return err
		}
	case Int8Negative:
		err = scanInt8Negative(reader, dest)
		if err != nil {
			return err
		}
	case Int16Type:
		err = scanInt16Type(reader, dest)
		if err != nil {
			return err
		}
	case Int16Positive:
		err = scanInt16Positive(reader, dest)
		if err != nil {
			return err
		}
	case Int16Negative:
		err = scanInt16Negative(reader, dest)
		if err != nil {
			return err
		}
	case Long8Positive:
		err = scanLong8Positive(reader, dest)
		if err != nil {
			return err
		}
	case Long8Negative:
		err = scanLong8Negative(reader, dest)
		if err != nil {
			return err
		}
	case Long16Positive:
		err = scanLong16Positive(reader, dest)
		if err != nil {
			return err
		}
	case Long16Negative:
		err = scanLong16Negative(reader, dest)
		if err != nil {
			return err
		}
	case StringType:
		err = scanString(reader, dest)
		if err != nil {
			return err
		}
	case CompressedInt32Type:
		err = scanCompressedInt32(reader, dest)
		if err != nil {
			return err
		}
	case CompressedInt64Type:
		err = scanCompressedInt64(reader, dest)
		if err != nil {
			return err
		}
	case Float32ArrayType:
		err = scanFloat32Array(reader, dest)
		if err != nil {
			return err
		}
	case Float32Type:
		err = scanFloat32(reader, dest)
		if err != nil {
			return err
		}
	case Float64Type:
		err = scanFloat64(reader, dest)
		if err != nil {
			return err
		}
	case BooleanTrueType:
		dest.Num = 1
	case BooleanFalseType:
		dest.Num = 0
	case IntZeroType, ShortZeroType, LongZeroType, FloatZeroType, DoubleZeroType, ByteZeroType:
		break
	case ArrayType:
		err = scanArray(reader, dest)
		if err != nil {
			return err
		}
	case ShortArrayType:
		err = scanShortArray(reader, dest)
		if err != nil {
			return err
		}
	case ByteArrayType:
		err = scanByteArray(reader, dest)
		if err != nil {
			return err
		}
	case BooleanArrayType:
		err = scanBooleanArray(reader, dest)
		if err != nil {
			return err
		}
	case StringArrayType:
		err = scanStringArray(reader, dest)
		if err != nil {
			return err
		}
	case DictionaryType:
		err = scanDictionary(reader, dest)
		if err != nil {
			return err
		}
	case CompressedIntArrayType:
		err = scanCompressedIntArray(reader, dest)
		if err != nil {
			return err
		}
	case CompressedLongArrayType:
		err = scanCompressedLongArray(reader, dest)
		if err != nil {
			return err
		}
	case NilType, UnknownType:
		break
	default:
		return fmt.Errorf("unsupported type: %d", dest.Kind)
	}
	return nil
}
