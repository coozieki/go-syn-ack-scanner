package scanner_test

import (
	"go-syn-ack-scanner/pkg/scanner"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var checkedPorts []uint
var countLoggerCalls uint
var mu sync.Mutex

type connectorMock struct {
	openPorts []uint
}

func (c connectorMock) Connect(ip string, port uint) bool {
	mu.Lock()
	checkedPorts = append(checkedPorts, port)
	mu.Unlock()
	for _, openPort := range c.openPorts {
		if openPort == port {
			return true
		}
	}
	return false
}

type loggerMock struct{}

func (l loggerMock) Log(message string) {
	mu.Lock()
	countLoggerCalls++
	mu.Unlock()
}

func TestNewScanner(t *testing.T) {
	t.Run("all params", func(t *testing.T) {
		var maxThreads uint = 7
		scanner.NewScanner(scanner.ScannerParams{Connector: connectorMock{}, MaxThreads: maxThreads})
	})

	t.Run("no params", func(t *testing.T) {
		scanner.NewScanner(scanner.ScannerParams{})
	})
}

func TestScan(t *testing.T) {
	t.Run("with open ports", func(t *testing.T) {
		checkedPorts = []uint{}
		countLoggerCalls = 0

		openPorts := []uint{1, scanner.MAX_PORT_NUMBER}

		s := scanner.NewScanner(
			scanner.ScannerParams{
				Connector: connectorMock{openPorts: openPorts},
				Logger:    loggerMock{},
			},
		)

		assert.Equal(t, openPorts, s.Scan("ip"))
		assert.Equal(t, scanner.MAX_PORT_NUMBER, countLoggerCalls)
		assert.Equal(t, scanner.MAX_PORT_NUMBER, uint(len(checkedPorts)))
	})

	t.Run("without open port", func(t *testing.T) {
		checkedPorts = []uint{}
		countLoggerCalls = 0

		openPorts := []uint{}

		s := scanner.NewScanner(
			scanner.ScannerParams{
				Connector: connectorMock{openPorts: openPorts},
				Logger:    loggerMock{},
			},
		)

		assert.Equal(t, openPorts, s.Scan("ip"))
		assert.Equal(t, scanner.MAX_PORT_NUMBER, countLoggerCalls)
		assert.Equal(t, scanner.MAX_PORT_NUMBER, uint(len(checkedPorts)))
	})
}
