package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/haveyoudebuggedit/trojansourcedetector"
)

const defaultConfigFile = ".trojansourcedetector.json"

func main() {
	configFile := defaultConfigFile

	flag.StringVar(&configFile, "config", configFile, "JSON file containing the configuration.")
	flag.Parse()

	cfg, err := readConfigFile(configFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) || configFile != defaultConfigFile {
			panic(err)
		}
	}

	detector := trojansourcedetector.New(cfg)
	errs := detector.Run()
	if err := writeErrors(errs); err != nil {
		panic(err)
	}
	if len(errs.Get()) > 0 {
		os.Exit(1)
	} else {
		fmt.Printf("No errors found.")
	}
}

func writeErrors(errors trojansourcedetector.Errors) error {
	for _, e := range errors.Get() {
		if os.Getenv("GITHUB_ACTIONS") != "" {
			fmt.Printf(
				"::error file=%s,line=%d,col=%d,title=%s::%s\n",
				e.File(),
				e.Line(),
				e.Column(),
				e.Code(),
				fmt.Sprintf("%s: %s", e.Code(), e.Details()),
			)
			fmt.Printf("%s in %s line %d column %d (%s)\n", e.Code(), e.File(), e.Line(), e.Column(), e.Details())
		} else {
			encoded, err := e.JSON()
			if err != nil {
				return fmt.Errorf("bug: failed to encode error entry (%w)", err)
			}
			fmt.Printf("%s\n", encoded)
		}
	}
	return nil
}

func readConfigFile(file string) (*trojansourcedetector.Config, error) {
	result := &trojansourcedetector.Config{}
	result.Defaults()
	fh, err := os.Open(file) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s (%w)", file, err)
	}
	//nolint:gosec // Ignore G307 for now.
	defer func() {
		_ = fh.Close()
	}()
	decoder := json.NewDecoder(fh)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode configuration file %s (%w)", file, err)
	}
	return result, nil
}
