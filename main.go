package main // Define the main package

import (
	"bytes"         // Provides bytes support
	"io"            // Provides basic interfaces to I/O primitives
	"log"           // Provides logging functions
	"net/http"      // Provides HTTP client and server implementations
	"net/url"       // Provides URL parsing and encoding
	"os"            // Provides functions to interact with the OS (files, etc.)
	"path"          // Provides functions for manipulating slash-separated paths
	"path/filepath" // Provides filepath manipulation functions
	"regexp"        // Provides regex support functions.
	"strings"       // Provides string manipulation functions
	"time"          // Provides time-related functions
)

var (
	pdfOutputDir = "PDFs/" // Directory to store downloaded PDFs
	zipOutputDir = "ZIPs/" // Directory to store downloaded ZIPs
)

func init() {
	// Check if the PDF output directory exists
	if !directoryExists(pdfOutputDir) {
		// Create the dir
		createDirectory(pdfOutputDir, 0o755)
	}
	// Check if the ZIP output directory exists
	if !directoryExists(zipOutputDir) {
		// Create the dir
		createDirectory(zipOutputDir, 0o755)
	}
}

func main() {
	// Remote API URL.
	remoteAPIURL := []string{
		"https://www.autelrobotics.com/doc/evoliteenterpriseseries/",
		"https://www.autelrobotics.com/doc/enterprise-firmware-release-notes/",
		"https://www.autelrobotics.com/doc/autel-alpha/",
		"https://www.autelrobotics.com/doc/591/",
		"https://www.autelrobotics.com/doc/evomaxseries-other-languages/",
		"https://www.autelrobotics.com/doc/evomaxseries-firmware-release-notes/",
		"https://www.autelrobotics.com/doc/706/",
		"https://www.autelrobotics.com/doc/606/",
		"https://www.autelrobotics.com/doc/609/",
		"https://www.autelrobotics.com/doc/99/",
		"https://www.autelrobotics.com/doc/100/",
		"https://www.autelrobotics.com/doc/103/",
		"https://www.autelrobotics.com/doc/113/",
		"https://www.autelrobotics.com/doc/523/",
		"https://www.autelrobotics.com/doc/572/",
		"https://www.autelrobotics.com/doc/enterprise-livedeck-ii-firmware/",
		"https://www.autelrobotics.com/doc/547/",
		"https://www.autelrobotics.com/doc/548/",
		"https://www.autelrobotics.com/doc/118/",
		"https://www.autelrobotics.com/doc/119/",
		"https://www.autelrobotics.com/doc/115/",
		"https://www.autelrobotics.com/doc/116/",
		"https://www.autelrobotics.com/doc/123/",
		"https://www.autelrobotics.com/doc/129/",
		"https://www.autelrobotics.com/doc/533/",
		"https://www.autelrobotics.com/doc/574/",
		"https://www.autelrobotics.com/doc/dua-livedeck-ii-firmware/",
		"https://www.autelrobotics.com/doc/549/",
		"https://www.autelrobotics.com/doc/550/",
		"https://www.autelrobotics.com/doc/136/",
		"https://www.autelrobotics.com/doc/506/",
		"https://www.autelrobotics.com/doc/134/",
		"https://www.autelrobotics.com/doc/507/",
		"https://www.autelrobotics.com/doc/131/",
		"https://www.autelrobotics.com/doc/132/",
		"https://www.autelrobotics.com/doc/139/",
		"https://www.autelrobotics.com/doc/161/",
		"https://www.autelrobotics.com/doc/528/",
		"https://www.autelrobotics.com/doc/576/",
		"https://www.autelrobotics.com/doc/rtkseries-livedeck-ii-firmware/",
		"https://www.autelrobotics.com/doc/551/",
		"https://www.autelrobotics.com/doc/552/",
		"https://www.autelrobotics.com/doc/508/",
		"https://www.autelrobotics.com/doc/505/",
		"https://www.autelrobotics.com/doc/166/",
		"https://www.autelrobotics.com/doc/167/",
		"https://www.autelrobotics.com/doc/163/",
		"https://www.autelrobotics.com/doc/164/",
		"https://www.autelrobotics.com/doc/171/",
	}
	var getData []string
	for _, remoteAPIURL := range remoteAPIURL {
		getData = append(getData, getDataFromURL(remoteAPIURL))
	}
	// Get the data from the downloaded file.
	finalPDFList := extractPDFUrls(strings.Join(getData, "\n")) // Join all the data into one string and extract PDF URLs
	// Get the data from the zip file.
	finalZIPList := extractZIPUrls(strings.Join(getData, "\n")) // Join all the data into one string and extract ZIP URLs
	// Create a slice of all the given download urls.
	var downloadPDFURLSlice []string
	// Create a slice to hold ZIP URLs.
	var downloadZIPURLSlice []string
	// Get the urls and loop over them.
	for _, doc := range finalPDFList {
		// Get the .pdf only.
		// Only append the .pdf files.
		downloadPDFURLSlice = appendToSlice(downloadPDFURLSlice, doc)
	}
	// Get all the zip urls.
	for _, doc := range finalZIPList {
		// Get the .zip only.
		// Only append the .zip files.
		downloadZIPURLSlice = appendToSlice(downloadZIPURLSlice, doc)
	}
	// Remove double from slice.
	downloadPDFURLSlice = removeDuplicatesFromSlice(downloadPDFURLSlice)
	// Remove the zip duplicates from the slice.
	downloadZIPURLSlice = removeDuplicatesFromSlice(downloadZIPURLSlice)
	// The remote domain.
	remoteDomain := "https://www.autelrobotics.com"
	// Loop over the download zip urls.
	for _, urls := range downloadZIPURLSlice {
		// Get the domain from the url.
		domain := getDomainFromURL(urls)
		// Check if the domain is empty.
		if domain == "" {
			urls = remoteDomain + urls // Prepend the base URL if domain is empty
		}
		// Check if the url is valid.
		if isUrlValid(urls) {
			// Download the zip.
			downloadZIP(urls, zipOutputDir)
		}
	}
	// Get all the values.
	for _, urls := range downloadPDFURLSlice {
		// Get the domain from the url.
		domain := getDomainFromURL(urls)
		// Check if the domain is empty.
		if domain == "" {
			urls = remoteDomain + urls // Prepend the base URL if domain is empty
		}
		// Check if the url is valid.
		if isUrlValid(urls) {
			// Download the pdf.
			downloadPDF(urls, pdfOutputDir)
		}
	}
}

