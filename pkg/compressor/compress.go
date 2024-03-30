package compressor

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/s2"
)

func CompressS2Bytes(rawBytes []byte) (compressedBytes []byte, err error) {
	var buf bytes.Buffer
	c := s2.NewWriter(&buf)
	if _, err := c.Write(rawBytes); err != nil {
		return nil, err
	}
	if err := c.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecompressS2Bytes(compressedBytes []byte) (rawBytes []byte, err error) {
	if len(compressedBytes) == 0 {
		return []byte{}, nil
	}
	reader := s2.NewReader(bytes.NewReader(compressedBytes))

	rawBytes, err = io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return rawBytes, nil

}

func TryCompressS2Bytes(rawBytes []byte) (compressedBytes []byte) {
	var err error
	compressedBytes, err = CompressS2Bytes(rawBytes)
	if err != nil {
		return rawBytes
	}
	return compressedBytes
}

func TryDecompressS2Bytes(compressedBytes []byte) (rawBytes []byte) {
	var err error
	rawBytes, err = DecompressS2Bytes(compressedBytes)
	if err != nil {
		return compressedBytes
	}
	return rawBytes
}

func CompressS2Str(rawStr string) (out []byte, err error) {
	rawBytes := []byte(rawStr)
	return CompressS2Bytes(rawBytes)
}

func DecompressS2Str(compressed []byte) (string, error) {
	rawBytes, err := DecompressS2Bytes(compressed)
	if err != nil {
		return "", err
	}
	return string(rawBytes), nil
}

// TryCompressS2Str try to compress the string with S2, but if it fails, return []byte(in)
func TryCompressS2Str(rawStr string) (compressedBytes []byte) {
	rawBytes := []byte(rawStr)
	return TryCompressS2Bytes(rawBytes)
}

// TryDecompressS2Str try to decompress the []byte with S2, but if it fails, return string(in)
func TryDecompressS2Str(in []byte) (out string) {
	return string(TryDecompressS2Bytes(in))
}

func CreateTgzArchive(sourceDir, outputFile string) error {
	// Create the output directory if it doesn't exist
	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Create a new tar.gz file for writing
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		_ = outputFile.Close()
	}(f)

	// Create a gzip writer
	gzipWriter := gzip.NewWriter(f)
	defer func(gzipWriter *gzip.Writer) {
		_ = gzipWriter.Close()
	}(gzipWriter)

	// Create a tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer func(tarWriter *tar.Writer) {
		_ = tarWriter.Close()
	}(tarWriter)

	// Recursively add the sourceDir content to the tar archive
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the relative path within the source directory
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Get the last directory name (basename)
		lastDir := filepath.Base(sourceDir)

		// Modify header name to include the last directory
		headerName := filepath.Join(lastDir, relPath)

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(headerName)

		// Write header to tar archive
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// If it's a regular file, copy its content to the tar archive
		if info.Mode().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			if _, err := io.Copy(tarWriter, file); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func ExtractTgzArchive(archivePath, destDir string) error {
	// Open the input .tgz file for reading
	inputFile, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer func(inputFile *os.File) {
		_ = inputFile.Close()
	}(inputFile)

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(inputFile)
	if err != nil {
		return err
	}
	defer func(gzipReader *gzip.Reader) {
		_ = gzipReader.Close()
	}(gzipReader)

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Calculate the destination path
		destPath := filepath.Join(destDir, header.Name)

		// Create directories if necessary
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(destPath, header.FileInfo().Mode()); err != nil {
				return err
			}
			continue
		}

		// Create the file
		file, err := os.Create(destPath)
		if err != nil {
			return err
		}

		// Copy content to the file
		if _, err := io.Copy(file, tarReader); err != nil {
			_ = file.Close()
			return err
		}

		_ = file.Close()
	}

	return nil
}
