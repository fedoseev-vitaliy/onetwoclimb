package storages

import "github.com/pkg/errors"

var (
	sqlPutColors = `INSERT INTO colors (name, pin_code, hex) VALUES (?, ?, ?)`
	sqlGetColors = `SELECT * FROM colors`
)

func (s *MySQLStorage) initStatement() error {
	var err error
	if s.stmtPutColors, err = s.DB().Prepare(sqlPutColors); err != nil {
		return errors.WithStack(err)
	}
	if s.stmtGetColots, err = s.DB().Prepare(sqlGetColors); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
