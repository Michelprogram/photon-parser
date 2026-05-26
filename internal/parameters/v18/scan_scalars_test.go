package v18_test

import (
	"bytes"
	"testing"

	v18 "github.com/AutoDruid/photon-parser/internal/parameters/v18"
	"github.com/AutoDruid/photon-parser/internal/reader"
)

func TestParseStringParameterAndAccessor(t *testing.T) {
	longString := bytes.Repeat([]byte{'a'}, 128)
	tests := []struct {
		name       string
		input      []byte
		wantID     uint8
		wantString string
		wantNum    uint64
		wantCursor int
	}{
		{
			name:       "string with 1 byte varint size",
			input:      append([]byte{0x01, byte(v18.StringType), 0x05}, []byte("hello")...),
			wantID:     1,
			wantString: "hello",
			wantNum:    5,
			wantCursor: 8,
		},
		{
			name:       "string with 2 byte varint size",
			input:      append([]byte{0x02, byte(v18.StringType), 0x80, 0x01}, longString...),
			wantID:     2,
			wantString: string(longString),
			wantNum:    128,
			wantCursor: 132,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v18.Parameter
			var got v18.Parameter
			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}
			if got.Kind != v18.StringType {
				t.Errorf("Kind = %d, want %d", got.Kind, v18.StringType)
			}
			if got.Num != tt.wantNum {
				t.Errorf("Num = %d, want %d", got.Num, tt.wantNum)
			}
			if gotString, ok := got.StringValue(); !ok || gotString != tt.wantString {
				t.Errorf("StringValue() = %q, want %q", gotString, tt.wantString)
			}
			if r.Cursor != tt.wantCursor {
				t.Errorf("Cursor = %d, want %d", r.Cursor, tt.wantCursor)
			}
		})
	}
}

func TestParseFloat32ParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		wantID     uint8
		wantValue  float32
		wantCursor int
	}{
		{
			name:       "positive float32",
			input:      []byte{0x01, byte(v18.Float32Type), 0x00, 0x00, 0xf7, 0x42},
			wantID:     1,
			wantValue:  123.5,
			wantCursor: 6,
		},
		{
			name:       "negative float32",
			input:      []byte{0x02, byte(v18.Float32Type), 0x00, 0x00, 0xf7, 0xc2},
			wantID:     2,
			wantValue:  -123.5,
			wantCursor: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}

			if got.Kind != v18.Float32Type {
				t.Errorf("Kind = %d, want %d", got.Kind, v18.Float32Type)
			}

			if value, ok := got.Float32Value(); !ok || value != tt.wantValue {
				t.Errorf("Float32Value() = %v, want %v", value, tt.wantValue)
			}

			if r.Cursor != tt.wantCursor {
				t.Errorf("Cursor = %d, want %d", r.Cursor, tt.wantCursor)
			}
		})
	}
}

func TestParseZeroFloatParametersAndAccessors(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		wantID     uint8
		wantKind   v18.ParameterType
		wantCursor int
	}{
		{
			name:       "float32 zero shorthand",
			input:      []byte{0x01, byte(v18.FloatZeroType)},
			wantID:     1,
			wantKind:   v18.FloatZeroType,
			wantCursor: 2,
		},
		{
			name:       "float64 zero shorthand",
			input:      []byte{0x02, byte(v18.DoubleZeroType)},
			wantID:     2,
			wantKind:   v18.DoubleZeroType,
			wantCursor: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)
			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}
			if got.Kind != tt.wantKind {
				t.Errorf("Kind = %d, want %d", got.Kind, tt.wantKind)
			}
			if r.Cursor != tt.wantCursor {
				t.Errorf("Cursor = %d, want %d", r.Cursor, tt.wantCursor)
			}
			if tt.wantKind == v18.FloatZeroType {
				if value, ok := got.Float32Value(); !ok || value != 0 {
					t.Errorf("Float32Value() = %v, %v; want 0, true", value, ok)
				}
			} else {
				if value, ok := got.Float64Value(); !ok || value != 0 {
					t.Errorf("Float64Value() = %v, %v; want 0, true", value, ok)
				}
			}
		})
	}
}

