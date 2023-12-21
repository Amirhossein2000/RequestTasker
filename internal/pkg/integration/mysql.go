package integration

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"github.com/testcontainers/testcontainers-go"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type testDB struct {
	conn       *dbr.Connection
	inUseCount int64
}

var tb *testDB

const databaseName = "test_db"

// TODO: return conn
func SetupMySQLContainer() (*dbr.Session, func(), error) {
	if tb != nil {
		atomic.AddInt64(&tb.inUseCount, 1)
		return tb.conn.NewSession(nil),
			func() {
				atomic.AddInt64(&tb.inUseCount, -1)
			},
			nil
	}

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mysql:latest",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      databaseName,
		},
	}

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	host, err := mysqlC.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	port, err := mysqlC.MappedPort(ctx, "3306")
	if err != nil {
		return nil, nil, err
	}
	connectionString := fmt.Sprintf("root:root@tcp(%s:%s)/%s?parseTime=true", host, port.Port(), databaseName)

	time.Sleep(time.Second * 10)

	conn, err := dbr.Open("mysql", connectionString, nil)
	if err != nil {
		return nil, nil, err
	}

	driver, err := mysql.WithInstance(conn.DB, &mysql.Config{
		MigrationsTable: mysql.DefaultMigrationsTable,
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
		conn:       conn,
		inUseCount: 1,
	}

	return conn.NewSession(nil),
		func() {
			inUseCount := atomic.AddInt64(&tb.inUseCount, -1)
			if inUseCount < 1 {
				mysqlC.Terminate(ctx)
				conn.Close()
			}
		},
		nil
}
