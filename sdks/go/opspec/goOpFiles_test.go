package opspec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoOpFilesMatchCreate(t *testing.T) {
	files, err := GoOpFiles("hello", "Greets the caller.", GoSDK{})
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(t.TempDir(), "hello")
	if err := CreateGo(dir, "hello", "Greets the caller.", GoSDK{}); err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		data, err := os.ReadFile(filepath.Join(dir, file.Name))
		if err != nil || string(data) != file.Content {
			t.Fatalf("%s: %v; %q", file.Name, err, data)
		}
	}
}
