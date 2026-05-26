package photon_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AutoDruid/photon-parser"
	v18 "github.com/AutoDruid/photon-parser/internal/parameters/v18"
)

type WiresharkFrame struct {
	Source struct {
		Layers struct {
			Frame struct {
				Number string `json:"frame.number"`
			} `json:"frame"`
			UDP struct {
				Payload string `json:"udp.payload"`
			} `json:"udp"`
		} `json:"layers"`
	} `json:"_source"`
}

func rowToPayload(item json.RawMessage) (payload []byte, frameSuffix string, err error) {
	var asString string
	if err := json.Unmarshal(item, &asString); err == nil {
		b, err := hex.DecodeString(strings.ReplaceAll(asString, ":", ""))
		return b, "", err
	}

	var row WiresharkFrame
	if err := json.Unmarshal(item, &row); err != nil {
		return nil, "", err
	}
	if n := row.Source.Layers.Frame.Number; n != "" {
		frameSuffix = " [wireshark frame " + n + "]"
	}
	p := row.Source.Layers.UDP.Payload
	if p == "" {
		return nil, frameSuffix, nil
	}
	b, err := hex.DecodeString(strings.ReplaceAll(p, ":", ""))
	return b, frameSuffix, err
}

func loadCaptures(path string, t *testing.T) []json.RawMessage {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("test data not available: %v", err)
	}
	var frames []json.RawMessage
	if err := json.Unmarshal(data, &frames); err != nil {
		t.Fatalf("invalid test data: %v", err)
	}
	return frames
}

func loadCapturesB(path string, b *testing.B) []json.RawMessage {
	b.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		b.Skipf("test data not available: %v", err)
	}
	var frames []json.RawMessage
	if err := json.Unmarshal(data, &frames); err != nil {
		b.Fatalf("invalid test data: %v", err)
	}
	return frames
}

func TestParseVersion18OnDataset1(t *testing.T) {

	parser := photon.NewParserV18()

	frames := loadCaptures("./tests/dataset/v18/1.json", t)

	session := photon.Session[v18.Parameter]{}

	for i, frame := range frames {
		payload, frameNo, err := rowToPayload(frame)
		if err != nil {
			t.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) == 0 {
			t.Logf("row %d%s: skip (empty payload)", i, frameNo)
			continue
		}

		name := fmt.Sprintf("row_%d", i)
		if s := strings.TrimSpace(frameNo); s != "" {
			name += "_" + strings.ReplaceAll(strings.ReplaceAll(s, " ", "_"), "__", "_")
		}

		t.Run(name, func(t *testing.T) {
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				t.Fatalf("parse (%d bytes): %v", len(payload), err)
			}
			if session.Timestamp == 0 {
				t.Fatalf("session must have a timestamp")
			}
		})
	}
}

func TestParseVersion18OnDataset2(t *testing.T) {

	parser := photon.NewParserV18()

	frames := loadCaptures("./tests/dataset/v18/2.json", t)

	session := photon.Session[v18.Parameter]{}
	for i, frame := range frames {
		payload, frameNo, err := rowToPayload(frame)
		if err != nil {
			t.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) == 0 {
			t.Logf("row %d%s: skip (empty payload)", i, frameNo)
			continue
		}

		name := fmt.Sprintf("row_%d", i)
		if s := strings.TrimSpace(frameNo); s != "" {
			name += "_" + strings.ReplaceAll(strings.ReplaceAll(s, " ", "_"), "__", "_")
		}

		t.Run(name, func(t *testing.T) {
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				t.Fatalf("parse (%d bytes): %v", len(payload), err)
			}
			if session.Timestamp == 0 {
				t.Fatalf("session must have a timestamp")
			}
		})
	}
}

