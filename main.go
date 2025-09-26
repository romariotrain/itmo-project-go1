package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	const url = "http://srv.msk01.gigacorp.local/_stats"
	errs := 0

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		resp, err := http.Get(url)
		if err != nil {
			errs++
			checkErrors(errs)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			errs++
			checkErrors(errs)
			resp.Body.Close()
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			errs++
			checkErrors(errs)
			continue
		}

		data := strings.TrimSpace(string(body))
		parts := strings.Split(data, ",")
		if len(parts) != 7 {
			errs++
			checkErrors(errs)
			continue
		}

		values := make([]int64, 7)
		for i, part := range parts {
			val, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				errs++
				checkErrors(errs)
				continue
			}
			values[i] = val
		}

		loadAvg := values[0]
		memTotal := values[1]
		memUsed := values[2]
		diskTotal := values[3]
		diskUsed := values[4]
		netTotal := values[5]
		netUsed := values[6]

		// Load Average
		if loadAvg > 30 {
			fmt.Printf("Load Average is too high: %d\n", loadAvg)
		}

		// Memory
		if memTotal > 0 {
			usage := memUsed * 100 / memTotal
			if usage > 80 {
				fmt.Printf("Memory usage too high: %d%%\n", usage)
			}
		}

		// Disk
		if diskTotal > 0 {
			usage := diskUsed * 100 / diskTotal
			if usage > 90 {
				freeDiskMb := (diskTotal - diskUsed) / 1024 / 1024
				fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskMb)
			}
		}

		// Network
		if netTotal > 0 {
			usage := netUsed * 100 / netTotal
			if usage > 90 {
				freeMbit := (netTotal - netUsed) / 1000 / 1000
				fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeMbit)
			}
		}
	}
}

func checkErrors(errs int) {
	if errs >= 3 {
		fmt.Println("Unable to fetch server statistic")
	}
}
