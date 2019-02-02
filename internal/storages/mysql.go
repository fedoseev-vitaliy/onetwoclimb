package storages

import "database/sql"

func New() *sql.DB {
	//db, err := sql.Open("mysql", "user:password@/dbname")
	return &sql.DB{} // todo add db implementation
}
