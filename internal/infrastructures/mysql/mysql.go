package mysql

import (
	"fmt"

	"github.com/gocraft/dbr/v2"
)

type MYSQLConfig struct {
	Addr     string
	User     string
	Password string
	Database string
}

func NewMYSQLConn(conf MYSQLConfig) (*dbr.Connection, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", conf.User, conf.Password, conf.Addr, conf.Database)

	conn, err := dbr.Open("mysql", connectionString, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
