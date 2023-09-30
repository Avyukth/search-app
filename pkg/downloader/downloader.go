package downloader

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/avyukth/search-app/pkg/queue"
)

type LinkProcessor interface {
	IsLinkProcessed(ctx context.Context, id string) (bool, error)
	MarkLinkAsCompleted(ctx context.Context, id string) error
	MarkLinkAsProcessing(ctx context.Context, id string) error
}

// Downloader is responsible for downloading and processing links.
type Downloader struct {
	processor LinkProcessor
}

// NewDownloader creates a new Downloader with the given LinkProcessor.
func NewDownloader(processor LinkProcessor) *Downloader {
	return &Downloader{processor: processor}
}

func (d *Downloader) DownloadAndProcess(ctx context.Context, link string) error {
	// Check if the link is live
	resp, err := http.Head(link)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err // or a custom error indicating the link is not live
	}

	// Hash the link to create a unique ID
	hash := sha256.Sum256([]byte(link))
	id := hex.EncodeToString(hash[:])

	// Check if the link has been processed
	processed, err := d.processor.IsLinkProcessed(ctx, id)
	if err != nil {
		return err
	}

	if processed {
		return nil
	}

	err = d.processor.MarkLinkAsProcessing(ctx, id)
	if err != nil {
		return err
	}

	// Here, send the link to the processing queue
	// TODO: Send link to processing queue

	return nil
}

func ExtractTarGz(gzipStream, dest string) error {
	reader, err := os.Open(gzipStream)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	tarReader := tar.NewReader(archive)

	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			writer, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(writer, tarReader); err != nil {
				return err
			}
			writer.Close()
		}
	}
}

func WalkDir(filePath string, queue *queue.TaskQueue, wg *sync.WaitGroup, errch chan<- error) error {
	defer wg.Done()
	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(strings.ToLower(info.Name()), ".xml") {
			return nil
		}

		// Enqueue the task to the processing queue
		task := queue.Task{FilePath: path}
		queue.Enqueue(task)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (d *Downloader) sendToProcessingQueue(ctx context.Context, link string, id string) error {
	destPath := "./downloads/" + id
	err := Download(ctx, link, destPath)
	if err != nil {
		log.Printf("Error downloading file: %v", err)
		return err
	}

	destDir := "./extracted/" + id
	err = ExtractTarGz(destPath, destDir)
	if err != nil {
		log.Printf("Error extracting tar.gz: %v", err)
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 1)
	err = WalkDir(destDir, queue, &wg, errCh)
	if err != nil {
		log.Printf("Error walking directory: %v", err)
		return err
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			log.Printf("Error processing file: %v", err)
			return err
		}
	}

	return nil
}

func Download(ctx context.Context, url string, destPath string) error {
	// Create the file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Get the data
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server return non-200 status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
