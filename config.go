package trojansourcedetector

import (
	"fmt"
	"os"
)

// Defaults adds the default value for each filed.
func (c *Config) Defaults() {
	c.Directory = "."
	c.DetectBIDI = true
	c.Exclude = []string{
		// Version control systems
		".git/**",
		".hg/**",

		// Image files. Do not include SVG here.
		"**/*.ai",
		"**/*.bmp",
		"**/*.gif",
		"**/*.jpg",
		"**/*.png",
		"**/*.psd",
		"**/*.ico",
		"**/*.icns",

		// Audio files.
		"**/*.wav",
		"**/*.mp3",
		"**/*.aiff",
		"**/*.au",
		"**/*.ogg",

		// Executable or object files
		"**/*.a",
		"**/*.bin",
		"**/*.exe",
		"**/*.so",
		"**/*.dll",
		"**/*.dylib",

		// Compressed files
		"**/*.tar",
		"**/*.zip",
		"**/*.gz",
		"**/*.bz",
		"**/*.bz2",
		"**/*.z",
		"**/*.lzo",
		"**/*.rar",
		"**/*.ear",
		"**/*.xz",
		"**/*.zstd",

		// Python
		"**/*.whl",
		"**/*.pyc",
		"**/__pycache__/**",

		// Java
		"**/*.class",
		"**/*.jar",

		// PHP
		"**/*.phar",

		// Fonts
		"**/*.woff",
		"**/*.woff2",

		// Translations
		"**/*.mo",
		"**/*.pot",

		// Documents
		"**/*.doc",
		"**/*.docx",
		"**/*.odf",
		"**/*.pdf",
		"**/*.chm",

		// Spreadsheets
		"**/*.xls",
		"**/*.xlsx",
		"**/*.ods",

		// Presentations
		"**/*.ppt",
		"**/*.pptx",
		"**/*.ods",

		// Packet capture
		"**/*.pcap",
	}
	c.DetectUnicode = false
	c.Parallelism = 10
}

// Config is the base configuration structure of the trojan source detector.
type Config struct {
	// Directory indicates the directory to scan.
	Directory string `json:"directory"`
	// Include limits the paths to only include the listed files or directories. Basic pattern matching is supported
	// as described in the Go filepath pattern matching docs: https://pkg.go.dev/path/filepath#Match
	Include []string `json:"include"`
	// Exclude contains a list of files or paths never to include in the scan. Basic pattern matching is supported as
	// described in the Go filepath pattern matching docs: https://pkg.go.dev/path/filepath#Match
	Exclude []string `json:"exclude"`
	// DetectUnicode will alert if any unicode characters are found. This may cause problems if comments, author
	// names, etc. are written in their national language. Use with caution.
	DetectUnicode bool `json:"detect_unicode"`
	// DetectBIDI will alert if any bidirectional control characters mentioned in the trojan source paper are found.
	// See https://trojansource.codes/trojan-source.pdf
	DetectBIDI bool `json:"detect_bidi"`
	// Parallelism sets how many parallel scans should run.
	Parallelism uint `json:"parallelism"`
}

// Validate validates the configuration and returns an error if the configuration is invalid.
func (c *Config) Validate() error {
	if _, err := os.Stat(c.Directory); err != nil {
		return fmt.Errorf("failed to stat directory %s (%w)", c.Directory, err)
	}
	if !c.DetectBIDI && !c.DetectUnicode {
		return fmt.Errorf("no detectors enabled")
	}
	if c.Parallelism == 0 {
		return fmt.Errorf("parallelism can't be zero")
	}
	return nil
}
