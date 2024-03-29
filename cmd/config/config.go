package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/onetwoclimb/internal/storages"
)

type Config struct {
	ServerConfig
	DB          storages.Config
	StaticDst   string
	MaxFileSize int
}

const (
	DebugMode   string = "debug"
	ReleaseMode string = "release"
	TestMode    string = "test"
)

type ServerConfig struct {
	Host string
	Port int
	Mode string

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (c *Config) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("APIConfig", pflag.PanicOnError)

	f.StringVar(&c.Host, "host", "0.0.0.0", "ip")
	f.IntVar(&c.Port, "port", 80, "port")
	f.StringVar(&c.StaticDst, "static_dst", "/Users/fedoseevvt/go/src/github.com/onetwoclimb", "path to static files")
	f.IntVar(&c.MaxFileSize, "max_file_size", 5000000, "max image size to upload")
	f.StringVar(&c.Mode, "mode", ReleaseMode, "release,debug,test")
	f.AddFlagSet(c.DB.Flags("mysql"))

	f.DurationVar(&c.ReadTimeout, "readtimeout", time.Duration(0), "api read timeout (default 0s)")
	f.DurationVar(&c.WriteTimeout, "writetimeout", time.Duration(0), "api write timeout (default 0s)")

	return f
}

func (c *Config) Validate() error {
	if _, err := os.Stat(c.StaticDst); os.IsNotExist(err) {
		return errors.WithStack(err)
	}
	return nil
}
