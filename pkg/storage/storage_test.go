// File: pkg/storage/storage_test.go

package storage

import "testing"

// TestSaveData verifies that SaveData always returns nil regardless of the input.
// This ensures full test coverage for the stub implementation.
func TestSaveData(t *testing.T) {
	// Test with non-empty data.
	testData := map[string]string{"example": "data"}
	options := []StorageOption{JSON, XML, Excel, MongoDB, MySQL}

	for _, opt := range options {
		if err := SaveData(testData, opt); err != nil {
			t.Errorf("SaveData returned an error for option %v: %v", opt, err)
		}
	}

	// Also test with an empty map.
	if err := SaveData(map[string]string{}, JSON); err != nil {
		t.Errorf("SaveData returned an error for empty map: %v", err)
	}
}
