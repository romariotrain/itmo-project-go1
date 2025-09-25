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

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		resp, err := client.Get(url)
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

		// Парсим все значения сначала
		values := make([]int64, 7)
		parseError := false
		for i, part := range parts {
			val, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
			if err != nil || val < 0 {
				parseError = true
				break
			}
			values[i] = val
		}

		if parseError {
			errs++
			checkErrors(errs)
			continue
		}

		// Сброс счетчика ошибок при успешном парсинге
		errs = 0

		// Проверка Load Average
		if values[0] > 30 {
			fmt.Printf("Load Average is too high: %d\n", values[0])
		}

		// Проверка памяти
		if values[1] > 0 {
			usage := (values[2]) * 100 / (values[1])
			if usage > 80 {
				fmt.Printf("Memory usage too high: %.0d%%\n", usage) // Целое число
			}
		}

		// Проверка диска
		if values[3] > 0 {
			usage := float64(values[4]) * 100 / float64(values[3])
			if usage > 90 {
				left := (values[3] - values[4]) / (1024 * 1024) // MB
				fmt.Printf("Free disk space is too low: %d Mb left\n", left)
			}
		}

		// Проверка сети - ИСПРАВЛЕНО!
		if values[5] > 0 {
			usage := float64(values[6]) * 100 / float64(values[5])
			if usage > 90 {
				// Правильный расчет: байты/с → биты/с → мегабиты/с
				freeMbit := float64(values[5]-values[6]) * 8 / (1024 * 1024)
				if freeMbit == float64(int64(freeMbit)) {
					fmt.Printf("Network bandwidth usage high: %.0f Mbit/s available\n", freeMbit)
				} else {
					fmt.Printf("Network bandwidth usage high: %.1f Mbit/s available\n", freeMbit)
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