func TestParseVersion18OnDataset3(t *testing.T) {

	parser := photon.NewParserV18()

	frames := loadCaptures("./tests/dataset/v18/3.json", t)

	session := photon.Session[v18.Parameter]{}
	for i, frame := range frames {
		payload, frameNo, err := rowToPayload(frame)
		if err != nil {
			t.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) == 0 {
			t.Logf("row %d%s: skip (empty payload)", i, frameNo)
			continue
		}

		name := fmt.Sprintf("row_%d", i)
		if s := strings.TrimSpace(frameNo); s != "" {
			name += "_" + strings.ReplaceAll(strings.ReplaceAll(s, " ", "_"), "__", "_")
		}

		t.Run(name, func(t *testing.T) {
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				t.Fatalf("parse (%d bytes): %v", len(payload), err)
			}
			if session.Timestamp == 0 {
				t.Fatalf("session must have a timestamp")
			}
		})
	}
}

func TestParseVersion18OnDataset3WithoutUnknownPayloads(t *testing.T) {

	parser := photon.NewParserV18(photon.SkipUnknownPayloads(true))

	frames := loadCaptures("./tests/dataset/v18/3.json", t)

	session := photon.Session[v18.Parameter]{}

	parser.OnCommandSync(func(command photon.Command[v18.Parameter]) {
		if !bytes.Equal(command.UnknownPayload.Raw, []byte{}) {
			t.Fatalf("unknown payload found")
		}
	})

	for i, frame := range frames {
		payload, frameNo, err := rowToPayload(frame)
		if err != nil {
			t.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) == 0 {
			t.Logf("row %d%s: skip (empty payload)", i, frameNo)
			continue
		}

		name := fmt.Sprintf("row_%d", i)
		if s := strings.TrimSpace(frameNo); s != "" {
			name += "_" + strings.ReplaceAll(strings.ReplaceAll(s, " ", "_"), "__", "_")
		}

		t.Run(name, func(t *testing.T) {
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				t.Fatalf("parse (%d bytes): %v", len(payload), err)
			}
			if session.Timestamp == 0 {
				t.Fatalf("session must have a timestamp")
			}
		})
	}
}

func TestParseVersion18OnDataset3WithoutCommands(t *testing.T) {

	parser := photon.NewParserV18(photon.SkipCommands(photon.SendReliableCommand))

	frames := loadCaptures("./tests/dataset/v18/3.json", t)

	session := photon.Session[v18.Parameter]{}

	parser.OnCommandSync(func(command photon.Command[v18.Parameter]) {
		if command.Type == photon.SendReliableCommand {
			t.Fatalf("send reliable command found")
		}
	})

	for i, frame := range frames {
		payload, frameNo, err := rowToPayload(frame)
		if err != nil {
			t.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) == 0 {
			t.Logf("row %d%s: skip (empty payload)", i, frameNo)
			continue
		}

		name := fmt.Sprintf("row_%d", i)
		if s := strings.TrimSpace(frameNo); s != "" {
			name += "_" + strings.ReplaceAll(strings.ReplaceAll(s, " ", "_"), "__", "_")
		}

		t.Run(name, func(t *testing.T) {
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				t.Fatalf("parse (%d bytes): %v", len(payload), err)
			}
			if session.Timestamp == 0 {
				t.Fatalf("session must have a timestamp")
			}
		})
	}
}

func TestParseVersion18OnDataset3WithoutParameters(t *testing.T) {

	parser := photon.NewParserV18(photon.SkipParameterParsing(true))

	frames := loadCaptures("./tests/dataset/v18/3.json", t)

	session := photon.Session[v18.Parameter]{}

	parser.OnCommandSync(func(command photon.Command[v18.Parameter]) {
		var parameters []v18.Parameter
		switch command.Type {
		case photon.SendReliableCommand:
			parameters = command.ReliablePayload.Parameters
		case photon.SendUnreliableCommand:
			parameters = command.UnreliablePayload.Parameters
		}
		if len(parameters) > 0 {
			t.Fatalf("parameters found")
		}
	})

	for i, frame := range frames {
		payload, frameNo, err := rowToPayload(frame)
		if err != nil {
			t.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) == 0 {
			t.Logf("row %d%s: skip (empty payload)", i, frameNo)
			continue
		}

		name := fmt.Sprintf("row_%d", i)
		if s := strings.TrimSpace(frameNo); s != "" {
			name += "_" + strings.ReplaceAll(strings.ReplaceAll(s, " ", "_"), "__", "_")
		}

		t.Run(name, func(t *testing.T) {
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				t.Fatalf("parse (%d bytes): %v", len(payload), err)
			}
			if session.Timestamp == 0 {
				t.Fatalf("session must have a timestamp")
			}
		})
	}
}

