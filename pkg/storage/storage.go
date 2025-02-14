// File: pkg/storage/storage.go
package storage

// StorageOption enumerates the types of storage we might support.
type StorageOption int

const (
	JSON StorageOption = iota
	XML
	Excel
	MongoDB
	MySQL
)

// SaveData will eventually accept the extracted data and store it in various formats.
// This could later become multiple functions or a strategy pattern.
func SaveData(data map[string]string, option StorageOption) error {
	// Stub: for now, do nothing.
	return nil
}
