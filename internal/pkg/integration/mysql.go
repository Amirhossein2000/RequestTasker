package integration

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"github.com/testcontainers/testcontainers-go"

	"github.com/golang-migrate/migrate/v4"
	migsql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type testDB struct {
	conn *dbr.Connection
}

var tb *testDB

const databaseName = "test_db"

// TODO: return conn
func SetupMySQLContainer() (*dbr.Connection, func(), error) {
	if tb != nil {
		return tb.conn, func() {}, nil
	}

	ctx := context.Background()

	mysqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:latest"),
		mysql.WithDatabase(databaseName),
		mysql.WithUsername("root"),
		mysql.WithPassword("root"),
	)
	if err != nil {
		return nil, nil, err
	}

	terminateFunc := func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	}

	host, err := mysqlContainer.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		return nil, nil, err
	}
	connectionString := fmt.Sprintf("root:root@tcp(%s:%s)/%s?parseTime=true", host, port.Port(), databaseName)

	time.Sleep(time.Second * 10)

	conn, err := dbr.Open("mysql", connectionString, nil)
	if err != nil {
		return nil, nil, err
	}

	driver, err := migsql.WithInstance(conn.DB, &migsql.Config{
		MigrationsTable: migsql.DefaultMigrationsTable,
		DatabaseName:    databaseName,
	})
	if err != nil {
		return nil, nil, err
	}

	_, currentFile, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(currentFile), "../../../migrations")

	m, err := migrate.NewWithDatabaseInstance(
		"file:"+migrationsPath,
		"mysql", driver,
	)
	if err != nil {
		return nil, nil, err
	}

	err = m.Up()
	if err != nil {
		return nil, nil, err
	}

	tb = &testDB{
		conn: conn,
	}

	return conn, terminateFunc, nil
}
