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

	stmtInsertColor *sql.Stmt
	stmtSelectColor *sql.Stmt
	stmtDeleteColor *sql.Stmt
}

func NewMySQLStorage(storage Storage) (*MySQLStorage, error) {
	s := &MySQLStorage{Storage: storage}

	if err := s.initStatement(); err != nil {
		return nil, errors.WithStack(err)
	}

	return s, nil
}

type Color struct {
	Id      int32  `json:"id"`
	Name    string `json:"name"`
	PinCode string `json:"pinCode"`
	Hex     string `json:"hex"`
}

func (s *MySQLStorage) GetColors() ([]*Color, error) {
	rows, err := s.stmtSelectColor.Query()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := make([]*Color, 0)
	for rows.Next() {
		color := &Color{}
		if err := rows.Scan(&color.Id, &color.Name, &color.PinCode, &color.Hex); err != nil {
			return nil, errors.WithStack(err)
		}
		result = append(result, color)
	}
	return result, nil
}

func (s *MySQLStorage) DelColor(id int) error {
	if _, err := s.stmtDeleteColor.Exec(id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *MySQLStorage) PutColor(color *Color) error {
	if _, err := s.stmtInsertColor.Exec(&color.Name, &color.PinCode, &color.Hex); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
