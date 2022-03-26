package scanner

import (
	"fmt"
	"net"
	"time"
)

const MAX_PORT_NUMBER int = 65536
const DEFAULT_MAX_THREADS int = 8

type MakesConnection interface {
	Connect(ip string, port int) bool
}

type scanner struct {
	connector          MakesConnection
	checkedPortChannel chan int
	openPortChannel    chan int
	portToCheckChannel chan int
	maxThreads         int
}

type ScannerParams struct {
	Connector  MakesConnection
	MaxThreads int
}

type NetPackageConnector struct {
	Timeout time.Duration
	Network string
}

func (connector NetPackageConnector) Connect(ip string, port int) bool {
	timeout := time.Second
	network := "tcp"
	address := fmt.Sprintf("%s:%d", ip, port)
	if connector.Timeout != 0 {
		timeout = connector.Timeout
	}
	if connector.Network != "" {
		network = connector.Network
	}
	connect, error := net.DialTimeout(network, address, timeout)
	if error == nil {
		connect.Close()
		return true
	}

	return false
}

func NewScanner(params ScannerParams) scanner {
	var connector MakesConnection
	connector = NetPackageConnector{}
	maxThreads := DEFAULT_MAX_THREADS

	if params.MaxThreads != 0 {
		maxThreads = params.MaxThreads
	}
	if params.Connector != nil {
		connector = params.Connector
	}

	portChannel := make(chan int, 1)
	openPortChannel := make(chan int, 1)
	portToCheckChannel := make(chan int, maxThreads)

	return scanner{
		connector:          connector,
		checkedPortChannel: portChannel,
		openPortChannel:    openPortChannel,
		maxThreads:         maxThreads,
		portToCheckChannel: portToCheckChannel,
	}
}

func (s *scanner) Scan(ip string) []int {
	var res []int
	for i := 1; i <= s.maxThreads; i++ {
		s.portToCheckChannel <- i
	}
	for i := 0; i < s.maxThreads; i++ {
		go s.runCheckPortsWorker(ip)
	}
	countPortsChecked := 0
	var scannedPort, openedPort int
	for {
		if countPortsChecked == MAX_PORT_NUMBER {
			s.closeChannels()
			break
		}
		select {
		case scannedPort = <-s.checkedPortChannel:
			fmt.Println("Отсканирован порт: ", scannedPort)
			countPortsChecked++
			if countPortsChecked+s.maxThreads > MAX_PORT_NUMBER {
				continue
			}
			s.portToCheckChannel <- countPortsChecked + s.maxThreads
		case openedPort = <-s.openPortChannel:
			res = append(res, openedPort)
		}
	}
	return res
}

func (s *scanner) runCheckPortsWorker(ip string) {
	for port := range s.portToCheckChannel {
		if s.connector.Connect(ip, port) {
			s.openPortChannel <- port
		}
		s.checkedPortChannel <- port
	}
}

func (s *scanner) closeChannels() {
	close(s.checkedPortChannel)
	close(s.openPortChannel)
	close(s.portToCheckChannel)
}
