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

	ticker := time.NewTicker(10 * time.Second) // опрос каждые 10 секунд
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
			// Load Average может быть дробным
			if i == 0 {
				load, err := strconv.ParseFloat(part, 64)
				if err != nil {
					errs++
					checkErrors(errs)
					break
				}
				mapa["Load Average"] = int64(load)
				if load > 30 {
					fmt.Printf("Load Average is too high: %.2f\n", load)
				}
				continue
			}

			val, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				errs++
				checkErrors(errs)
				break
			}

			switch i {
			case 1:
				mapa["Mem Total"] = val
			case 2:
				mapa["Mem Used"] = val
				if mapa["Mem Total"] > 0 {
					usage := float64(val) / float64(mapa["Mem Total"]) * 100
					if usage > 80 {
						fmt.Printf("Memory usage too high: %.1f%%\n", usage)
					}
				}
			case 3:
				mapa["Disk Total"] = val
			case 4:
				mapa["Disk Used"] = val
				if mapa["Disk Total"] > 0 {
					usage := float64(val) / float64(mapa["Disk Total"])
					if usage > 0.9 {
						left := (mapa["Disk Total"] - val) / 1024 / 1024
						fmt.Printf("Free disk space is too low: %d Mb left\n", left)
					}
				}
			case 5:
				mapa["Net Total"] = val
			case 6:
				mapa["Net Used"] = val
				if mapa["Net Total"] > 0 {
					usage := float64(val) / float64(mapa["Net Total"])
					if usage > 0.9 {
						freeMbit := float64((mapa["Net Total"]-val)*8) / 1024 / 1024
						fmt.Printf("Network bandwidth usage high: %.2f Mbit/s available\n", freeMbit)
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
