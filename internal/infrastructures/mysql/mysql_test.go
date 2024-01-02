package mysql

import (
	"log"
	"os"
	"testing"

	"github.com/Amirhossein2000/RequestTasker/internal/pkg/integration"
)

func TestMain(m *testing.M) {
	conn, tearDown, err := integration.SetupMySQLContainer()
	if err != nil {
		log.Panic(err)
	}
	defer tearDown()

	err = conn.Ping()
	if err != nil {
		log.Panic(err)
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}
