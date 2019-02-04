package storages

import "github.com/spf13/pflag"

type Config struct {
	Host         string
	Port         uint16
	User         string
	Password     string
	Database     string
	MaxOpenConns int
	MaxIdleConns int
}

func (c *Config) Flags(name string) *pflag.FlagSet {
	f := pflag.NewFlagSet(name, pflag.PanicOnError)
	f.StringVar(&c.Host, "db_host", "127.0.0.1", "")
	f.Uint16Var(&c.Port, "db_port", 3306, "")
	f.StringVar(&c.User, "db_user", "root", "")
	f.StringVar(&c.Password, "db_password", "", "[secret]")
	f.StringVar(&c.Database, "db_database", "onetwoclimb", "")
	f.IntVar(&c.MaxOpenConns, "db_maxopenconns", 16, "max open connections")
	f.IntVar(&c.MaxIdleConns, "db_maxidleconns", 16, "max idle connections")

	return f
}