func TestParseInt8ParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		wantID     uint8
		wantKind   v18.ParameterType
		wantValue  int64
		wantCursor int
	}{
		{
			name:       "int8 raw value",
			input:      []byte{0x01, byte(v18.Int8Type), 0x7f},
			wantID:     1,
			wantKind:   v18.Int8Type,
			wantValue:  127,
			wantCursor: 3,
		},
		{
			name:       "int8 positive",
			input:      []byte{0x02, byte(v18.Int8Positive), 0x7f},
			wantID:     2,
			wantKind:   v18.Int8Positive,
			wantValue:  127,
			wantCursor: 3,
		},
		{
			name:       "int8 negative",
			input:      []byte{0x03, byte(v18.Int8Negative), 0x7f},
			wantID:     3,
			wantKind:   v18.Int8Negative,
			wantValue:  -127,
			wantCursor: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}

			if got.Kind != tt.wantKind {
				t.Errorf("Kind = %d, want %d", got.Kind, tt.wantKind)
			}

			if value, ok := got.IntValue(); !ok || value != tt.wantValue {
				t.Errorf("IntValue() = %d, want %d", value, tt.wantValue)
			}

			if r.Cursor != tt.wantCursor {
				t.Errorf("Cursor = %d, want %d", r.Cursor, tt.wantCursor)
			}
		})
	}
}

func TestParseInt16ParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		wantID     uint8
		wantKind   v18.ParameterType
		wantValue  int64
		wantCursor int
	}{
		{
			name:       "int16 raw positive value",
			input:      []byte{0x01, byte(v18.Int16Type), 0x39, 0x30},
			wantID:     1,
			wantKind:   v18.Int16Type,
			wantValue:  12345,
			wantCursor: 4,
		},
		{
			name:       "int16 raw negative value",
			input:      []byte{0x02, byte(v18.Int16Type), 0xc7, 0xcf},
			wantID:     2,
			wantKind:   v18.Int16Type,
			wantValue:  -12345,
			wantCursor: 4,
		},
		{
			name:       "int16 positive shorthand",
			input:      []byte{0x03, byte(v18.Int16Positive), 0x39, 0x30},
			wantID:     3,
			wantKind:   v18.Int16Positive,
			wantValue:  12345,
			wantCursor: 4,
		},
		{
			name:       "int16 negative shorthand",
			input:      []byte{0x04, byte(v18.Int16Negative), 0x39, 0x30},
			wantID:     4,
			wantKind:   v18.Int16Negative,
			wantValue:  -12345,
			wantCursor: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}

			if got.Kind != tt.wantKind {
				t.Errorf("Kind = %d, want %d", got.Kind, tt.wantKind)
			}

			if value, ok := got.IntValue(); !ok || value != tt.wantValue {
				t.Errorf("IntValue() = %d, want %d", value, tt.wantValue)
			}

			if r.Cursor != tt.wantCursor {
				t.Errorf("Cursor = %d, want %d", r.Cursor, tt.wantCursor)
			}
		})
	}
}

func TestParseLong8AndLong16ParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		wantID     uint8
		wantKind   v18.ParameterType
		wantValue  int64
		wantCursor int
	}{
		{
			name:       "long8 positive",
			input:      []byte{0x01, byte(v18.Long8Positive), 0x7f},
			wantID:     1,
			wantKind:   v18.Long8Positive,
			wantValue:  127,
			wantCursor: 3,
		},
		{
			name:       "long8 negative",
			input:      []byte{0x02, byte(v18.Long8Negative), 0x7f},
			wantID:     2,
			wantKind:   v18.Long8Negative,
			wantValue:  -127,
			wantCursor: 3,
		},
		{
			name:       "long16 positive",
			input:      []byte{0x03, byte(v18.Long16Positive), 0x39, 0x30},
			wantID:     3,
			wantKind:   v18.Long16Positive,
			wantValue:  12345,
			wantCursor: 4,
		},
		{
			name:       "long16 negative",
			input:      []byte{0x04, byte(v18.Long16Negative), 0x39, 0x30},
			wantID:     4,
			wantKind:   v18.Long16Negative,
			wantValue:  -12345,
			wantCursor: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}

			if got.Kind != tt.wantKind {
				t.Errorf("Kind = %d, want %d", got.Kind, tt.wantKind)
			}

			if value, ok := got.IntValue(); !ok || value != tt.wantValue {
				t.Errorf("IntValue() = %d, want %d", value, tt.wantValue)
			}

			if r.Cursor != tt.wantCursor {
				t.Errorf("Cursor = %d, want %d", r.Cursor, tt.wantCursor)
			}
		})
	}
}

