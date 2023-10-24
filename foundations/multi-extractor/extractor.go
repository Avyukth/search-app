package main

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
			// Ensure parent directories exist
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

func extractRecursively(archivePath string, destPath string) error {
	var childArchives []Archive

	// Strip the file extension from the archive name to form the extraction directory name.
	archiveNameWithoutExtension := strings.TrimSuffix(filepath.Base(archivePath), filepath.Ext(archivePath))
	extractDestPath := filepath.Join(destPath, archiveNameWithoutExtension)

	if strings.HasSuffix(strings.ToLower(archivePath), ".tar") {
		if err := extractTar(Archive{Path: archivePath, Dest: extractDestPath}); err != nil {
			return err
		}
	} else if strings.HasSuffix(strings.ToLower(archivePath), ".zip") {
		if err := extractZip(Archive{Path: archivePath, Dest: extractDestPath}); err != nil {
			return err
		}
	}

	// Check the extracted content for child archives
	err := filepath.Walk(extractDestPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(strings.ToLower(info.Name()), ".tar") || strings.HasSuffix(strings.ToLower(info.Name()), ".zip") {
			// Use the directory containing the current archive as the destination for child archives.
			childArchives = append(childArchives, Archive{Path: path, Dest: filepath.Dir(path)})
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Process child archives
	for _, childArchive := range childArchives {
		if err := extractRecursively(childArchive.Path, childArchive.Dest); err != nil {
			return err
		}
		// Remove the archive file after extracting to prevent re-processing in subsequent iterations
		os.Remove(childArchive.Path)
	}

	return nil
}

func main() {
	archivePath := "/home/amrit/Documents/100Days/golang/Search-app/tmp-out/I20230711.tar"
	destPath := "/home/amrit/Documents/100Days/golang/Search-app/tmp-out/output"

	if err := extractRecursively(archivePath, destPath); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Extraction completed successfully.")
	}
}
