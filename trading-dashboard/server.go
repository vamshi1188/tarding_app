package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	clientID    = "J0SM53J6A6-100" // Replace with your client_id
	appSecret   = "XVENU2DVWH"
	redirectURI = "https://trading-app-ihw6.onrender.com/"
	authCodeURL = "https://api-t1.fyers.in/api/v3/generate-authcode"
	validateURL = "https://api-t1.fyers.in/api/v3/validate-authcode"
	refreshURL  = "https://api-t1.fyers.in/api/v3/validate-refresh-token"
	holdingsURL = "https://api.fyers.in/api/v3/holdings"
	ordersURL   = "https://api-t1.fyers.in/api/v3/orders"
)

var accessToken string

// Define structs for authentication and responses
type ValidateAuthCodeResponse struct {
	Status       string `json:"s"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type HoldingsResponse struct {
	Status string `json:"s"`
	Data   []struct {
		HoldingType string  `json:"holdingType"`
		Quantity    int     `json:"quantity"`
		Symbol      string  `json:"symbol"`
		LTP         float64 `json:"ltp"`
	} `json:"data"`
}

type OrderResponse struct {
	Status string `json:"s"`
	Data   []struct {
		OrderID     string  `json:"orderId"`
		Symbol      string  `json:"symbol"`
		OrderType   string  `json:"orderType"`
		Quantity    int     `json:"quantity"`
		Price       float64 `json:"price"`
		OrderStatus string  `json:"orderStatus"`
	} `json:"data"`
}

func getAuthCodeURL() string {
	state := "random_state" // Generate a random state value
	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&state=%s", authCodeURL, clientID, redirectURI, state)
}

func validateAuthCode(authCode string) (string, string, error) {
	appIDHash := sha256.Sum256([]byte(clientID + appSecret))
	appIDHashStr := hex.EncodeToString(appIDHash[:])

	payload := map[string]string{
		"grant_type": "authorization_code",
		"appIdHash":  appIDHashStr,
		"code":       authCode,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(validateURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %v", err)
	}

	var response ValidateAuthCodeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Status != "ok" {
		return "", "", fmt.Errorf("validation failed: %s", response.Message)
	}

	return response.AccessToken, response.RefreshToken, nil
}

func refreshAccessToken(refreshToken string) (string, error) {
	appIDHash := sha256.Sum256([]byte(clientID + appSecret))
	appIDHashStr := hex.EncodeToString(appIDHash[:])

	payload := map[string]string{
		"grant_type":    "refresh_token",
		"appIdHash":     appIDHashStr,
		"refresh_token": refreshToken,
		"pin":           "your_pin", // Add user-specific pin if required
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(refreshURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var response ValidateAuthCodeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Status != "ok" {
		return "", fmt.Errorf("refresh failed: %s", response.Message)
	}

	return response.AccessToken, nil
}

func fetchHoldings() (*HoldingsResponse, error) {
	req, err := http.NewRequest("GET", holdingsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var holdings HoldingsResponse
	err = json.Unmarshal(body, &holdings)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal holdings: %v", err)
	}

	if holdings.Status != "ok" {
		return nil, fmt.Errorf("failed to fetch holdings: %s", holdings.Status)
	}

	return &holdings, nil
}

func fetchOrders() (*OrderResponse, error) {
	req, err := http.NewRequest("GET", ordersURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var orders OrderResponse
	err = json.Unmarshal(body, &orders)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal orders: %v", err)
	}

	if orders.Status != "ok" {
		return nil, fmt.Errorf("failed to fetch orders: %s", orders.Status)
	}

	return &orders, nil
}

func main() {
	fmt.Println("Navigate to this URL to get auth code:")
	fmt.Println(getAuthCodeURL())

	// Prompt user for the auth_code (can be automated in production)
	var authCode string
	fmt.Print("Enter the auth code: ")
	fmt.Scan(&authCode)

	var err error
	accessToken, _, err = validateAuthCode(authCode)
	if err != nil {
		log.Fatalf("Failed to validate auth code: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		holdings, err := fetchHoldings()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch holdings: %v", err), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		orders, err := fetchOrders()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch orders: %v", err), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		// Log for debugging
		log.Printf("Holdings: %+v\n", holdings)
		log.Printf("Orders: %+v\n", orders)

		// Render holdings and orders
		fmt.Fprintf(w, "Holdings and Orders fetched successfully\n")
		fmt.Fprintf(w, "Holdings: %+v\n", holdings)
		fmt.Fprintf(w, "Orders: %+v\n", orders)
	})

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