func TestParseVersion18OnDataset3WithoutEvents(t *testing.T) {

	parser := photon.NewParserV18(photon.SkipTargetEventCodes(photon.OperationRequest, photon.OperationResponse))

	frames := loadCaptures("./tests/dataset/v18/3.json", t)

	session := photon.Session[v18.Parameter]{}

	for i, frame := range frames {
		payload, frameNo, err := rowToPayload(frame)
		if err != nil {
			t.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) == 0 {
			t.Logf("row %d%s: skip (empty payload)", i, frameNo)
			continue
		}

		name := fmt.Sprintf("row_%d", i)
		if s := strings.TrimSpace(frameNo); s != "" {
			name += "_" + strings.ReplaceAll(strings.ReplaceAll(s, " ", "_"), "__", "_")
		}

		t.Run(name, func(t *testing.T) {
			parser.OnCommandSync(func(command photon.Command[v18.Parameter]) {
				switch command.Type {
				case photon.SendReliableCommand:
					switch command.ReliablePayload.Type {
					case photon.OperationRequest, photon.OperationResponse:
						if len(command.ReliablePayload.Parameters) > 0 {
							t.Fatalf("decoded parameters on skipped operation type %v", command.ReliablePayload.Type)
						}
					}
				case photon.SendUnreliableCommand:
					switch command.UnreliablePayload.Type {
					case photon.OperationRequest, photon.OperationResponse:
						if len(command.UnreliablePayload.Parameters) > 0 {
							t.Fatalf("decoded parameters on skipped operation type %v", command.UnreliablePayload.Type)
						}
					}
				}
			})

			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				t.Fatalf("parse (%d bytes): %v", len(payload), err)
			}
			if session.Timestamp == 0 {
				t.Fatalf("session must have a timestamp")
			}
		})
	}
}

func TestParserOnASinglePacketVersion18(t *testing.T) {
	parser := photon.NewParserV18()
	frames := loadCaptures("./tests/dataset/v18/2.json", t)

	payload, _, err := rowToPayload(frames[181])
	if err != nil {
		t.Fatalf("decode row: %v", err)
	}

	session := photon.Session[v18.Parameter]{}

	err = parser.ParsePacketInto(payload, &session)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if session.Timestamp == 0 {
		t.Fatalf("session must have a timestamp")
	}
}

func TestKnownCommandPayloadRecoveryDoesNotOverSkipNextCommand(t *testing.T) {
	payload := []byte{
		0x00, 0x00, 0x00, 0x02, // session: peer ID, CRC disabled, 2 commands
		0x00, 0x00, 0x00, 0x01, // timestamp
		0x00, 0x00, 0x00, 0x00, // challenge

		0x06, 0x00, 0x00, 0x00, // send reliable header prefix
		0x00, 0x00, 0x00, 0x12, // command length: 12-byte header + 6-byte payload
		0x00, 0x00, 0x00, 0x01, // reliable sequence number
		0xf3, 0x02, 0x01, 0x01, // operation request with one parameter
		0x00, 0xff, // unsupported parameter type consumes the full payload before failing

		0x05, 0xff, 0x01, 0x04, // ping command header prefix
		0x00, 0x00, 0x00, 0x0c, // command length: header only
		0x00, 0x00, 0x00, 0x02, // reliable sequence number
	}

	parser := photon.NewParserV18(photon.SkipUnknownPayloads(true))
	var session photon.SessionV18

	if err := parser.ParsePacketInto(payload, &session); err != nil {
		t.Fatalf("ParsePacketInto() error = %v", err)
	}
	if got := session.Commands[1].Type; got != photon.PingCommand {
		t.Fatalf("second command type = %v, want %v", got, photon.PingCommand)
	}
}

