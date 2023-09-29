package downloader

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

// LinkProcessor is an interface that defines methods to interact with the storage of processed links.
type LinkProcessor interface {
	IsLinkProcessed(ctx context.Context, id string) (bool, error)
	MarkLinkAsProcessed(ctx context.Context, id string) error
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
		return nil // or a custom error indicating the link has been processed
	}

	// Send the link to the processing queue and mark it as processed in the database
	err = d.processor.MarkLinkAsProcessed(ctx, id)
	if err != nil {
		return err
	}

	// Here, send the link to the processing queue
	// ...

	return nil
}
