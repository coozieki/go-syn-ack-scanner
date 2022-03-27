package main

import (
	"fmt"
	"go-syn-ack-scanner/pkg/scanner"
	"os"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	var maxThreads int
	if len(os.Args) >= 3 {
		maxThreads, _ = strconv.Atoi(os.Args[2])
	}
	scanner := scanner.NewScanner(
		scanner.ScannerParams{
			Connector:  scanner.NetPackageConnector{Timeout: time.Millisecond * 300, Network: "tcp"},
			MaxThreads: maxThreads,
		},
	)
	openedPorts := scanner.Scan(os.Args[1])
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("Открытые порты: ", openedPorts)
	fmt.Println("--------------------------------------------------------------")
	fmt.Printf("Продолжительность сканирования: %.2f сек", time.Since(start).Seconds())
	fmt.Println()
	fmt.Println("--------------------------------------------------------------")
}
