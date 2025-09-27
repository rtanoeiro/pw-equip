package equip

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

// SubscriptionResponse represents the API response structure
type SubscriptionResponse struct {
	Status bool `json:"status"`
}

// GetHWID generates a unique hardware ID for the current machine
func GetHWID() (string, error) {
	// Get system information
	hostInfo, err := host.Info()
	if err != nil {
		return "", fmt.Errorf("failed to get host info: %v", err)
	}

	// Create a unique identifier based on system information
	identifier := fmt.Sprintf("%s-%s-%s-%s",
		hostInfo.HostID,
		hostInfo.Platform,
		hostInfo.PlatformFamily,
		runtime.GOARCH,
	)

	// Generate MD5 hash of the identifier
	hash := md5.Sum([]byte(identifier))
	hwid := fmt.Sprintf("%x", hash)

	return hwid, nil
}

// CheckSubscription verifies if the current machine has an active subscription
func CheckSubscription() (bool, error) {
	// Get HWID
	hwid, err := GetHWID()
	if err != nil {
		return false, fmt.Errorf("failed to generate HWID: %v", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make GET request to subscription API
	url := fmt.Sprintf("http://200.1.1.1/subscription?hwid=%s", hwid)
	resp, err := client.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to check subscription: %v", err)
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("subscription API returned status: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse JSON response
	var subscriptionResp SubscriptionResponse
	if err := json.Unmarshal(body, &subscriptionResp); err != nil {
		return false, fmt.Errorf("failed to parse response: %v", err)
	}

	return subscriptionResp.Status, nil
}

// ValidateSubscriptionWithRetry checks subscription with retry logic
func ValidateSubscriptionWithRetry(maxRetries int) (bool, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		isActive, err := CheckSubscription()
		if err == nil {
			return isActive, nil
		}

		lastErr = err
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second) // Exponential backoff
		}
	}

	return false, fmt.Errorf("subscription check failed after %d retries: %v", maxRetries, lastErr)
}

// DisplayHWID shows the current machine's HWID for debugging/registration purposes
func DisplayHWID() {
	hwid, err := GetHWID()
	if err != nil {
		fmt.Printf("Erro ao obter HWID: %v\n", err)
		return
	}
	fmt.Printf("HWID da máquina: %s\n", hwid)
}
