package cmd

import (
	"time"

	"github.com/spf13/pflag"
)

const (
	DebugMode   string = "debug"
	ReleaseMode string = "release"
	TestMode    string = "test"
)

type Config struct {
	Bind string
	Mode string

	RequestTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

func (c *Config) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("APIConfig", pflag.PanicOnError)

	f.StringVar(&c.Bind, "bind", "127.0.0.1:8082", "ip:port")
	f.StringVar(&c.Mode, "mode", ReleaseMode, "release,debug,test")

	f.DurationVar(&c.RequestTimeout, "requesttimeout", 60*time.Second, "api request timeout (default 60s)")
	f.DurationVar(&c.ReadTimeout, "readtimeout", time.Duration(0), "api read timeout (default 0s)")
	f.DurationVar(&c.WriteTimeout, "writetimeout", time.Duration(0), "api write timeout (default 0s)")

	return f
}