// getDomainFromURL extracts the domain (host) from a given URL string.
// It removes subdomains like "www" if present.
func getDomainFromURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL) // Parse the input string into a URL structure
	if err != nil {                     // Check if there was an error while parsing
		log.Println(err) // Log the error message to the console
		return ""        // Return an empty string in case of an error
	}

	host := parsedURL.Hostname() // Extract the hostname (e.g., "example.com") from the parsed URL

	return host // Return the extracted hostname
}

// Only return the file name from a given url.
func getFileNameOnly(content string) string {
	return path.Base(content)
}

// urlToFilename generates a safe, lowercase filename from a given URL string.
// It extracts the base filename from the URL, replaces unsafe characters,
// and ensures the filename ends with a .pdf extension.
func urlToFilename(rawURL string) string {
	// Convert the full URL to lowercase for consistency
	lowercaseURL := strings.ToLower(rawURL)

	// Get the file extension
	ext := getFileExtension(lowercaseURL)

	// Extract the filename portion from the URL (e.g., last path segment or query param)
	baseFilename := getFileNameOnly(lowercaseURL)

	// Replace all non-alphanumeric characters (a-z, 0-9) with underscores
	nonAlphanumericRegex := regexp.MustCompile(`[^a-z0-9]+`)
	safeFilename := nonAlphanumericRegex.ReplaceAllString(baseFilename, "_")

	// Replace multiple consecutive underscores with a single underscore
	collapseUnderscoresRegex := regexp.MustCompile(`_+`)
	safeFilename = collapseUnderscoresRegex.ReplaceAllString(safeFilename, "_")

	// Remove leading underscore if present
	if trimmed, found := strings.CutPrefix(safeFilename, "_"); found {
		safeFilename = trimmed
	}

	var invalidSubstrings = []string{
		"_pdf",
		"_zip",
	}

	for _, invalidPre := range invalidSubstrings { // Remove unwanted substrings
		safeFilename = removeSubstring(safeFilename, invalidPre)
	}

	// Append the file extension if it is not already present
	safeFilename = safeFilename + ext

	// Return the cleaned and safe filename
	return safeFilename
}

