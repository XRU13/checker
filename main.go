package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func checkSite(url string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		results <- fmt.Sprintf("Error checking %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	results <- fmt.Sprintf("%s - Status: %d", url, resp.StatusCode)
}

func main() {
	// Открываем файл с URL-адресами
	file, err := os.Open("urls.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Читаем URL-адреса из файла
	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		if url != "" {
			urls = append(urls, url)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	const maxWorkers = 10
	var wg sync.WaitGroup
	results := make(chan string, len(urls))
	semaphore := make(chan struct{}, maxWorkers)

	fmt.Println("Start site review...")
	startTime := time.Now()

	for _, url := range urls {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(url string) {
			defer func() {
				<-semaphore
			}()
			checkSite(url, &wg, results)
		}(url)
	}

	// Запускаем горутину для закрытия канала после завершения всех проверок
	go func() {
		wg.Wait()
		close(results)
	}()

	// Выводим результаты по мере их поступления
	for result := range results {
		fmt.Println(result)
	}

	fmt.Printf("\nVerification completed for %v\n", time.Since(startTime))
}
