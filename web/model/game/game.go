package game

import (
	"github.com/stellarisJAY/nesgo/cartridge"
	"github.com/stellarisJAY/nesgo/config"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type CartridgeInfo struct {
	Name      string `json:"name"`
	Mapper    byte   `json:"mapper"`
	Mirroring string `json:"mirroring"`
	Size      int    `json:"size"`
}

func GetGameInfo(name string) (*CartridgeInfo, error) {
	conf := config.GetEmulatorConfig()
	dir := conf.GameDirectory
	path := filepath.Join(dir, name)
	file, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	info, err := cartridge.ParseCartridgeInfo(data)
	if err != nil {
		return nil, err
	}
	c := toCartridgeInfo(info)
	c.Size = len(data)
	return c, nil
}

func ListGames() ([]*CartridgeInfo, error) {
	conf := config.GetEmulatorConfig()
	dir := conf.GameDirectory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	games := make([]*CartridgeInfo, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".nes") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		info, err := cartridge.ParseCartridgeInfo(data)
		if err != nil {
			continue
		}
		c := toCartridgeInfo(info)
		c.Size = len(data)
		c.Name = entry.Name()
		games = append(games, c)
	}
	// order by name
	slices.SortStableFunc(games, func(a, b *CartridgeInfo) int {
		return strings.Compare(a.Name, b.Name)
	})
	return games, nil
}

func toCartridgeInfo(info *cartridge.Info) *CartridgeInfo {
	c := &CartridgeInfo{
		Name:   info.Name,
		Mapper: info.Mapper,
	}
	switch info.Mirroring {
	case cartridge.Vertical:
		c.Mirroring = "vertical"
	case cartridge.Horizontal:
		c.Mirroring = "horizontal"
	case cartridge.FourScreen:
		c.Mirroring = "four-screen"
	case cartridge.OneScreenLow, cartridge.OneScreenHigh:
		c.Mirroring = "one-screen"
	default:
		c.Mirroring = "unknown"
	}
	return c
}
