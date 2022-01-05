package trojansourcedetector_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/haveyoudebuggedit/trojansourcedetector"
)

func TestE2E(t *testing.T) {
	testDir := t.TempDir()

	createTestData(t, testDir)

	assertFileExists(t, filepath.Join(testDir, "bidi.txt"))
	assertFileExists(t, filepath.Join(testDir, "unicode.txt"))

	detector := trojansourcedetector.New(&trojansourcedetector.Config{
		Directory:     testDir,
		DetectUnicode: true,
		DetectBIDI:    true,
		Parallelism:   10,
	})

	errs := detector.Run()

	for _, e := range errs.Get() {
		data, err := e.JSON()
		if err != nil {
			t.Fatalf("failed to generate JSON from error (%v)", err)
		}
		fmt.Printf("%s\n", data)
	}

	assertHasError(t, errs, trojansourcedetector.ErrBIDI, "bidi.txt", 1, 44)
	assertHasError(t, errs, trojansourcedetector.ErrUnicode, "unicode.txt", 1, 29)
	assertHasError(t, errs, trojansourcedetector.ErrUnicode, "unicode.txt", 1, 30)
	assertHasNoErrors(t, errs, "testsymlink")
}

func assertHasNoErrors(t *testing.T, errs trojansourcedetector.Errors, file string) {
	for _, err := range errs.Get() {
		if filepath.ToSlash(err.File()) == file {
			t.Fatalf("unexpected %s error in file %s (%s)", err.Code(), err.File(), err.Details())
		}
	}
}

func assertFileExists(t *testing.T, file string) {
	if _, err := os.Stat(file); err != nil {
		t.Fatalf("file does not exist: %s (%v)", file, err)
	}
}

func assertHasError(
	t *testing.T,
	errs trojansourcedetector.Errors,
	code trojansourcedetector.ErrorCode,
	file string,
	line uint,
	column uint,
) {
	for _, err := range errs.Get() {
		if err.Code() == code && filepath.ToSlash(err.File()) == file && err.Line() == line && err.Column() == column {
			return
		}
	}
	t.Fatalf("Did not find expected '%s' error in %s line %d column %d", code, file, line, column)
}