func BenchmarkParserOn343PacketsVersion18(b *testing.B) {
	parser := photon.NewParserV18()
	frames := loadCapturesB("./tests/dataset/v18/1.json", b)

	var totalBytes int64
	payloads := make([][]byte, 0, len(frames))

	for i, frame := range frames {
		payload, _, err := rowToPayload(frame)
		if err != nil {
			b.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) > 0 {
			payloads = append(payloads, payload)
			totalBytes += int64(len(payload))
		}
	}

	session := photon.Session[v18.Parameter]{}

	b.ReportAllocs()
	b.SetBytes(totalBytes)
	b.ResetTimer()

	for b.Loop() {
		for i, payload := range payloads {
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				b.Fatalf("row %d: error: %v", i, err)
			}
		}
	}
}

func BenchmarkParserOn343PacketsVersion18WithSyncHooks(b *testing.B) {
	parser := photon.NewParserV18()
	frames := loadCapturesB("./tests/dataset/v18/1.json", b)

	var totalBytes int64
	payloads := make([][]byte, 0, len(frames))

	for i, frame := range frames {
		payload, _, err := rowToPayload(frame)
		if err != nil {
			b.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) > 0 {
			payloads = append(payloads, payload)
			totalBytes += int64(len(payload))
		}
	}

	session := photon.Session[v18.Parameter]{}

	parser.OnSessionSync(func(session photon.Session[v18.Parameter]) {
	})

	b.ReportAllocs()
	b.SetBytes(totalBytes)
	b.ResetTimer()

	for b.Loop() {
		for i, payload := range payloads {
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				b.Fatalf("row %d: error: %v", i, err)
			}
		}
	}
}

func BenchmarkParserOn343PacketsVersion18WithAsyncHooks(b *testing.B) {
	parser := photon.NewParserV18()
	frames := loadCapturesB("./tests/dataset/v18/1.json", b)

	var totalBytes int64
	payloads := make([][]byte, 0, len(frames))

	for i, frame := range frames {
		payload, _, err := rowToPayload(frame)
		if err != nil {
			b.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) > 0 {
			payloads = append(payloads, payload)
			totalBytes += int64(len(payload))
		}
	}

	session := photon.Session[v18.Parameter]{}

	ch := parser.OnSessionAsync(photon.HookOptions{Size: uint16(b.N)})

	b.ReportAllocs()
	b.SetBytes(totalBytes)
	b.ResetTimer()

	for b.Loop() {
		for i, payload := range payloads {
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				b.Fatalf("row %d: error: %v", i, err)
			}
		}
	}

	b.StopTimer()
	parser.Close()
	for range ch {
		// drain channel
	}
}

func BenchmarkParserOn343PacketsVersion18WithoutUnknownPayloads(b *testing.B) {
	parser := photon.NewParserV18(photon.SkipUnknownPayloads(true))
	frames := loadCapturesB("./tests/dataset/v18/1.json", b)

	var totalBytes int64
	payloads := make([][]byte, 0, len(frames))

	for i, frame := range frames {
		payload, _, err := rowToPayload(frame)
		if err != nil {
			b.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) > 0 {
			payloads = append(payloads, payload)
			totalBytes += int64(len(payload))
		}
	}

	session := photon.Session[v18.Parameter]{}

	b.ReportAllocs()
	b.SetBytes(totalBytes)
	b.ResetTimer()

	for b.Loop() {
		for i, payload := range payloads {
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				b.Fatalf("row %d: error: %v", i, err)
			}
		}
	}
}

