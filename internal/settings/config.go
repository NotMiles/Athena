/* Athena - A server for Attorney Online 2 written in Go
Copyright (C) 2022 MangosArentLiterature <mango@transmenace.dev>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>. */

// Package settings handles reading and writing to Athena's configuration files.
package settings

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/MangosArentLiterature/Athena/internal/area"
	"github.com/MangosArentLiterature/Athena/internal/permissions"
)

// Stores the path to the config directory
var ConfigPath string

type Config struct {
	ServerConfig `toml:"Server"`
	MSConfig     `toml:"MasterServer"`
}

type ServerConfig struct {
	Addr       string `toml:"addr"`
	Port       int    `toml:"port"`
	Name       string `toml:"name"`
	Desc       string `toml:"description"`
	MaxPlayers int    `toml:"max_players"`
	MaxMsg     int    `toml:"max_message_length"`
	BufSize    int    `toml:"log_buffer_size"`
}
type MSConfig struct {
	Advertise bool   `toml:"advertise"`
	MSAddr    string `toml:"addr"`
}

// Returns a default configuration.
func defaultConfig() *Config {
	return &Config{
		ServerConfig{
			Addr:       "",
			Port:       27016,
			Name:       "Unnamed Server",
			Desc:       "",
			MaxPlayers: 100,
			MaxMsg:     256,
			BufSize:    150,
		},
		MSConfig{
			Advertise: false,
			MSAddr:    "https://servers.aceattorneyonline.com/servers",
		},
	}
}

// Loads the configuation from config/config.toml.
func (conf *Config) Load() error {
	_, err := toml.DecodeFile(ConfigPath+"/config.toml", conf)
	if err != nil {
		return err
	}
	return nil
}

// Returns a loaded configuration
func GetConfig() (*Config, error) {
	conf := defaultConfig()
	err := conf.Load()

	if err != nil {
		return nil, err
	}

	return conf, nil
}

// Loads the music list from config/music.txt.
func LoadMusic() ([]string, error) {
	var musicList []string
	f, err := os.Open(ConfigPath + "/music.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewScanner(f)
	for in.Scan() {
		musicList = append(musicList, in.Text())
	}
	if len(musicList) == 0 {
		return nil, fmt.Errorf("empty musiclist")
	}
	if strings.ContainsRune(musicList[0], '.') {
		musicList = append([]string{"Songs"}, musicList...)
	}
	return musicList, nil
}

// Loads the character list from config/characters.txt.
func LoadCharacters() ([]string, error) {
	var charList []string
	f, err := os.Open(ConfigPath + "/characters.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewScanner(f)
	for in.Scan() {
		charList = append(charList, in.Text())
	}
	if len(charList) == 0 {
		return nil, fmt.Errorf("empty charlist")
	}
	return charList, nil
}

// Loads the area list from config/areas.toml.
func LoadAreas() ([]area.AreaData, error) {
	var conf struct {
		Area []area.AreaData
	}
	_, err := toml.DecodeFile(ConfigPath+"/areas.toml", &conf)
	if err != nil {
		return conf.Area, err
	}
	if len(conf.Area) == 0 {
		return conf.Area, fmt.Errorf("empty arealist")
	}
	return conf.Area, nil
}

func LoadRoles() ([]permissions.Role, error) {
	var conf struct {
		Role []permissions.Role
	}
	_, err := toml.DecodeFile(ConfigPath+"/roles.toml", &conf)
	if err != nil {
		return conf.Role, err
	}
	if len(conf.Role) == 0 {
		return conf.Role, fmt.Errorf("empty rolelist")
	}
	return conf.Role, nil
}
