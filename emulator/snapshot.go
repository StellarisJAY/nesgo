package emulator

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/stellarisJAY/nesgo/bus"
	"github.com/stellarisJAY/nesgo/cpu"
	"github.com/stellarisJAY/nesgo/ppu"
	"time"
)

type Snapshot struct {
	Processor cpu.Snapshot
	PPU       ppu.Snapshot
	Bus       bus.Snapshot
	Timestamp time.Time
}

type SnapshotSerializer interface {
	Serialize(Snapshot) ([]byte, error)
	Deserialize([]byte) (Snapshot, error)
}

var GlobalSerializers = make(map[string]SnapshotSerializer)

func init() {
	GlobalSerializers["json"] = JSONSnapshotSerializer{}
	GlobalSerializers["gob"] = GobSnapshotSerializer{}
}

func (e *RawEmulator) PushSnapshot() {
	if time.Now().After(e.lastSnapshotTime.Add(snapshotInterval)) {
		s := e.createSnapshot()
		e.m.Lock()
		defer e.m.Unlock()
		e.snapshots = append(e.snapshots, s)
		if len(e.snapshots) > maxSnapshots {
			e.snapshots = e.snapshots[1:]
		}
	}
}

func (e *RawEmulator) createSnapshot() Snapshot {
	e.lastSnapshotTime = time.Now()
	// 必须保证每个组件MakeSnapshot时没有线程安全问题
	return Snapshot{
		Processor: e.processor.MakeSnapshot(),
		PPU:       e.ppu.MakeSnapshot(),
		Bus:       e.bus.MakeSnapshot(),
		Timestamp: e.lastSnapshotTime,
	}
}

type GobSnapshotSerializer struct{}

func (gs GobSnapshotSerializer) Serialize(s Snapshot) ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(s)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (gs GobSnapshotSerializer) Deserialize(data []byte) (Snapshot, error) {
	s := Snapshot{}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&s)
	return s, err
}

type JSONSnapshotSerializer struct{}

func (js JSONSnapshotSerializer) Serialize(s Snapshot) ([]byte, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (js JSONSnapshotSerializer) Deserialize(data []byte) (Snapshot, error) {
	s := Snapshot{}
	err := json.Unmarshal(data, &s)
	return s, err
}

func GetSnapshotSerializer(format string) SnapshotSerializer {
	if s, ok := GlobalSerializers[format]; ok {
		return s
	}
	panic("no such serializer format")
}
