package currency

import (
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	time.Local = time.UTC
	m.Run()
}
