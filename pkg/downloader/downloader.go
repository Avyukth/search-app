package downloader

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

type LinkProcessor interface {
	IsLinkProcessed(ctx context.Context, id string) (bool, error)
	MarkLinkAsCompleted(ctx context.Context, id string) error
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

	err = d.processor.MarkLinkAsCompleted(ctx, id)
	if err != nil {
		return err
	}

	// Here, send the link to the processing queue
	// ...

	return nil
}
