// package main

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"sync"
// )

// func main() {
// 	url := "https://bulkdata.uspto.gov/data/patent/grant/redbook/2023/I20230502.tar"

// 	numConnections := 200

// 	// Create a WaitGroup to wait for all goroutines to finish
// 	var wg sync.WaitGroup

// 	// Create a channel to collect downloaded file parts
// 	fileParts := make(chan []byte, numConnections)

// 	// Create a channel to signal download errors
// 	errCh := make(chan error, numConnections)

// 	// Calculate the size of each file part
// 	fileSize, err := getFileSize(url)
// 	if err != nil {
// 		fmt.Println("Error getting file size:", err)
// 		return
// 	}

// 	fmt.Println("fileSize", fileSize)

// 	partSize := fileSize / int64(numConnections)

// 	// Get the original file's extension from the URL
// 	fileExt := filepath.Ext(url)
// 	if fileExt == "" {
// 		fileExt = ".data" // Default extension if not found
// 	}

// 	// Create and start goroutines for concurrent downloads
// 	for i := 0; i < numConnections; i++ {
// 		startByte := int64(i) * partSize
// 		endByte := int64((i + 1)) * partSize

// 		// Create a new goroutine for downloading a file part
// 		wg.Add(1)
// 		go func(start, end int64, partNum int) {
// 			defer wg.Done()

// 			part, err := downloadPart(url, start, end)
// 			if err != nil {
// 				fmt.Println("Error downloading part:", err)
// 				errCh <- err
// 				return
// 			}

// 			fmt.Printf("Downloaded part %d\n", partNum)
// 			fileParts <- part
// 		}(startByte, endByte, i+1)
// 	}

// 	// Close the channels when all downloads are done
// 	go func() {
// 		wg.Wait()
// 		close(fileParts)
// 		close(errCh)
// 	}()

// 	// Create the output file with a ".part" extension
// 	dir, _ := os.Getwd()
// 	outputFilePath := filepath.Join(dir, "output"+fileExt+".part")
// 	outputFile, err := os.Create(outputFilePath)
// 	if err != nil {
// 		fmt.Println("Error creating output file:", err)
// 		return
// 	}
// 	defer outputFile.Close()

// 	// Write the downloaded parts to the output file
// 	for part := range fileParts {
// 		_, err := outputFile.Write(part)
// 		if err != nil {
// 			fmt.Println("Error writing part to output file:", err)
// 			return
// 		}
// 	}

// 	// Check for any download errors
// 	if err := <-errCh; err != nil {
// 		fmt.Println("Download failed:", err)
// 		return
// 	}

// 	// Rename the final file to remove the ".part" extension
// 	finalOutputFilePath := filepath.Join(dir, "output"+fileExt)
// 	err = os.Rename(outputFilePath, finalOutputFilePath)
// 	if err != nil {
// 		fmt.Println("Error renaming final output file:", err)
// 		return
// 	}

// 	fmt.Println("File downloaded successfully:", finalOutputFilePath)
// }

// // Rest of your code remains the same

// func getFileSize(url string) (int64, error) {
// 	resp, err := http.Head(url)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer resp.Body.Close()

// 	fmt.Println("resp.ContentLength", resp.ContentLength)
// 	return resp.ContentLength, nil
// }

// func downloadPart(url string, start, end int64) ([]byte, error) {
// 	client := &http.Client{}

// 	// Create a request for the specific range of bytes
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	part, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println("part", part)
// 	return part, nil
// }

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	url := "https://bulkdata.uspto.gov/data/patent/grant/redbook/2023/I20230502.tar"
	numConnections := 200

	var wg sync.WaitGroup
	errCh := make(chan error, numConnections)
	doneCh := make(chan bool, numConnections)
	tmpFiles := make([]string, numConnections)

	fileSize, err := getFileSize(url)
	if err != nil {
		fmt.Println("Error getting file size:", err)
		return
	}

	partSize := fileSize / int64(numConnections)

	// Create temporary directory in the current working directory
	dir, _ := os.Getwd()
	tmpDir, err := os.MkdirTemp(dir, "tmp-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}

	for i := 0; i < numConnections; i++ {
		startByte := int64(i) * partSize
		endByte := startByte + partSize - 1
		if i == numConnections-1 {
			endByte = fileSize - 1
		}

		wg.Add(1)
		go func(start, end int64, partNum int) {
			defer wg.Done()

			part, err := downloadPart(url, start, end, partNum)
			if err != nil {
				errCh <- err
				return
			}
			fmt.Println("got part")
			tmpFile, err := os.Create(filepath.Join(tmpDir, fmt.Sprintf("part-%d", partNum)))
			if err != nil {
				errCh <- err
				return
			}

			_, err = tmpFile.Write(part)
			if err != nil {
				errCh <- err
				return
			}

			tmpFiles[partNum-1] = tmpFile.Name()
			tmpFile.Close()
			doneCh <- true
		}(startByte, endByte, i+1)
	}

	go func() {
		wg.Wait()
		close(doneCh)
		close(errCh)
	}()

	// Wait until all goroutines have signaled their completion
	completedDownloads := 0
	for range doneCh {
		completedDownloads++
		if completedDownloads == numConnections {
			break
		}
	}

	fileExt := filepath.Ext(url)
	finalOutputFilePath := filepath.Join(dir, "output"+fileExt)
	finalOutputFile, err := os.Create(finalOutputFilePath)
	if err != nil {
		fmt.Println("Error creating final output file:", err)
		return
	}
	defer finalOutputFile.Close()

	// Once all parts are downloaded, combine the files
	for _, tmpFileName := range tmpFiles {
		tmpFile, err := os.Open(tmpFileName)
		if err != nil {
			fmt.Println("Error reading temporary file:", err)
			return
		}

		_, err = io.Copy(finalOutputFile, tmpFile)
		if err != nil {
			fmt.Println("Error writing to final output file:", err)
			tmpFile.Close()
			return
		}

		tmpFile.Close()
		os.Remove(tmpFileName)
	}

	// Remove temporary directory
	os.Remove(tmpDir)

	if err, ok := <-errCh; ok {
		fmt.Println("Download failed:", err)
		return
	}

	fmt.Println("File downloaded successfully:", finalOutputFilePath)
}

func getFileSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.ContentLength, nil
}

// func downloadPart(url string, start, end int64) ([]byte, error) {
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
// 	fmt.Println("req range : ", fmt.Sprintf("bytes=%d-%d", start, end))
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	return io.ReadAll(resp.Body)
// }

func downloadPart(url string, start, end int64, partNum int) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	startTime := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	part, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()

	partSize := float64(end - start + 1)
	speed := partSize / duration
	speedKB := speed / 1024

	var speedStr string
	if speedKB < 1024 {
		speedStr = fmt.Sprintf("%.2f KB/sec", speedKB)
	} else if speedKB >= 1024 && speedKB < 1024*1024 {
		speedMB := speedKB / 1024
		speedStr = fmt.Sprintf("%.2f MB/sec", speedMB)
	} else {
		speedMbit := speedKB * 8 / 1024
		speedStr = fmt.Sprintf("%.2f Mbit/sec", speedMbit)
	}

	fmt.Printf("Part %d downloaded in %.2f seconds at %s\n", partNum, duration, speedStr)

	return part, nil
}
