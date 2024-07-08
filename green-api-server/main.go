package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Constants for API URL
const (
	BaseAPIURL = "https://api.green-api.com"
)

type APIResponse struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error,omitempty"`
}

// getAPIUrl returns the full URL for the API endpoint
func getAPIUrl(idInstance, apiTokenInstance, endpoint string) string {
	return fmt.Sprintf("%s/waInstance%s/%s/%s", BaseAPIURL, idInstance, endpoint, apiTokenInstance)
}

// fetchAPI performs an HTTP request to the API endpoint
func fetchAPI(url string, method string, body interface{}) (interface{}, error) {
	var req *http.Request
	var err error

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error! Status: %s", resp.Status)
	}

	var result interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// makeHandlerFunc returns an HTTP handler function for API endpoints
func makeHandlerFunc(endpointFunc func(w http.ResponseWriter, r *http.Request, idInstance, apiTokenInstance string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idInstance := r.URL.Query().Get("idInstance")
		apiTokenInstance := r.URL.Query().Get("apiTokenInstance")
		endpointFunc(w, r, idInstance, apiTokenInstance)
	}
}

func getSettingsHandler(w http.ResponseWriter, r *http.Request, idInstance, apiTokenInstance string) {
	url := getAPIUrl(idInstance, apiTokenInstance, "getSettings")
	handleAPICall(w, r, url, "GET", nil)
}

func getStateInstanceHandler(w http.ResponseWriter, r *http.Request, idInstance, apiTokenInstance string) {
	url := getAPIUrl(idInstance, apiTokenInstance, "getStateInstance")
	handleAPICall(w, r, url, "GET", nil)
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request, idInstance, apiTokenInstance string) {
	url := getAPIUrl(idInstance, apiTokenInstance, "sendMessage")
	handleAPICall(w, r, url, "POST", nil)
}

func sendFileByUrlHandler(w http.ResponseWriter, r *http.Request, idInstance, apiTokenInstance string) {
	url := getAPIUrl(idInstance, apiTokenInstance, "sendFileByUrl")
	handleAPICall(w, r, url, "POST", nil)
}

func handleAPICall(w http.ResponseWriter, r *http.Request, url, method string, body interface{}) {
	var requestBody interface{}
	if r.Method == "POST" || r.Method == "PUT" {
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	result, err := fetchAPI(url, method, requestBody)
	apiResponse := APIResponse{Result: result}
	if err != nil {
		apiResponse.Error = err.Error()
	}
	json.NewEncoder(w).Encode(apiResponse)
}

func main() {
	http.HandleFunc("/getSettings", makeHandlerFunc(getSettingsHandler))
	http.HandleFunc("/getStateInstance", makeHandlerFunc(getStateInstanceHandler))
	http.HandleFunc("/sendMessage", makeHandlerFunc(sendMessageHandler))
	http.HandleFunc("/sendFileByUrl", makeHandlerFunc(sendFileByUrlHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
