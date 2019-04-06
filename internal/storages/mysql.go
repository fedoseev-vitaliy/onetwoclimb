package storages

import (
	"database/sql"
	"fmt"
	"syscall"

	"github.com/pkg/errors"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type MySQL struct {
	db *sql.DB
}

func (s *MySQL) DB() *sql.DB {
	return s.db
}

func (s *MySQL) Close() error {
	db := s.db
	s.db = nil
	return db.Close()
}

func New(conf *Config) (*MySQL, error) {
	var err error

	c := mysql.NewConfig()
	c.DBName = conf.Database
	c.User = conf.User
	c.Passwd = conf.Password
	c.Net = "tcp"
	c.Addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	c.ParseTime = true

	s := &MySQL{}
	dsn := c.FormatDSN()
	if s.db, err = sql.Open("mysql", dsn); err != nil {
		return nil, fmt.Errorf("failed to open MySQL DSN: %v", err)
	}

	if err := s.db.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}

	if conf.MaxOpenConns != 0 {
		s.db.SetMaxOpenConns(conf.MaxOpenConns)
	}
	if conf.MaxIdleConns != 0 {
		s.db.SetMaxIdleConns(conf.MaxIdleConns)
	}

	log.WithField("pid", syscall.Getpid()).
		WithField("dsn", fmt.Sprintf("%s:%s@%s:%d/%s?sslmode=disable", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)).
		Info("new storage")

	return s, nil
}
