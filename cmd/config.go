package cmd

import (
	"time"

	"github.com/onetwoclimb/internal/storages"

	"github.com/spf13/pflag"
)

type Config struct {
	ServerConfig
	DB storages.Config
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

func (c *ServerConfig) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("APIConfig", pflag.PanicOnError)

	f.StringVar(&c.Host, "host", "127.0.0.1", "ip")
	f.IntVar(&c.Port, "port", 8081, "port")
	f.StringVar(&c.Mode, "mode", ReleaseMode, "release,debug,test")
	f.AddFlagSet(config.DB.Flags("mysql"))

	f.DurationVar(&c.ReadTimeout, "readtimeout", time.Duration(0), "api read timeout (default 0s)")
	f.DurationVar(&c.WriteTimeout, "writetimeout", time.Duration(0), "api write timeout (default 0s)")

	return f
}
