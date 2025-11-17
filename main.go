package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	serverURL := "http://srv.msk01.gigacorp.local/_stats"
	errorCount := 0

	for {
		resp, err := http.Get(serverURL)
		if err != nil {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				errorCount = 0
			}
			time.Sleep(10 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			errorCount++
			resp.Body.Close()
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				errorCount = 0
			}
			time.Sleep(10 * time.Second)
			continue
		}

		scanner := bufio.NewScanner(resp.Body)
		if scanner.Scan() {
			data := scanner.Text()
			stats := strings.Split(data, ",")

			if len(stats) >= 7 {
				processStats(stats)
				errorCount = 0
			} else {
				errorCount++
			}
		} else {
			errorCount++
		}

		resp.Body.Close()

		if errorCount >= 3 {
			fmt.Println("Unable to fetch server statistic")
			errorCount = 0
		}

		time.Sleep(10 * time.Second)
	}
}

func processStats(stats []string) {
	// Load Average - выводим без десятичных
	if loadAvg, err := strconv.ParseFloat(stats[0], 64); err == nil {
		if loadAvg > 30 {
			fmt.Printf("Load Average is too high: %.0f\n", loadAvg)
		}
	}

	// Memory usage - выводим проценты без десятичных
	memTotal, err1 := strconv.ParseUint(stats[1], 10, 64)
	memUsed, err2 := strconv.ParseUint(stats[2], 10, 64)
	if err1 == nil && err2 == nil && memTotal > 0 {
		memUsagePercent := float64(memUsed) / float64(memTotal) * 100
		if memUsagePercent > 80 {
			fmt.Printf("Memory usage too high: %.0f%%\n", memUsagePercent)
		}
	}

	// Disk space - выводим мегабайты без десятичных
	diskTotal, err1 := strconv.ParseUint(stats[3], 10, 64)
	diskUsed, err2 := strconv.ParseUint(stats[4], 10, 64)
	if err1 == nil && err2 == nil && diskTotal > 0 {
		diskUsagePercent := float64(diskUsed) / float64(diskTotal) * 100
		if diskUsagePercent > 90 {
			freeSpaceMB := float64(diskTotal-diskUsed) / 1024 / 1024
			fmt.Printf("Free disk space is too low: %.0f Mb left\n", freeSpaceMB)
		}
	}

	// Network bandwidth - выводим мегабиты без десятичных
	netTotal, err1 := strconv.ParseUint(stats[5], 10, 64)
	netUsed, err2 := strconv.ParseUint(stats[6], 10, 64)
	if err1 == nil && err2 == nil && netTotal > 0 {
		netUsagePercent := float64(netUsed) / float64(netTotal) * 100
		if netUsagePercent > 90 {
			availableBandwidthMbit := float64(netTotal-netUsed) * 8 / 1000000
			fmt.Printf("Network bandwidth usage high: %.0f Mbit/s available\n", availableBandwidthMbit)
		}
	}
}