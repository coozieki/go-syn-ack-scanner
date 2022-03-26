package scanner_test

import (
	"go-syn-ack-scanner/pkg/scanner"
	"testing"

	"github.com/stretchr/testify/assert"
)

type connectorMock struct {
	openedPort1 int
	openedPort2 int
}

func (c connectorMock) Connect(ip string, port int) bool {
	return c.openedPort1 == port || c.openedPort2 == port
}

func TestNewScannerWithRightParams(t *testing.T) {
	t.Run("all params", func(t *testing.T) {
		maxThreads := 7
		scanner.NewScanner(scanner.ScannerParams{Connector: connectorMock{}, MaxThreads: maxThreads})
	})

	t.Run("no params", func(t *testing.T) {
		scanner.NewScanner(scanner.ScannerParams{})
	})
}

func TestScan(t *testing.T) {
	t.Run("with opened ports", func(t *testing.T) {
		openedPort1 := 1
		openedPort2 := scanner.MAX_PORT_NUMBER
		s := scanner.NewScanner(
			scanner.ScannerParams{
				Connector: connectorMock{openedPort1: openedPort1, openedPort2: openedPort2},
			},
		)

		assert.Equal(t, []int{openedPort1, openedPort2}, s.Scan("ip"))
	})

	t.Run("without opened port", func(t *testing.T) {
		s := scanner.NewScanner(
			scanner.ScannerParams{
				Connector: connectorMock{},
			},
		)

		assert.Equal(t, []int(nil), s.Scan("ip"))
	})
}
