package trojansourcedetector //nolint:testpackage // Excluded because the test is internal.

import (
	"fmt"
	"testing"
)

type globTestData struct {
	pattern        string
	shouldMatch    []string
	shouldNotMatch []string
}

func TestGlob(t *testing.T) {
	globData := []globTestData{
		{
			pattern: "*",
			shouldMatch: []string{
				"test.txt",
				"foo",
			},
			shouldNotMatch: []string{
				".git/config",
			},
		},
		{
			pattern: "?",
			shouldMatch: []string{
				"t",
				"9",
			},
			shouldNotMatch: []string{
				"ab",
				"test",
				"a/b",
				".git/config",
			},
		},
		{
			pattern: "*.txt",
			shouldMatch: []string{
				"test.txt",
				"foo.bar.txt",
			},
			shouldNotMatch: []string{
				"test.bar",
				"test/test.txt",
			},
		},
		{
			pattern: "\\*.txt",
			shouldMatch: []string{
				"*.txt",
			},
			shouldNotMatch: []string{
				"test.txt",
				"foo.bar.txt",
				"test.bar",
				"test/test.txt",
			},
		},
		{
			pattern: "**",
			shouldMatch: []string{
				"test.txt",
				"foo",
				".git/config",
			},
			shouldNotMatch: []string{},
		},
		{
			pattern: "**/*.txt",
			shouldMatch: []string{
				"test.txt",
				"test/test.txt",
				"test/test2/test3.txt",
			},
			shouldNotMatch: []string{
				"foo",
				".git/config",
			},
		},
		{
			pattern: "[a-z]*",
			shouldMatch: []string{
				"t",
				"test",
			},
			shouldNotMatch: []string{
				"test.txt",
				"test/test",
				"test/test2/test3",
				".git/config",
			},
		},
		{
			pattern: "?/[a-z]*/**/foo",
			shouldMatch: []string{
				"t/a/abc/foo",
				"b/test/abc/foo",
			},
			shouldNotMatch: []string{
				"t/a/abc",
				"b/test/abc",
				"test.txt",
				"test/test",
				"test/test2/test3",
				".git/config",
			},
		},
	}

	for _, entry := range globData {
		t.Run(fmt.Sprintf("Pattern %s", entry.pattern), func(t *testing.T) {
			t.Logf("Compiling pattern %s...", entry.pattern)
			compiledPattern, err := compile(entry.pattern)
			if err != nil {
				t.Fatalf("Failed to compile pattern %s (%v).", entry.pattern, err)
			}
			for _, shouldMatch := range entry.shouldMatch {
				t.Run(fmt.Sprintf("Should match %s", shouldMatch), func(t *testing.T) {
					t.Logf("Checking if %s matches %s...", entry.pattern, shouldMatch)
					if !compiledPattern.match(shouldMatch) {
						t.Fatalf("Pattern %s should have matched %s, but it didn't. Compiled regexp was %s", entry.pattern, shouldMatch, compiledPattern.string())
					}
				})
			}
			for _, shouldNotMatch := range entry.shouldNotMatch {
				t.Run(fmt.Sprintf("Should not match %s", shouldNotMatch), func(t *testing.T) {
					t.Logf("Checking if %s does not match %s...", entry.pattern, shouldNotMatch)
					if compiledPattern.match(shouldNotMatch) {
						t.Fatalf("Pattern %s should not have matched %s, but it did. Compiled regexp was %s", entry.pattern, shouldNotMatch, compiledPattern.string())
					}
				})
			}
		})
	}
}
