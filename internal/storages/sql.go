package storages

import "github.com/pkg/errors"

var (
	sqlInsertColors = `INSERT INTO colors (name, pin_code, hex) VALUES (?, ?, ?);`
	sqlSelectColors = `SELECT * FROM colors;`
	sqlDeleteColor  = `DELETE FROM colors WHERE id = ?;`
)

func (s *MySQLStorage) initStatement() error {
	var err error
	if s.stmtInsertColor, err = s.DB().Prepare(sqlInsertColors); err != nil {
		return errors.WithStack(err)
	}
	if s.stmtSelectColor, err = s.DB().Prepare(sqlSelectColors); err != nil {
		return errors.WithStack(err)
	}
	if s.stmtDeleteColor, err = s.DB().Prepare(sqlDeleteColor); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
