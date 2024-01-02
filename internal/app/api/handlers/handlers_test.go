package handlers

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_, terminate, err := setUpTestEnv()
	if err != nil {
		log.Panic(err)
	}
	defer terminate()

	exitCode := m.Run()

	os.Exit(exitCode)
}
