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

		if resp.StatusCode != 200 {
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

		mapa := make(map[string]int64)

		for i, part := range parts {
			val, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				errs++
				checkErrors(errs)
				break
			}

			switch i {
			case 0: // Load Average
				mapa["Load Average"] = val
				if val > 30 {
					fmt.Printf("Load Average is too high: %d\n", val)
				}

			case 1: // Mem Total
				mapa["Mem Total"] = val

			case 2: // Mem Used
				mapa["Mem Used"] = val
				if mapa["Mem Total"] > 0 {
					usage := val * 100 / mapa["Mem Total"]
					if usage > 80 {
						fmt.Printf("Memory usage too high: %d%%\n", usage)
					}
				}

			case 3: // Disk Total
				mapa["Disk Total"] = val

			case 4: // Disk Used
				mapa["Disk Used"] = val
				if mapa["Disk Total"] > 0 {
					usage := val * 100 / mapa["Disk Total"]
					if usage > 90 {
						left := (mapa["Disk Total"] - val) / 1024 / 1024
						fmt.Printf("Free disk space is too low: %d Mb left\n", left)
					}
				}

			case 5: // Net Total
				mapa["Net Total"] = val

			case 6: // Net Used
				mapa["Net Used"] = val
				if mapa["Net Total"] > 0 {
					usage := val * 100 / mapa["Net Total"]
					if usage > 90 {
						freeMbit := float64(mapa["Net Total"]-val) / 1024 / 1024
						fmt.Printf("Network bandwidth usage high: %.0f Mbit/s available\n", freeMbit)
					}
				}
			}
		}
	}
}

func checkErrors(errs int) {
	if errs >= 3 {
		fmt.Println("Unable to fetch server statistic")
	}
}
