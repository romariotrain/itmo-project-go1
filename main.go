package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", "http://srv.msk01.gigacorp.local/_stats", nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Читаем всё тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Выводим статус и тело
	fmt.Println("Статус:", resp.Status)
	// fmt.Println("Тело ответа:")
	fmt.Println(string(body))
}
