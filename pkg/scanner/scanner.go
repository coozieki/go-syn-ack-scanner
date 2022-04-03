package scanner

import (
	"fmt"
	"net"
	"time"
)

const MAX_PORT_NUMBER uint = 65535
const DEFAULT_MAX_THREADS uint = 8

type Connector interface {
	Connect(ip string, port uint) bool
}

type NetPackageConnector struct {
	Timeout time.Duration
	Network string
}

func (connector NetPackageConnector) Connect(ip string, port uint) bool {
	timeout := time.Second
	network := "tcp"
	address := fmt.Sprintf("%s:%d", ip, port)
	if connector.Timeout != 0 {
		timeout = connector.Timeout
	}
	if connector.Network != "" { // TODO: Такая инверсия не позволяет проверить поведение данного условия в тестах
		network = connector.Network
	}
	connect, error := net.DialTimeout(network, address, timeout) // TODO: Тут нет инверссии и интерфейсов, поэтому, ты не смог это написать по TDD
	if error == nil {
		defer connect.Close() // TODO: Это основная ветвь программы, я бы инвертировал проверку
		return true
	}

	return false
}

type Logger interface {
	Log(message string)
}

type ConsoleLogger struct{} // TODO: Это адаптер

func (cl ConsoleLogger) Log(message string) {
	fmt.Println(message)
}

type scanner struct {
	connector          Connector
	logger             Logger
	checkedPortChannel chan uint
	openPortChannel    chan uint
	portToCheckChannel chan uint
	maxThreads         uint
}

type ScannerParams struct {
	Connector  Connector
	MaxThreads uint
	Logger     Logger
}

func NewScanner(params ScannerParams) scanner { // TODO: Exported function with the unexported return type
	var connector Connector // TODO: Жесткая связь, будет мешать инверсии зависимостей
	var logger Logger

	connector = NetPackageConnector{}
	logger = ConsoleLogger{}
	maxThreads := DEFAULT_MAX_THREADS

	if params.MaxThreads != 0 {
		maxThreads = params.MaxThreads
	}
	if params.Connector != nil { // TODO: Тоже плохая форма инверсии, не позволяет проверить поведение программы полностью
		connector = params.Connector
	}
	if params.Logger != nil {
		logger = params.Logger
	}

	portChannel := make(chan uint, 1)
	openPortChannel := make(chan uint, 1)
	portToCheckChannel := make(chan uint, maxThreads)

	return scanner{
		connector:          connector,
		logger:             logger,
		checkedPortChannel: portChannel,
		openPortChannel:    openPortChannel,
		maxThreads:         maxThreads,
		portToCheckChannel: portToCheckChannel,
	}
}

func (s *scanner) Scan(ip string) []uint {
	res := []uint{}
	for i := uint(1); i <= s.maxThreads; i++ {
		s.portToCheckChannel <- i // TODO: Нарушение DRY, ниже еще набиваем его портами, лучше было выделить в отдельный поток
	}
	for i := uint(0); i < s.maxThreads; i++ {
		go s.runCheckPortsWorker(ip)
	}
	var countPortsChecked uint = 0
	for {
		if countPortsChecked == MAX_PORT_NUMBER {
			s.closeChannels()
			break
		}
		select {
		case checkedPort := <-s.checkedPortChannel:
			s.logger.Log(fmt.Sprintf("Отсканирован порт: %d", checkedPort))
			countPortsChecked++
			if countPortsChecked+s.maxThreads > MAX_PORT_NUMBER { // TODO: Сложно понять эту проверку, но я справился
				continue
			}
			s.portToCheckChannel <- countPortsChecked + s.maxThreads // TODO: Тут тоже сложновато, через меняц будет сложно понять, почему мы их складываем
		case openedPort := <-s.openPortChannel:
			res = append(res, openedPort) // TODO: Только поедаем память, могли возвращать данные через канал, не занимая столько памяти
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
