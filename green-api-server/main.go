package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type APIResponse struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error,omitempty"`
}

func getAPIUrl(idInstance, apiTokenInstance, endpoint string) string {
	return fmt.Sprintf("https://api.green-api.com/waInstance%s/%s/%s", idInstance, endpoint, apiTokenInstance)
}

func fetchAPI(url string, method string, body interface{}) (interface{}, error) {
	var req *http.Request
	var err error

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, ioutil.NopCloser(bytes.NewReader(jsonData)))
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

func getSettings(w http.ResponseWriter, r *http.Request) {
	idInstance := r.URL.Query().Get("idInstance")
	apiTokenInstance := r.URL.Query().Get("apiTokenInstance")
	url := getAPIUrl(idInstance, apiTokenInstance, "getSettings")

	result, err := fetchAPI(url, "GET", nil)
	apiResponse := APIResponse{Result: result}
	if err != nil {
		apiResponse.Error = err.Error()
	}
	json.NewEncoder(w).Encode(apiResponse)
}

func getStateInstance(w http.ResponseWriter, r *http.Request) {
	idInstance := r.URL.Query().Get("idInstance")
	apiTokenInstance := r.URL.Query().Get("apiTokenInstance")
	url := getAPIUrl(idInstance, apiTokenInstance, "getStateInstance")

	result, err := fetchAPI(url, "GET", nil)
	apiResponse := APIResponse{Result: result}
	if err != nil {
		apiResponse.Error = err.Error()
	}
	json.NewEncoder(w).Encode(apiResponse)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	idInstance := r.URL.Query().Get("idInstance")
	apiTokenInstance := r.URL.Query().Get("apiTokenInstance")
	url := getAPIUrl(idInstance, apiTokenInstance, "sendMessage")

	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := fetchAPI(url, "POST", body)
	apiResponse := APIResponse{Result: result}
	if err != nil {
		apiResponse.Error = err.Error()
	}
	json.NewEncoder(w).Encode(apiResponse)
}

func sendFileByUrl(w http.ResponseWriter, r *http.Request) {
	idInstance := r.URL.Query().Get("idInstance")
	apiTokenInstance := r.URL.Query().Get("apiTokenInstance")
	url := getAPIUrl(idInstance, apiTokenInstance, "sendFileByUrl")

	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := fetchAPI(url, "POST", body)
	apiResponse := APIResponse{Result: result}
	if err != nil {
		apiResponse.Error = err.Error()
	}
	json.NewEncoder(w).Encode(apiResponse)
}

func main() {
	http.HandleFunc("/getSettings", getSettings)
	http.HandleFunc("/getStateInstance", getStateInstance)
	http.HandleFunc("/sendMessage", sendMessage)
	http.HandleFunc("/sendFileByUrl", sendFileByUrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
