package trojansourcedetector_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func createTestData(t *testing.T, targetDir string) {
	data := map[string]string{
		"unicode.txt": "Hello world with an accent: \u00E1",
		"bidi.txt":    "Hello world with a BIDI control character: \u202A",
	}

	for file, contents := range data {
		dir := filepath.Join(targetDir, filepath.Dir(file))
		if err := os.MkdirAll(dir, 0750); err != nil {
			t.Fatalf("failed to create test data directory %s (%v)", dir, err)
		}
		if err := ioutil.WriteFile(
			filepath.Join(targetDir, file),
			[]byte(contents),
			0600,
		); err != nil {
			t.Fatalf("failed to write %s (%v)", file, err)
		}
	}

	dir := filepath.Join(targetDir, "testdir")
	if err := os.MkdirAll(dir, 0750); err != nil {
		t.Fatalf("failed to create test directory %s (%v)", dir, err)
	}

	symlink := filepath.Join(targetDir, "testsymlink")
	if err := os.Symlink(dir, symlink); err != nil {
		t.Fatalf("failed to create symlink from %s to %s (%v)", symlink, dir, err)
	}

	symlink2 := filepath.Join(targetDir, "testdanglingsymlink")
	if err := os.Symlink("nonexistent", symlink2); err != nil {
		t.Fatalf("failed to create symlink from %s to nonexistent (%v)", symlink, err)
	}
}
