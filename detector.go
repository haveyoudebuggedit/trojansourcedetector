package trojansourcedetector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// New creates a composite detector based on a configuration passed.
func New(config *Config) Detector {
	var fileDetectors []SingleFileDetector

	if config == nil {
		config = &Config{}
		config.Defaults()
	}

	if config.DetectBIDI {
		fileDetectors = append(fileDetectors, &bidiDetector{})
	}
	if config.DetectUnicode {
		fileDetectors = append(fileDetectors, &unicodeDetector{})
	}

	return &detector{
		config:        config,
		fileDetectors: fileDetectors,
	}
}

// Detector detects malicious unicode code points in a directory based on the configuration.
type Detector interface {
	// Run runs the detection algorithm. The returned list of errors contains the violations of the
	// passed rule set. If reading the directory fails, the list will contain an error entry without a file name.
	Run() Errors
}

type detector struct {
	fileDetectors []SingleFileDetector
	config        *Config
}

func (d *detector) Run() Errors {
	lock := make(chan struct{}, d.config.Parallelism)
	wg := &sync.WaitGroup{}
	container := NewErrors()
	err := filepath.Walk(d.config.Directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			switch info.Mode().Type() {
			case os.ModeDir:
				// Subdirectories will be listed anyway.
				return nil
			case os.ModeSymlink:
				// We ignore symlinks because they either point to a file in the project, which is checked
				// anyway, or it points outside of the project, in which case we don't care.
				return nil
			case os.ModeDevice:
				// We can't do anything with device nodes.
				return nil
			case os.ModeNamedPipe:
				// We can't do anything with named pipes.
				return nil
			case os.ModeSocket:
				// We can't do anything with sockets.
				return nil
			case os.ModeCharDevice:
				// We can't do anything with character devices.
				return nil
			case os.ModeIrregular:
				// Irregular file?
				return nil
			}
			path = filepath.ToSlash(path)
			removePrefix := filepath.ToSlash(d.config.Directory)
			if !strings.HasSuffix(removePrefix, "/") {
				removePrefix = fmt.Sprintf("%s/", removePrefix)
			}
			reportedPath := strings.TrimPrefix(path, removePrefix)

			if len(d.config.Include) != 0 {
				match := false
				for _, include := range d.config.Include {
					matches, err := filepath.Match(include, reportedPath)
					if err != nil {
						return err
					}
					if matches {
						match = true
						break
					}
				}
				if !match {
					return nil
				}
			}
			for _, exclude := range d.config.Exclude {
				matches, err := filepath.Match(exclude, reportedPath)
				if err != nil {
					return err
				}
				if matches {
					return nil
				}
			}
			wg.Add(1)
			go d.processFile(path, reportedPath, lock, wg, container)
			return nil
		})
	wg.Wait()
	if err != nil {
		container.Add(
			ErrIODirectory,
			err.Error(),
			"",
			0,
			0,
		)
	}
	return container
}

func (d *detector) processFile(
	path string,
	reportedPath string,
	lock chan struct{},
	wg *sync.WaitGroup,
	container Errors,
) {
	lock <- struct{}{}
	defer func() {
		<-lock
		wg.Done()
	}()

	fh, err := os.Open(path) //nolint:gosec
	if err != nil {
		container.Add(
			ErrIOFile,
			err.Error(),
			path,
			0,
			0,
		)
	}
	//nolint:gosec // Ignore G307 for now.
	defer func() {
		_ = fh.Close()
	}()

	for _, detector := range d.fileDetectors {
		if _, err := fh.Seek(0, io.SeekStart); err != nil {
			container.Add(ErrIOSeek, err.Error(), path, 0, 0)
			continue
		}
		container.AddAll(detector.Detect(reportedPath, bufio.NewReader(fh)))
	}
}
