package storages

import (
	"database/sql"

	"github.com/pkg/errors"
)

type Storage interface {
	DB() *sql.DB
	Close() error
}

type MySQLStorage struct {
	Storage

	stmtPutColors *sql.Stmt
	stmtGetColots *sql.Stmt
}

func NewMySQLStorage(storage Storage) (*MySQLStorage, error) {
	s := &MySQLStorage{Storage: storage}

	if err := s.initStatement(); err != nil {
		return nil, errors.WithStack(err)
	}

	return s, nil
}
