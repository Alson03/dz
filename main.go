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
	if loadAvg, err := strconv.ParseFloat(stats[0], 64); err == nil {
		if loadAvg > 30 {
			fmt.Printf("Load Average is too high: %.0f\n", loadAvg)
		}
	}

	memTotal, err1 := strconv.ParseUint(stats[1], 10, 64)
	memUsed, err2 := strconv.ParseUint(stats[2], 10, 64)
	if err1 == nil && err2 == nil && memTotal > 0 {
		memUsagePercent := (memUsed * 100) / memTotal
		if memUsagePercent > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", memUsagePercent)
		}
	}

	diskTotal, err1 := strconv.ParseUint(stats[3], 10, 64)
	diskUsed, err2 := strconv.ParseUint(stats[4], 10, 64)
	if err1 == nil && err2 == nil && diskTotal > 0 {
		diskUsagePercent := (diskUsed * 100) / diskTotal
		if diskUsagePercent > 90 {
			freeSpaceMB := (diskTotal - diskUsed) / (1024 * 1024)
			fmt.Printf("Free disk space is too low: %d Mb left\n", freeSpaceMB)
		}
	}

	netTotal, err1 := strconv.ParseUint(stats[5], 10, 64)
	netUsed, err2 := strconv.ParseUint(stats[6], 10, 64)
	if err1 == nil && err2 == nil && netTotal > 0 {
		netUsagePercent := (netUsed * 100) / netTotal
		if netUsagePercent > 90 {
			availableBandwidthBytes := netTotal - netUsed
			availableBandwidthMbit := availableBandwidthBytes / (1024 * 1024 * 8)
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", availableBandwidthMbit)
		}
	}
}