func BenchmarkParserOn343PacketsVersion18WithoutParameters(b *testing.B) {
	parser := photon.NewParserV18(photon.SkipParameterParsing(true))
	frames := loadCapturesB("./tests/dataset/v18/1.json", b)

	var totalBytes int64
	payloads := make([][]byte, 0, len(frames))

	for i, frame := range frames {
		payload, _, err := rowToPayload(frame)
		if err != nil {
			b.Fatalf("row %d: decode row: %v", i, err)
		}
		if len(payload) > 0 {
			payloads = append(payloads, payload)
			totalBytes += int64(len(payload))
		}
	}

	session := photon.Session[v18.Parameter]{}

	b.ReportAllocs()
	b.SetBytes(totalBytes)
	b.ResetTimer()

	for b.Loop() {
		for i, payload := range payloads {
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				b.Fatalf("row %d: error: %v", i, err)
			}
		}
	}
}

func BenchmarkParserOnASinglePacketVersion18(b *testing.B) {
	parser := photon.NewParserV18()
	frames := loadCapturesB("./tests/dataset/v18/1.json", b)

	payload, _, err := rowToPayload(frames[200])
	if err != nil {
		b.Fatalf("decode row: %v", err)
	}

	totalBytes := int64(len(payload))

	session := photon.Session[v18.Parameter]{}

	b.ReportAllocs()
	b.SetBytes(totalBytes)
	b.ResetTimer()

	for b.Loop() {
		err := parser.ParsePacketInto(payload, &session)
		if err != nil {
			b.Fatalf("error: %v", err)
		}
		if session.Timestamp == 0 {
			b.Fatalf("session must have a timestamp")
		}
	}
}

func BenchmarkParserOnASingleWithAsyncHooksPacketVersion18(b *testing.B) {
	parser := photon.NewParserV18()
	frames := loadCapturesB("./tests/dataset/v18/1.json", b)

	payload, _, err := rowToPayload(frames[200])
	if err != nil {
		b.Fatalf("decode row: %v", err)
	}

	totalBytes := int64(len(payload))

	ch := parser.OnSessionAsync(photon.HookOptions{Size: uint16(b.N)})

	session := photon.Session[v18.Parameter]{}

	b.ReportAllocs()
	b.SetBytes(totalBytes)
	b.ResetTimer()

	for b.Loop() {
		err := parser.ParsePacketInto(payload, &session)
		if err != nil {
			b.Fatalf("error: %v", err)
		}
		if session.Timestamp == 0 {
			b.Fatalf("session must have a timestamp")
		}
	}

	b.StopTimer()
	parser.Close()
	for range ch {
		// drain channel
	}
}

func BenchmarkParserOnASingleWithSyncHooksPacketVersion18(b *testing.B) {
	parser := photon.NewParserV18()
	frames := loadCapturesB("./tests/dataset/v18/1.json", b)

	payload, _, err := rowToPayload(frames[200])
	if err != nil {
		b.Fatalf("decode row: %v", err)
	}

	totalBytes := int64(len(payload))

	parser.OnSessionSync(func(session photon.Session[v18.Parameter]) {
	})

	session := photon.Session[v18.Parameter]{}

	b.ReportAllocs()
	b.SetBytes(totalBytes)
	b.ResetTimer()

	for b.Loop() {
		err := parser.ParsePacketInto(payload, &session)
		if err != nil {
			b.Fatalf("error: %v", err)
		}
		if session.Timestamp == 0 {
			b.Fatalf("session must have a timestamp")
		}
	}

	b.StopTimer()
}

func BenchmarkParserOnASinglePacketVersion18Parallel(b *testing.B) {
	frames := loadCapturesB("./tests/dataset/v18/1.json", b)
	payload, _, err := rowToPayload(frames[6])
	if err != nil {
		b.Fatalf("decode row: %v", err)
	}

	totalBytes := int64(len(payload))

	b.ReportAllocs()
	b.SetBytes(totalBytes)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		parser := photon.NewParserV18()
		session := photon.Session[v18.Parameter]{}

		for pb.Next() {
			err := parser.ParsePacketInto(payload, &session)
			if err != nil {
				b.Errorf("error: %v", err)
			}
			if session.PeerID == 0 {
				b.Errorf("session must have a peer id")
			}
		}
	})
}
