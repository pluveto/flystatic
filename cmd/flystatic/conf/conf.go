package conf

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pluveto/flystatic/pkg/logger"
)

type HashMethond string

const BcryptHash HashMethond = "bcrypt"
const SHA256Hash HashMethond = "sha256"

func GetDefaultConf() Conf {
	defaultFsDir, _ := os.Getwd()
	if !strings.HasPrefix(defaultFsDir, "/home") {
		serveDir := filepath.Join(os.TempDir(), "flystatic")
		err := os.MkdirAll(serveDir, 0755)
		if err != nil {
			logger.Fatal("Failed to create tmp dir", err)
		}
		defaultFsDir = serveDir
	}
	return Conf{
		Log: Log{
			Level: "Warning",
			Stdout: []Stdout{
				{
					Format: "text",
					Output: "stdout",
				},
			},
			File: []File{},
		},
		Server: Server{
			Host:              "127.0.0.1",
			Port:              7086,
			Path:              "",
			FsDir:             defaultFsDir,
			DefaultSpeedLimit: 0,
		},
		Auth: Auth{
			User: []User{},
		},
		UI: UI{
			Enabled: false,
			Path:    "/ui",
			Source:  "",
		},
		CORS: CORS{
			Enabled: false,
		},
		SSL: SSL{
			Enabled: false,
		},
	}
}

type Conf struct {
	Log    Log    `toml:"log"`
	Server Server `toml:"server"`
	Auth   Auth   `toml:"auth"`
	UI     UI     `toml:"ui"`
	CORS   CORS   `toml:"cors"`
	SSL    SSL    `toml:"ssl"`
}

type CORS struct {
	Enabled          bool     `toml:"enabled"`
	AllowedOrigins   []string `toml:"allowed_origins"`
	AllowedMethods   []string `toml:"allowed_methods"`
	AllowedHeaders   []string `toml:"allowed_headers"`
	ExposedHeaders   []string `toml:"exposed_headers"`
	AllowCredentials bool     `toml:"allow_credentials"`
	MaxAge           int
}

type Server struct {
	Host              string `toml:"host"`
	Port              int    `toml:"port"`
	Path              string `toml:"path"`
	FsDir             string `toml:"fs_dir"`
	DefaultSpeedLimit uint64 `toml:"default_speed_limit"` // octet/s
}

type UI struct {
	Enabled bool   `toml:"enabled"`
	Path    string `toml:"path"`   // Path prefix. TODO: ui.path cannot equals to server.path
	Source  string `toml:"source"` // Source location of the UI
}

type User struct {
	SubPath       string      `toml:"sub_path"`
	SubFsDir      string      `toml:"sub_fs_dir"`
	Username      string      `toml:"username"`
	SpeedLimit    uint64      `toml:"speed_limit"` // octet/s
	ShowDirectory bool        `toml:"show_directory"` // show directory index
	PasswordHash  string      `toml:"password_hash"`
	PasswordCrypt HashMethond `toml:"password_crypt"`
}
type Auth struct {
	User []User `toml:"user"`
}

type File struct {
	Format  LogFormat `toml:"format"`
	Path    string    `toml:"path"`
	MaxSize int       `toml:"max_size"`
	MaxAge  int       `toml:"max_age"`
}

type LogFormat string

const (
	LogFormatJSON LogFormat = "json"
	LogFormatText LogFormat = "text"
)

type LogOutput string

const (
	LogOutputStdout LogOutput = "stdout"
	LogOutputStderr LogOutput = "stderr"
	LogOutputFile   LogOutput = "file"
)

type Stdout struct {
	Format LogFormat `toml:"format"`
	Output LogOutput `toml:"output"`
}
type Log struct {
	Level  string   `toml:"level"`
	File   []File   `toml:"file"`
	Stdout []Stdout `toml:"stdout"`
}

type SSL struct {
	Enabled bool   `toml:"enabled"`
	Cert    string `toml:"cert"`
	Key     string `toml:"key"`
}
