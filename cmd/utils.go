package cmd

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Zip struct {
	z     *zip.ReadCloser
	files map[string]*zip.File
}

func OpenZip(filename string) (*Zip, error) {
	z, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}

	var ret = &Zip{
		z:     z,
		files: make(map[string]*zip.File, len(z.File)),
	}
	for _, f := range z.File {
		ret.files[f.Name] = f
	}
	return ret, nil
}

func (m *Zip) Close() {
	if m.z != nil {
		m.z.Close()
	}
}

func (m *Zip) Exists(item string) bool {
	_, ok := m.files[item]
	return ok
}

func (m *Zip) ReadBytes(item string) ([]byte, error) {
	rc, err := m.files[item].Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	b, _ := io.ReadAll(rc)
	return b, nil
}

func (m *Zip) ReadJson(item string, v interface{}) error {
	data, err := m.ReadBytes(item)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &v)
	return err
}

func (m *Zip) CalcHash256(item string) string {
	rc, err := m.files[item].Open()
	if err != nil {
		return ""
	}
	defer rc.Close()

	// Create a new SHA-256 hash object
	hasher := sha256.New()

	// Copy the file content to the hasher
	if _, err := io.Copy(hasher, rc); err != nil {
		return ""
	}

	// Get the final hash sum and encode it to a hexadecimal string
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func (m *Zip) Extract(item string, filename string) error {
	rc, err := m.files[item].Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return err
	}
	return nil
}

// humanReadableBytes converts a number of bytes to a human-readable string.
func HumanReadableBytes(bytes uint64) string {
	const (
		_  = iota // ignore first value by assigning to blank identifier
		KB = 1 << (10 * iota)
		MB
		GB
		TB
		PB
	)

	var size string

	switch {
	case bytes >= PB:
		size = fmt.Sprintf("%.2f PB", float64(bytes)/float64(PB))
	case bytes >= TB:
		size = fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		size = fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		size = fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		size = fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		size = fmt.Sprintf("%d B", bytes)
	}

	return size
}
