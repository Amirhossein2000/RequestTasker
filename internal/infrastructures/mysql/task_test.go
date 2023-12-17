package mysql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"RequestTasker/internal/domian/entities"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func setupMySQLContainer(t *testing.T) (*dbr.Connection, func()) {
	ctx := context.Background()

	databaseName := "test_db"

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
	require.NoError(t, err)

	host, err := mysqlC.Host(ctx)
	require.NoError(t, err)

	port, err := mysqlC.MappedPort(ctx, "3306")
	require.NoError(t, err)

	connectionString := fmt.Sprintf("root:root@tcp(%s:%s)/%s?parseTime=true", host, port.Port(), databaseName)

	time.Sleep(time.Second * 10)

	conn, err := dbr.Open("mysql", connectionString, nil)
	require.NoError(t, err)

	driver, err := mysql.WithInstance(conn.DB, &mysql.Config{
		MigrationsTable: mysql.DefaultMigrationsTable,
		DatabaseName:    databaseName,
	})
	require.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance(
		"file:../../../migrations",
		"mysql", driver,
	)
	require.NoError(t, err)

	err = m.Up()
	require.NoError(t, err)

	return conn, func() {
		mysqlC.Terminate(ctx)
		conn.Close()
	}
}

func TestMySQLTaskRepository(t *testing.T) {
	// Set up MySQL container
	conn, tearDown := setupMySQLContainer(t)
	defer tearDown()

	// Create a dbr.Session with the actual database connection
	session := conn.NewSession(nil)

	// Create the MySQLTaskRepository with the actual session
	repo := NewMySQLTaskRepository(session, "tasks")

	t.Run("Create and GetByPublicID", func(t *testing.T) {
		// Create a sample task
		task := entities.NewTask(
			"https://example.com",
			"GET",
			map[string]interface{}{"Authorization": "Bearer token"},
			map[string]interface{}{"key": "value"},
		)

		// Perform the Create operation
		createdTask, err := repo.Create(context.Background(), task)
		require.NoError(t, err)

		// Verify the result
		assert.Equal(t, task.Url(), createdTask.Url())
		assert.Equal(t, task.Method(), createdTask.Method())
		assert.Equal(t, task.Headers(), createdTask.Headers())
		assert.Equal(t, task.Body(), createdTask.Body())

		// Perform the GetByPublicID operation
		foundTask, err := repo.GetByPublicID(context.Background(), createdTask.PublicID())
		require.NoError(t, err)

		// Verify the result
		assert.Equal(t, createdTask.PublicID(), foundTask.PublicID())
		assert.Equal(t, createdTask.Url(), foundTask.Url())
		assert.Equal(t, createdTask.Method(), foundTask.Method())
		assert.Equal(t, createdTask.Headers(), foundTask.Headers())
		assert.Equal(t, createdTask.Body(), foundTask.Body())
	})
}
