package main

import (
	"flag"
	"time"

	dmtUtils "domain-monitoring-tool/utils"
)

// cmd to run program: go run your_program.go -file path/to/your_file.yaml

// This go program performs the following steps:
// 1. Parse command line arguments to get the YAML file path.
// 2. Verify the file path for existence and valid YAML extension.
// 3. Extract requests from the YAML file.
// 4. Continuously check endpoint health every 15 seconds.
// 5. Print the availability percentage for each domain based on health checks.
func main() {
	// Parse command line arguments
	filePath := flag.String("file", "", "Path to YAML file")
	flag.Parse()

	// Verify the validity of the provided file
	verifyFileError := dmtUtils.VerifyFileType(filePath)
	if verifyFileError != nil {
		panic(verifyFileError)
	}

	// Get requests from the specified YAML file
	requests, getReqError := dmtUtils.GetRequestsFromFile(*filePath)
	if getReqError != nil {
		panic(getReqError)
	}

	// Check the health of endpoints every 15 seconds
	for {
		resultMap, checkEndpointsError := dmtUtils.CheckEndpointsHealths(requests)
		if checkEndpointsError != nil {
			panic(checkEndpointsError)
		}
		// Calculate and print the availability percentage for each domain
		dmtUtils.CalculateResult(resultMap)

		// Sleep for 15 seconds before the next iteration
		time.Sleep(15 * time.Second)
	}
}
