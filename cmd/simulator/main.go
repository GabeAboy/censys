package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type CreateAssetRequest struct {
	IPAddress   string    `json:"ip_address"`
	Hostname    string    `json:"hostname"`
	PortNumbers []int     `json:"port_numbers"`
	Tags        []*string `json:"tags"`
}

type AssetCountResponse struct {
	Total int `json:"total"`
}

var (
	commonPorts = []int{22, 80, 443, 3306, 5432, 6379, 8080, 3389, 21, 25}
	commonTags  = []string{"production", "staging", "development", "web-server", "database", "cache", "api", "monitoring"}
	hostnames   = []string{"web-server", "db-server", "cache-server", "api-server", "worker", "load-balancer", "proxy"}
)

func main() {
	log.Println("Starting Asset Simulator...")

	// Get API URL from environment variable or use default
	apiBaseURL := os.Getenv("API_URL")
	if apiBaseURL == "" {
		apiBaseURL = "http://localhost:8080"
	}
	apiURL := fmt.Sprintf("%s/api/v1/assets", apiBaseURL)
	
	log.Printf("Using API URL: %s", apiURL)

	// Check asset count and create initial batch if needed
	count, err := getAssetCount(apiURL)
	if err != nil {
		log.Printf("Failed to get asset count: %v", err)
	} else {
		log.Printf("Current asset count: %d", count)
		if count < 5 {
			log.Printf("Asset count is less than 5, creating 10 initial assets...")
			for i := 0; i < 10; i++ {
				if err := createDummyAsset(apiURL); err != nil {
					log.Printf("Failed to create asset %d: %v", i+1, err)
				} else {
					log.Printf("Created initial asset %d/10", i+1)
				}
				time.Sleep(100 * time.Millisecond)
			}
			log.Println("Initial batch creation complete")
		}
	}

	// Schedule asset creation every minute
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("Simulator running. Creating new asset every minute...")

	for range ticker.C {
		if err := createDummyAsset(apiURL); err != nil {
			log.Printf("Failed to create asset: %v", err)
		}
	}
}

func createDummyAsset(apiURL string) error {
	// Generate random IP address
	ip := fmt.Sprintf("192.168.%d.%d", rand.Intn(256), rand.Intn(256))

	// Generate random hostname
	hostname := fmt.Sprintf("%s-%d", hostnames[rand.Intn(len(hostnames))], rand.Intn(100))

	// Generate random ports (1-5 ports)
	numPorts := rand.Intn(5) + 1
	ports := make([]int, numPorts)
	for i := 0; i < numPorts; i++ {
		ports[i] = commonPorts[rand.Intn(len(commonPorts))]
	}

	// Generate random tags (0-3 tags)
	numTags := rand.Intn(4)
	var tags []*string
	if numTags > 0 {
		tags = make([]*string, numTags)
		for i := 0; i < numTags; i++ {
			tag := commonTags[rand.Intn(len(commonTags))]
			tags[i] = &tag
		}
	}

	// Create request payload
	assetReq := CreateAssetRequest{
		IPAddress:   ip,
		Hostname:    hostname,
		PortNumbers: ports,
		Tags:        tags,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(assetReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP POST request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	log.Printf("✓ Created asset: %s (%s) with %d ports and %d tags", hostname, ip, len(ports), len(tags))
	return nil
}

func getAssetCount(apiURL string) (int, error) {
	countURL := apiURL + "/count"
	resp, err := http.Get(countURL)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response AssetCountResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Total, nil
}
