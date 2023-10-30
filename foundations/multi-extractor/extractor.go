package main

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Archive struct {
	Path string
	Dest string
}

func extractTar(archive Archive) error {
	file, err := os.Open(archive.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := tar.NewReader(file)

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(archive.Dest, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			fmt.Printf("Creating directory: %s\n", targetPath)
			if err := os.MkdirAll(targetPath, info.Mode()); err != nil {
				return err
			}
			fmt.Println("Created directory:", targetPath)
		} else {
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			fmt.Printf("Creating file: %s\n", targetPath)
			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, reader); err != nil {
				outFile.Close()
				return err
			}
			fmt.Println("Extracted file:", targetPath)
			outFile.Close()
		}
	}
	return nil
}

func extractZip(archive Archive) error {
	reader, err := zip.OpenReader(archive.Path)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		targetPath := filepath.Join(archive.Dest, file.Name)
		if file.FileInfo().IsDir() {
			fmt.Printf("Creating directory (zip): %s\n", targetPath)
			if err := os.MkdirAll(targetPath, file.Mode()); err != nil {
				return err
			}
			fmt.Println("Created directory:", targetPath)
		} else {
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			fmt.Printf("Creating file (zip): %s\n", targetPath)
			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			rc, err := file.Open()
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, rc); err != nil {
				rc.Close()
				outFile.Close()
				return err
			}
			fmt.Println("Extracted file:", targetPath)
			rc.Close()
			outFile.Close()
		}
	}
	return nil
}

func extractRecursively(archivePath string, destPath string, doneCh chan bool, errCh chan error) {
	var childArchives []Archive

	archiveNameWithoutExtension := strings.TrimSuffix(filepath.Base(archivePath), filepath.Ext(archivePath))
	extractDestPath := filepath.Join(destPath, archiveNameWithoutExtension)

	var err error
	if strings.HasSuffix(strings.ToLower(archivePath), ".tar") {
		err = extractTar(Archive{Path: archivePath, Dest: extractDestPath})
	} else if strings.HasSuffix(strings.ToLower(archivePath), ".zip") {
		err = extractZip(Archive{Path: archivePath, Dest: extractDestPath})
	}

	if err != nil {
		errCh <- err
		doneCh <- true
		return
	}

	err = filepath.Walk(extractDestPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".tar") || strings.HasSuffix(strings.ToLower(info.Name()), ".zip") {
			childArchives = append(childArchives, Archive{Path: path, Dest: filepath.Dir(path)})
		}
		return nil
	})

	if err != nil {
		errCh <- err
		doneCh <- true
		return
	}

	wg := &sync.WaitGroup{}
	for _, childArchive := range childArchives {
		wg.Add(1)
		go func(archive Archive) {
			defer wg.Done()
			extractRecursively(archive.Path, archive.Dest, doneCh, errCh)
		}(childArchive)
		os.Remove(childArchive.Path)
	}

	// Wait for child extractions to complete
	wg.Wait()

	doneCh <- true
}

func main() {
	archivePath := "/home/amrit/Documents/100Days/golang/Search-app/tmp-out/I20230711.tar"
	destPath := "/home/amrit/Documents/100Days/golang/Search-app/tmp-out/output"

	doneCh := make(chan bool, 10)
	errCh := make(chan error, 10)

	go extractRecursively(archivePath, destPath, doneCh, errCh)

	// Monitor channels
	for {
		select {
		case <-doneCh:
			return
		case err := <-errCh:
			if err != nil {
				fmt.Println("Error:", err)
			}
		}
	}
}
