// setup_test.go
package mysql

import (
	"log"
	"os"
	"testing"

	"github.com/Amirhossein2000/RequestTasker/internal/pkg/integration"
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