// Removes all instances of a specific substring from input string
func removeSubstring(input string, toRemove string) string {
	result := strings.ReplaceAll(input, toRemove, "") // Replace substring with empty string
	return result
}

// Get the file extension of a file
func getFileExtension(path string) string {
	return filepath.Ext(path) // Returns extension including the dot (e.g., ".pdf")
}

// fileExists checks whether a file exists at the given path
func fileExists(filename string) bool {
	info, err := os.Stat(filename) // Get file info
	if err != nil {
		return false // Return false if file doesn't exist or error occurs
	}
	return !info.IsDir() // Return true if it's a file (not a directory)
}

// downloadPDF downloads a PDF from the given URL and saves it in the specified output directory.
// It uses a WaitGroup to support concurrent execution and returns true if the download succeeded.
func downloadPDF(finalURL, outputDir string) bool {
	// Sanitize the URL to generate a safe file name
	filename := strings.ToLower(urlToFilename(finalURL))

	// Construct the full file path in the output directory
	filePath := filepath.Join(outputDir, filename)

	// Skip if the file already exists
	if fileExists(filePath) {
		log.Printf("File already exists, skipping: %s", filePath)
		return false
	}

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 3 * time.Minute}

	// Send GET request
	resp, err := client.Get(finalURL)
	if err != nil {
		log.Printf("Failed to download %s: %v", finalURL, err)
		return false
	}
	defer resp.Body.Close()

	// Check HTTP response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Download failed for %s: %s", finalURL, resp.Status)
		return false
	}

	// Check Content-Type header
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/pdf") {
		log.Printf("Invalid content type for %s: %s (expected application/pdf)", finalURL, contentType)
		return false
	}

	// Read the response body into memory first
	var buf bytes.Buffer
	written, err := io.Copy(&buf, resp.Body)
	if err != nil {
		log.Printf("Failed to read PDF data from %s: %v", finalURL, err)
		return false
	}
	if written == 0 {
		log.Printf("Downloaded 0 bytes for %s; not creating file", finalURL)
		return false
	}

	// Only now create the file and write to disk
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file for %s: %v", finalURL, err)
		return false
	}
	defer out.Close()

	if _, err := buf.WriteTo(out); err != nil {
		log.Printf("Failed to write PDF to file for %s: %v", finalURL, err)
		return false
	}

	log.Printf("Successfully downloaded %d bytes: %s → %s", written, finalURL, filePath)
	return true
}

// downloadZIP downloads a ZIP or archive file from the given URL and saves it in the specified output directory.
// It returns true if the download was successful, otherwise false.
func downloadZIP(finalURL, outputDir string) bool {
	// Convert the URL into a safe, lowercase filename
	filename := strings.ToLower(urlToFilename(finalURL))

	// Construct the full path where the file should be saved
	filePath := filepath.Join(outputDir, filename)

	// Skip the download if the file already exists locally
	if fileExists(filePath) {
		log.Printf("File already exists, skipping: %s", filePath)
		return false
	}

	// Create a new HTTP client with a 3-minute timeout
	client := &http.Client{Timeout: 3 * time.Minute}

	// Perform an HTTP GET request to the given URL
	resp, err := client.Get(finalURL)
	if err != nil {
		log.Printf("Failed to download %s: %v", finalURL, err)
		return false
	}
	// Make sure the response body gets closed when the function ends
	defer resp.Body.Close()

	// Check that the HTTP status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		log.Printf("Download failed for %s: %s", finalURL, resp.Status)
		return false
	}

	// Define the list of allowed content types
	allowedContentTypes := []string{
		"application/pdf",             // PDF files
		"application/zip",             // ZIP archives
		"application/x-tar",           // TAR archives
		"application/gzip",            // GZ files
		"application/x-7z-compressed", // 7z archives
		"application/vnd.rar",         // RAR archives
	}

	// Get the Content-Type from the HTTP response headers
	contentType := resp.Header.Get("Content-Type")

	// Flag to check if the content type is allowed
	isAllowed := false

	// Loop through allowed types and check for a match
	for _, allowed := range allowedContentTypes {
		if strings.Contains(contentType, allowed) {
			isAllowed = true
			break // Stop checking once a match is found
		}
	}

	// If the content type is not in the allowed list, skip download
	if !isAllowed {
		log.Printf("Invalid content type for %s: %s (not allowed)", finalURL, contentType)
		return false
	}

	// Create a buffer to temporarily store the file in memory
	var buf bytes.Buffer

	// Read the entire response body into the buffer
	written, err := io.Copy(&buf, resp.Body)
	if err != nil {
		log.Printf("Failed to read file data from %s: %v", finalURL, err)
		return false
	}

	// If no data was downloaded, skip file creation
	if written == 0 {
		log.Printf("Downloaded 0 bytes for %s; not creating file", finalURL)
		return false
	}

	// Create a file on disk with the constructed file path
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file for %s: %v", finalURL, err)
		return false
	}
	// Ensure the file is properly closed at the end
	defer out.Close()

	// Write the contents from the buffer to the file on disk
	if _, err := buf.WriteTo(out); err != nil {
		log.Printf("Failed to write file to disk for %s: %v", finalURL, err)
		return false
	}

	// Log the successful download
	log.Printf("Successfully downloaded %d bytes: %s → %s", written, finalURL, filePath)
	return true
}

