package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"
)

func speedTestHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Received request")
	serverURL := "https://research.nhm.org/pdfs/10840/10840.pdf" // Replace with a large file URL for testing

	// Start a timer
	startTime := time.Now()

	// Make an HTTP GET request to download a large file
	response, err := http.Get(serverURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Create a temporary file to store the downloaded data
	tmpFile, err := os.CreateTemp("", "download-*.tmp")
	if err != nil {
		fmt.Printf("Error creating temporary file: %v\n", err)
		http.Error(w, "Error creating temporary file", http.StatusInternalServerError)
		return
	}
	defer tmpFile.Close()

	// Copy the response body to the temporary file
	_, err = io.Copy(tmpFile, response.Body)
	if err != nil {
		fmt.Printf("Error copying data: %v\n", err)
		http.Error(w, "Error copying data", http.StatusInternalServerError)
		return
	}

	// Stop the timer
	duration := time.Since(startTime)

	// Calculate download speed in Mbps
	fileSizeBytes, _ := tmpFile.Seek(0, io.SeekEnd)
	fileSizeBits := fileSizeBytes * 8
	durationSeconds := duration.Seconds()
	downloadSpeed := float64(fileSizeBits) / durationSeconds / 1000000
	fmt.Printf("Download Speed: %.2f Mbps \n", downloadSpeed)
	// conver download speed to 2 decimal places

	// Delete the temporary file
	tmpFile.Close()
	err = os.Remove(tmpFile.Name())
	if err != nil {
		fmt.Printf("Error deleting temporary file: %v\n", err)
		http.Error(w, "Error deleting temporary file", http.StatusInternalServerError)
		return
	}

	// Create a template to render the HTML response
	tmpl, err := template.New("speedtest").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Internet Speed Test</title>
	</head>
	<body>
		<h1>Internet Speed Test</h1>
		<p>Download Speed: {{printf "%.2f" .DownloadSpeed}} Mbps</p>
	</body>
	</html>
	`)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Render the HTML response with the download speed
	data := struct {
		DownloadSpeed float64
	}{
		DownloadSpeed: downloadSpeed,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering HTML", http.StatusInternalServerError)
		return
	}

}

func main() {
	// Serve the speed test handler
	http.HandleFunc("/speedtest", speedTestHandler)

	// Serve the static index.html file
	http.Handle("/", http.FileServer(http.Dir("static")))

	// Start the HTTP server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
