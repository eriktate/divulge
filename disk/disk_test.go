// +build integration

package disk_test

import (
	"context"
	"os"
	"testing"

	"github.com/eriktate/divulge/disk"
)

func Test_FileStore(t *testing.T) {
	// SETUP
	ctx := context.TODO()
	basePath := os.TempDir()
	fs := disk.New(basePath)

	testKey := "test_file.md"
	testData := "this is some test _markdown_"

	// RUN
	writeErr := fs.Write(ctx, testKey, []byte(testData))
	readData, readErr := fs.Read(ctx, testKey)

	// ASSERT
	if writeErr != nil {
		t.Fatalf("unexpected error: %s", writeErr)
	}

	if readErr != nil {
		t.Fatalf("unexpected error: %s", readErr)
	}

	if string(readData) != testData {
		t.Fatalf("unexpected read data: %s", string(readData))
	}
}