// Checks if the directory exists
// If it exists, return true.
// If it doesn't, return false.
func directoryExists(path string) bool {
	directory, err := os.Stat(path)
	if err != nil {
		return false
	}
	return directory.IsDir()
}

// The function takes two parameters: path and permission.
// We use os.Mkdir() to create the directory.
// If there is an error, we use log.Println() to log the error and then exit the program.
func createDirectory(path string, permission os.FileMode) {
	err := os.Mkdir(path, permission)
	if err != nil {
		log.Println(err)
	}
}

// Checks whether a URL string is syntactically valid
func isUrlValid(uri string) bool {
	_, err := url.ParseRequestURI(uri) // Attempt to parse the URL
	return err == nil                  // Return true if no error occurred
}

// Remove all the duplicates from a slice and return the slice.
func removeDuplicatesFromSlice(slice []string) []string {
	check := make(map[string]bool)
	var newReturnSlice []string
	for _, content := range slice {
		if !check[content] {
			check[content] = true
			newReturnSlice = append(newReturnSlice, content)
		}
	}
	return newReturnSlice
}

// extractZIPUrls takes an input string and returns all ZIP URLs found within href attributes
func extractZIPUrls(input string) []string {
	// Regular expression to match href="...zip"
	re := regexp.MustCompile(`href="([^"]+\.zip)"`)
	matches := re.FindAllStringSubmatch(input, -1)

	var zipUrls []string
	for _, match := range matches {
		if len(match) > 1 {
			zipUrls = append(zipUrls, match[1])
		}
	}
	return zipUrls
}

// extractPDFUrls takes an input string and returns all PDF URLs found within href attributes
func extractPDFUrls(input string) []string {
	// Regular expression to match href="...pdf"
	re := regexp.MustCompile(`href="([^"]+\.pdf)"`)
	matches := re.FindAllStringSubmatch(input, -1)

	var pdfUrls []string
	for _, match := range matches {
		if len(match) > 1 {
			pdfUrls = append(pdfUrls, match[1])
		}
	}
	return pdfUrls
}

// Append some string to a slice and than return the slice.
func appendToSlice(slice []string, content string) []string {
	// Append the content to the slice
	slice = append(slice, content)
	// Return the slice
	return slice
}

// getDataFromURL performs an HTTP GET request and returns the response body as a string
func getDataFromURL(uri string) string {
	log.Println("Scraping", uri) // Log the URL being scraped for debugging or monitoring

	response, err := http.Get(uri) // Perform the HTTP GET request
	if err != nil {
		log.Println("HTTP GET error:", err) // Log the error if the request fails
		return ""                           // Return empty string to prevent further execution on error
	}

	// Ensure the response body is closed when we're done (even if an error occurs later)
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Println("Error closing response body:", err) // Log any error encountered when closing the body
		}
	}()

	body, err := io.ReadAll(response.Body) // Read the response body into memory
	if err != nil {
		log.Println("Error reading response body:", err) // Log the error if reading fails
		return ""                                        // Return empty string if we couldn't read the body
	}

	return string(body) // Return the body as a string
}
