// setup_test.go
package mysql

import (
	"RequestTasker/internal/pkg/integration"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	session, tearDown, err := integration.SetupMySQLContainer()
	if err != nil {
		log.Panic(err)
	}
	defer tearDown()

	err = session.Ping()
	if err != nil {
		log.Panic(err)
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}
