package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/likexian/whois"
)

func runWhois(domain string) {
	result, err := whois.Whois(domain)
	if err != nil {
		fmt.Println("[X] WHOIS Hatası:", err)
		return
	}
	fmt.Println("[✓] WHOIS Sonucu:\n", result)
}

func runIPLookup(ip string) {
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		fmt.Println("[X] IP sorgu hatası:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	fmt.Printf("[✓] IP Bilgisi: %s (%s, %s)\n", data["query"], data["country"], data["org"])
}

func runShodan(ip string, apiKey string) {
	url := fmt.Sprintf("https://api.shodan.io/shodan/host/%s?key=%s", ip, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("[X] Shodan bağlantı hatası:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	fmt.Println("[✓] Shodan Cevabı:")
	for _, item := range result["data"].([]interface{}) {
		port := item.(map[string]interface{})["port"]
		product := item.(map[string]interface{})["product"]
		fmt.Printf("→ Port %v - %v\n", port, product)
	}
}

func main() {
	domain := flag.String("domain", "", "Domain için WHOIS sorgusu")
	ip := flag.String("ip", "", "IP için coğrafi bilgi")
	shodan := flag.String("shodan", "", "Shodan IP taraması")
	apiKey := flag.String("apikey", "", "Shodan API key")

	flag.Parse()

	if *domain != "" {
		runWhois(*domain)
	}

	if *ip != "" {
		runIPLookup(*ip)
	}

	if *shodan != "" && *apiKey != "" {
		runShodan(*shodan, *apiKey)
	}

	if *domain == "" && *ip == "" && *shodan == "" {
		fmt.Println("Kullanım örnekleri:")
		fmt.Println("  go run intelsentry.go -domain example.com")
		fmt.Println("  go run intelsentry.go -ip 1.2.3.4")
		fmt.Println("  go run intelsentry.go -shodan 1.2.3.4 -apikey YOUR_KEY")
	}
}
