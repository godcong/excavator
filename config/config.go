package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/godcong/fate/config"
)

const JSONName = "config.json"

const (
	ActionNon      string = ""
	ActionSimplify string = "简化"
	ActionVariant  string = "整理"
	ActionParse    string = "解析"
	ActionFill     string = "修复"
	ActionGrab     string = "获取"
)

type Config struct {
	TmpDir       string                `json:"tmp_dir"`
	BaseUrl      string                `json:"base_url"`
	UnicodeFile  string                `json:"unicode_file"`
	Action       string                `json:"action"`
	DatabaseExc  config.DatabaseConfig `json:"database_exc"`
	DatabaseFate config.DatabaseConfig `json:"database_fate"`
}

var DefaultJSONPath = ""

func init() {
	if DefaultJSONPath == "" {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		s, err := filepath.Abs(dir)
		if err != nil {
			panic(err)
		}
		DefaultJSONPath = s
	}
}

func LoadConfig() (c *Config) {
	c = &Config{}
	//def := DefaultConfig()
	f := filepath.Join(DefaultJSONPath, JSONName)
	bys, e := os.ReadFile(f)
	if e != nil {
		panic(f)
	}
	e = json.Unmarshal(bys, &c)
	if e != nil {
		panic(bys)
	}
	return c
}

func OutputConfig(config *Config) error {
	bys, e := json.MarshalIndent(config, "", " ")
	if e != nil {
		panic(e)
	}

	return os.WriteFile(filepath.Join(DefaultJSONPath, JSONName), bys, 0755)
}

func DefaultConfig() *Config {
	return &Config{
		TmpDir:      "tmp",
		BaseUrl:     "http://tool.httpcn.com/Zi/So.asp",
		UnicodeFile: "GB.txt",
		Action:      ActionNon,
		DatabaseExc: config.DatabaseConfig{
			Host:         "localhost",
			Port:         "3306",
			User:         "root",
			Pwd:          "111111",
			Name:         "excavator",
			MaxIdleCon:   0,
			MaxOpenCon:   0,
			Driver:       "sqlite3",
			File:         "exc.db",
			Dsn:          "",
			ShowSQL:      false,
			ShowExecTime: false,
		},
		DatabaseFate: config.DatabaseConfig{
			Host:         "localhost",
			Port:         "3306",
			User:         "root",
			Pwd:          "111111",
			Name:         "fate",
			MaxIdleCon:   0,
			MaxOpenCon:   0,
			Driver:       "sqlite3",
			File:         "ft.db",
			Dsn:          "",
			ShowSQL:      false,
			ShowExecTime: false,
		},
	}
}
