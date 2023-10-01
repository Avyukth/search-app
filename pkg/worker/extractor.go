package worker

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExtractAndDeleteArchive(archivePath string, info os.FileInfo) error {
	switch {
	case strings.HasSuffix(strings.ToLower(info.Name()), ".zip"):
		return extractZipAndDelete(archivePath)
	case strings.HasSuffix(strings.ToLower(info.Name()), ".tar"):
		return extractTarAndDelete(archivePath, false)
	case strings.HasSuffix(strings.ToLower(info.Name()), ".tar.gz"):
		return extractTarAndDelete(archivePath, true)
	}
	return fmt.Errorf("unsupported archive format: %s", archivePath)
}

// extractZipAndDelete function
func extractZipAndDelete(zipPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	extractDir := filepath.Dir(zipPath)
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		extractPath := filepath.Join(extractDir, f.Name)
		if !strings.HasPrefix(extractPath, filepath.Clean(extractDir)+string(os.PathSeparator)) {
			rc.Close()
			return fmt.Errorf("%s: illegal file path", extractPath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(extractPath, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(extractPath), f.Mode())
			file, err := os.OpenFile(extractPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				rc.Close()
				return err
			}
			_, err = io.Copy(file, rc)
			file.Close()
			rc.Close()
			if err != nil {
				return err
			}
		}
	}

	return os.Remove(zipPath)
}

// extractTarAndDelete function
func extractTarAndDelete(tarPath string, isGzipped bool) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var r io.Reader = file
	if isGzipped {
		r, err = gzip.NewReader(file)
		if err != nil {
			return err
		}
	}

	extractDir := filepath.Dir(tarPath)
	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		extractPath := filepath.Join(extractDir, header.Name)
		if !strings.HasPrefix(extractPath, filepath.Clean(extractDir)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", extractPath)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(extractPath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			file, err := os.OpenFile(extractPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(file, tr); err != nil {
				file.Close()
				return err
			}
			file.Close()
		}
	}

	return os.Remove(tarPath)
}
