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
	for {
		resp, err := http.Get("http://127.0.0.1:80/_stats")
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			time.Sleep(time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		fields := strings.Split(strings.TrimSpace(string(body)), ",")
		if len(fields) != 7 {
			time.Sleep(time.Second)
			continue
		}

		// парсим в int
		vals := make([]int, 7)
		ok := true
		for i, f := range fields {
			val, err := strconv.Atoi(f)
			if err != nil {
				ok = false
				break
			}
			vals[i] = val
		}
		if !ok {
			time.Sleep(time.Second)
			continue
		}

		// раскладываем
		loadAvg := vals[0]
		memAvail := vals[1]
		memUsed := vals[2]
		diskAvail := vals[3]
		diskUsed := vals[4]
		netAvail := vals[5]
		netUsed := vals[6]

		// проверки
		if loadAvg > 30 {
			fmt.Printf("Load Average is too high: %d\n", loadAvg)
		}

		memPercent := memUsed * 100 / memAvail
		if memPercent > 85 {
			fmt.Printf("Memory usage too high: %d%%\n", memPercent)
		}

		freeDiskMb := (diskAvail - diskUsed) / 1024 / 1024
		if freeDiskMb < 1024 {
			fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskMb)
		}

		freeNetMbit := (netAvail - netUsed) / 1000 / 1000 // ВАЖНО: именно 1000, а не 1024
		if freeNetMbit < 1000 {
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeNetMbit)
		}

		time.Sleep(time.Second)
	}
}
