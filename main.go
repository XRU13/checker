package main

import (
	"fmt"
	"net/http"
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
		results <- fmt.Sprintf("Ошибка при проверке %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	results <- fmt.Sprintf("%s - Статус: %d", url, resp.StatusCode)
}

func main() {
	urls := []string{
		"https://example.com",
		"https://google.com",
		"https://github.com",
		"https://openai.com",
		"https://golang.org",
		"https://news.ycombinator.com",
		"https://stackoverflow.com",
		"https://fake-site123.com",
		"https://real-site.net",
		"https://mycoolblog.io",
		"https://thisisfake.co",
		"https://reddit.com",
		"https://linkedin.com",
		"https://microsoft.com",
		"https://amaz0n-shop.net", // fake-like
		"https://netflix.com",
		"https://go.dev",
		"https://superdata.ai",
		"https://techreview.fake",
		"https://openstream.tech",
		"https://yahoo.com",
		"https://tw1tter.net",     // fake-like
		"https://newssource.org",
		"https://weatherhub123.info",
		"https://somethingreal.edu",
		"https://cloudflare.com",
		"https://fakesite-login.org",
		"https://secure-account.help", // fake
		"https://devblog.dev",
		"https://untrustedsite.biz",
		"https://trustedsource.org",
		"https://funnycatvideos.tv",
		"https://g00gle.com", // fake-like
		"https://mybank-securelogin.com", // fake-like
		"https://python.org",
		"https://music-zone.cloud",
		"https://archive.org",
		"https://blog.fakesite.dev",
		"https://dailytips.info",
		"https://science-research.xyz",
		"https://universityportal.edu",
		"https://updatesoft.site",
		"https://sportsfan.live",
		"https://real-news.co.uk",
		"https://fakebank.support",
		"https://true-techworld.com",
		"https://ai-for-everyone.fake",
		"https://bestrecipes.io",
		"https://nasa.gov",
		"https://craigslist.org",
	}
	var wg sync.WaitGroup
	results := make(chan string, len(urls))

	fmt.Println("Начинаем проверку сайтов...")
	startTime := time.Now()

	for _, url := range urls {
		wg.Add(1)
		go checkSite(url, &wg, results)
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

	fmt.Printf("\nПроверка завершена за %v\n", time.Since(startTime))
}