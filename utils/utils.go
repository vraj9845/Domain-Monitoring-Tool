package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Request struct {
	Headers map[string]string `yaml:"headers"`
	Method  string            `yaml:"method"`
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Body    string            `yaml:"body,omitempty"`
	Domain  string            `yaml:"domain"`
}

// GetRequestsFromFile reads a YAML file at the specified filePath and
//	returns a slice of Request structs
func GetRequestsFromFile(filePath string) ([]Request, error) {
	// Read YAML file
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading YAML file: %v", err)
	}

	// Unmarshal YAML content
	var requests []Request
	err = yaml.Unmarshal(yamlFile, &requests)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling YAML: %v", err)
	}

	return requests, nil
}

// CheckEndpointsHealths takes a slice of Request structs representing different endpoints,
// sends HTTP requests to those endpoints, measures their response times, and checks if they are healthy
func CheckEndpointsHealths(requests []Request) (map[string][]int, error) {
	// Initialize a map to store health status for each domain
	resultMap := make(map[string][]int)

	for _, req := range requests {
		// Extract domain name from the URL
		domainName, domainNameError := getDomainName(req.URL)
		if domainNameError != nil {
			return nil, fmt.Errorf("Error in getting Domain Name for URL: %v", req.URL)
		}
		req.Domain = domainName

		// Prepare HTTP request
		client := &http.Client{}
		var httpRequest *http.Request

		httpRequest, err := http.NewRequest(req.Method, req.URL, strings.NewReader(req.Body))
		if err != nil {
			return nil, fmt.Errorf("Error creating HTTP request: %v", err)
		}

		// Add headers to the request
		for key, value := range req.Headers {
			httpRequest.Header.Add(key, value)
		}

		// Execute HTTP request
		// Start measuring request time
		startTime := time.Now().UnixNano() / int64(time.Millisecond)
		response, err := client.Do(httpRequest)
		if err != nil {
			return nil, fmt.Errorf("Error executing HTTP request: %v", err)
		}
		defer response.Body.Close()

		// Stop measuring request time
		endTime := time.Now().UnixNano() / int64(time.Millisecond)

		// Check if the response status code indicates success
		if response.StatusCode >= 200 && response.StatusCode <= 299 {
			// Calculate the response time
			diff := endTime - startTime
			// If the response time is less than 500 milliseconds, the service UP
			if diff < 500 {
				resultMap[req.Domain] = append(resultMap[req.Domain], 1)
			} else {
				// The service is DOWN as response time exceeded
				resultMap[req.Domain] = append(resultMap[req.Domain], 0)
			}
		} else {
			// The service is DOWN
			resultMap[req.Domain] = append(resultMap[req.Domain], 0)
		}

	}
	return resultMap, nil
}

// VerifyFileType checks the provided file path to ensure it is not empty and has a valid YAML extension
func VerifyFileType(filePath *string) error {
	// Check if file path is provided
	if *filePath == "" {
		return fmt.Errorf("Please provide the path to the YAML file using the -file flag.")
	}

	// Check if the file has a YAML extension
	if !strings.HasSuffix(strings.ToLower(*filePath), ".yaml") && !strings.HasSuffix(strings.ToLower(*filePath), ".yml") {
		return fmt.Errorf("Invalid file format. The input file must have a YAML (.yaml or .yml) extension.")
	}

	return nil
}

// getDomainName takes a raw URL as input, parses it, and extracts the domain name
func getDomainName(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("Error parsing URL: %v", err)
	}
	return parsedURL.Hostname(), nil
}

// CalculateResult takes a map where each key is the domain name and the corresponding value is a slice of integers
// indicating the health status (1 for UP, 0 for DOWN) over multiple checks
func CalculateResult(resultMap map[string][]int) {
	for key, val := range resultMap {
		var availability int
		var availabilityPercentage float32
		// Calculate the total availability count for the domain
		for i := 0; i < len(val); i++ {
			availability = availability + val[i]
		}

		// Calculating domain availability percentage
		availabilityPercentage = float32(availability) / float32(len(val))
		fmt.Printf("%v has %.0f%% availability percentage\n", key, availabilityPercentage*100)
	}
}
