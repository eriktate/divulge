package disk

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
)

// A FileStore is an on-disk implementation of divulge.FileStore.
type FileStore struct {
	basePath string
}

// New returns a new FileStore given a base path to operate on.
func New(basePath string) FileStore {
	return FileStore{basePath}
}

// Write a file.
func (fs FileStore) Write(ctx context.Context, key string, data []byte) error {
	fullPath := fmt.Sprintf("%s/%s", fs.basePath, key)
	if err := ioutil.WriteFile(fullPath, data, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Read a file.
func (fs FileStore) Read(ctx context.Context, key string) ([]byte, error) {
	fullPath := fmt.Sprintf("%s/%s", fs.basePath, key)
	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}
