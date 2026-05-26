package v18

import "math"

func (p Parameter) StringValue() (string, bool) {
	if p.Kind != StringType {
		return "", false
	}
	return string(p.Blob), true
}

func (p Parameter) Float32Value() (float32, bool) {
	if p.Kind == FloatZeroType {
		return 0, true
	}
	if p.Kind != Float32Type {
		return 0, false
	}
	return math.Float32frombits(uint32(p.Num)), true
}

func (p Parameter) Float64Value() (float64, bool) {
	if p.Kind == DoubleZeroType {
		return 0, true
	}
	if p.Kind != Float64Type {
		return 0, false
	}
	return math.Float64frombits(uint64(p.Num)), true
}

func (p Parameter) IntValue() (int64, bool) {
	switch p.Kind {
	case Int8Type:
		return int64(int8(uint8(p.Num))), true
	case Int8Positive:
		return int64(uint8(p.Num)), true
	case Int8Negative:
		return int64(int32(p.Num)), true
	case Int16Type:
		return int64(int16(uint16(p.Num))), true
	case Int16Positive:
		return int64(uint16(p.Num)), true
	case Int16Negative, CompressedInt32Type:
		return int64(int32(uint32(p.Num))), true
	case Long8Positive:
		return int64(uint8(p.Num)), true
	case Long8Negative:
		return int64(p.Num), true
	case Long16Positive:
		return int64(uint16(p.Num)), true
	case Long16Negative, CompressedInt64Type:
		return int64(p.Num), true
	case IntZeroType, ShortZeroType, LongZeroType, ByteZeroType:
		return 0, true
	default:
		return 0, false
	}
}

func (p Parameter) BooleanValue() (bool, bool) {
	if p.Kind != BooleanTrueType && p.Kind != BooleanFalseType && p.Kind != BooleanType {
		return false, false
	}
	return p.Num == 1, true
}
