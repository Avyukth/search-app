package downloader

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/avyukth/search-app/pkg/config"
)

type Downloader interface {
	Download(ctx context.Context, link string) (string, error)
	ExtractTarGz(filePath string) (string, error)
}

type HTTPDownloader struct {
	client       *http.Client
	serverConfig *config.ServerConfig
}

func NewDownloader(client *http.Client, serverConfig *config.ServerConfig) Downloader {
	return &HTTPDownloader{
		client:       client,
		serverConfig: serverConfig,
	}
}

func (d *HTTPDownloader) Download(ctx context.Context, link string) (string, error) {
	log.Printf("Starting download for link: %s", link)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return "", fmt.Errorf("creating request for link %s: %w", link, err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing request for link %s: %w", link, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 status code received for link %s: %s", link, resp.Status)
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current working directory: %w", err)
	}

	destFolder := filepath.Join(dir, d.serverConfig.DataStoreDirectory)
	if err := os.MkdirAll(destFolder, 0755); err != nil {
		return "", fmt.Errorf("creating directory %s: %w", destFolder, err)
	}

	destPath := filepath.Join(destFolder, "abc.tar.gz")
	out, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("creating file at %s: %w", destPath, err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("writing to file at %s: %w", destPath, err)
	}

	log.Printf("Successfully downloaded and wrote to file at: %s", destPath)
	return destPath, nil
}

func (d *HTTPDownloader) ExtractTarGz(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("opening file %s: %w", filePath, err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return "", fmt.Errorf("creating gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	destDir := filepath.Dir(filePath)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("reading tar header: %w", err)
		}

		target := filepath.Join(destDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return "", fmt.Errorf("creating directory %s: %w", target, err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return "", fmt.Errorf("creating file %s: %w", target, err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return "", fmt.Errorf("writing to file %s: %w", target, err)
			}
			outFile.Close()
		}
	}
	return destDir, nil
}