func TestParseCompressedInt32AndInt64ParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		wantID     uint8
		wantKind   v18.ParameterType
		wantValue  int64
		wantCursor int
	}{
		{
			name:       "compressed int32 3 byte positive",
			input:      []byte{0x01, byte(v18.CompressedInt32Type), 0x80, 0x80, 0x01},
			wantID:     1,
			wantKind:   v18.CompressedInt32Type,
			wantValue:  8192,
			wantCursor: 5,
		},
		{
			name:       "compressed int32 3 byte negative",
			input:      []byte{0x02, byte(v18.CompressedInt32Type), 0x81, 0x80, 0x01},
			wantID:     2,
			wantKind:   v18.CompressedInt32Type,
			wantValue:  -8193,
			wantCursor: 5,
		},
		{
			name:       "compressed int64 4 byte positive",
			input:      []byte{0x03, byte(v18.CompressedInt64Type), 0x80, 0x80, 0x80, 0x01},
			wantID:     3,
			wantKind:   v18.CompressedInt64Type,
			wantValue:  1048576,
			wantCursor: 6,
		},
		{
			name:       "compressed int64 4 byte negative",
			input:      []byte{0x04, byte(v18.CompressedInt64Type), 0x81, 0x80, 0x80, 0x01},
			wantID:     4,
			wantKind:   v18.CompressedInt64Type,
			wantValue:  -1048577,
			wantCursor: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}

			if got.Kind != tt.wantKind {
				t.Errorf("Kind = %d, want %d", got.Kind, tt.wantKind)
			}

			if value, ok := got.IntValue(); !ok || value != tt.wantValue {
				t.Errorf("IntValue() = %d, want %d", value, tt.wantValue)
			}

			if r.Cursor != tt.wantCursor {
				t.Errorf("Cursor = %d, want %d", r.Cursor, tt.wantCursor)
			}
		})
	}
}

func TestParseZeroTypeParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		wantID    uint8
		wantKind  v18.ParameterType
		wantValue int64
	}{
		{
			name:      "int zero",
			input:     []byte{0x01, byte(v18.IntZeroType)},
			wantID:    1,
			wantKind:  v18.IntZeroType,
			wantValue: 0,
		},
		{
			name:      "short zero",
			input:     []byte{0x02, byte(v18.ShortZeroType)},
			wantID:    2,
			wantKind:  v18.ShortZeroType,
			wantValue: 0,
		},
		{
			name:      "long zero",
			input:     []byte{0x03, byte(v18.LongZeroType)},
			wantID:    3,
			wantKind:  v18.LongZeroType,
			wantValue: 0,
		},
		{
			name:      "byte zero",
			input:     []byte{0x04, byte(v18.ByteZeroType)},
			wantID:    4,
			wantKind:  v18.ByteZeroType,
			wantValue: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}

			if got.Kind != tt.wantKind {
				t.Errorf("Kind = %d, want %d", got.Kind, tt.wantKind)
			}

			if value, ok := got.IntValue(); !ok || value != tt.wantValue {
				t.Errorf("IntValue() = %d, want %d", value, tt.wantValue)
			}
		})
	}
}

func TestParseBooleanParameterAndAccessor(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		wantID    uint8
		wantKind  v18.ParameterType
		wantValue bool
	}{
		{
			name:      "boolean true",
			input:     []byte{0x01, byte(v18.BooleanTrueType)},
			wantID:    1,
			wantKind:  v18.BooleanTrueType,
			wantValue: true,
		},
		{
			name:      "boolean false",
			input:     []byte{0x02, byte(v18.BooleanFalseType)},
			wantID:    2,
			wantKind:  v18.BooleanFalseType,
			wantValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.NewReader(tt.input)

			var parser v18.Parameter
			var got v18.Parameter

			if err := parser.ParseInto(r, nil, &got); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got.ID() != tt.wantID {
				t.Errorf("ID() = %d, want %d", got.ID(), tt.wantID)
			}

			if got.Kind != tt.wantKind {
				t.Errorf("Kind = %d, want %d", got.Kind, tt.wantKind)
			}

			if value, ok := got.BooleanValue(); !ok || value != tt.wantValue {
				t.Errorf("BooleanValue() = %t, want %t", value, tt.wantValue)
			}
		})
	}
}
