package nes

import (
	"crypto"
	_ "crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/stellarisJAY/nesgo/nes/cartridge"
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type Save struct {
	Game       string // 存档对应的游戏
	Serializer string // 快照的序列化格式
	Hash       []byte // 存档checksum Hash 防止玩家修改存档数据
	Snapshot   []byte // 该存档的模拟器快照
}

func (e *RawEmulator) SaveToFile() error {
	timestamp := time.Now()
	saveData, err := e.GetSaveData()
	if err != nil {
		return err
	}
	path := filepath.Join(e.config.SaveDirectory, getSaveFileName(filepath.Base(e.config.Game), timestamp))
	if err := os.WriteFile(path, saveData, 0644); err != nil {
		return fmt.Errorf("write save file error %s", err)
	}
	log.Println("game saved at:", path)
	return nil
}

func (e *RawEmulator) GetSaveData() ([]byte, error) {
	s := e.createSnapshot()
	data, err := GetSnapshotSerializer(e.config.SnapshotSerializer).Serialize(s)
	if err != nil {
		return nil, fmt.Errorf("serializer error: %s", err)
	}
	save := createSave(e.config.Game, e.config.SnapshotSerializer, data)
	saveData, _ := json.Marshal(save)
	return saveData, nil
}

func (e *RawEmulator) Load(savedGame []byte) error {
	save := Save{}
	if err := json.Unmarshal(savedGame, &save); err != nil {
		return err
	}
	if verifyChecksum(save) {
		s, err := GetSnapshotSerializer(save.Serializer).Deserialize(save.Snapshot)
		if err != nil {
			return fmt.Errorf("invalid save data %s", err)
		}
		e.processor.Reverse(s.Processor)
		e.bus.Reverse(s.Bus)
		_ = e.ppu.Reverse(s.PPU)
		if err := cartridge.Load(e.cartridge, s.Cartridge); err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("verify checksum failed, corrupted save data")
	}
}

func createSave(game, format string, snapshot []byte) Save {
	h := crypto.SHA256.New()
	s := Save{
		Game:       game,
		Serializer: format,
		Snapshot:   snapshot,
	}
	h.Write(snapshot)
	h.Write([]byte(game))
	h.Write([]byte(format))
	s.Hash = h.Sum(nil)
	return s
}

func verifyChecksum(save Save) bool {
	h := crypto.SHA256.New()
	h.Write(save.Snapshot)
	h.Write([]byte(save.Game))
	h.Write([]byte(save.Serializer))
	return slices.Equal(save.Hash, h.Sum(nil))
}

func getSaveFileName(game string, timestamp time.Time) string {
	return game + "-" + timestamp.Format(time.DateTime)
}
