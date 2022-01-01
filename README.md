# Trojan Source Detector

This application detects [Trojan Source](https://trojansource.codes) attacks in source code. It can be used as part of the CI system to make sure there are no trojan source / unicode bi-directional text attacks in a pull request.

## Usage

This utility can be used either on GitHub Actions:

```yaml
jobs:
  trojansource:
    name: Trojan Source Detection
    runs-on: ubuntu-latest
    steps:
      # Checkout your project with git
      - name: Checkout
        uses: actions/checkout@v2
      # Run trojansourcedetector
      - name: Trojan Source Detector
        uses: haveyoudebuggedit/trojansourcedetector@v1
```

You can also run it on any CI system by simply downloading the [released binary](https://github.com/haveyoudebuggedit/trojansourcedetector/releases) and running:

```
./trojansourcedetector
```

## Configuration

You can customize the behavior by providing a config file. This file is named `.trojansourcedetector.json` by default and has the following fields:

| Field            | Description                                                                                                                                                                                                                                                                                                         |
|------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `directory`      | Directory to run the check on. Defaults to the current directory.                                                                                                                                                                                                                                                   |
| `include`        | A list of files to include in the scan. Paths should always be written in Linux syntax with forward slashes and begin with the project directory. For supported patterns see the [Globbing section](#globbing) below. Defaults to empty (all files).                                                                |
| `exclude`        | A list of files to exclude from the scan. Paths should always be written in Linux syntax with forward slashes and begin with the project directory. For supported patterns see the [Globbing section](#globbing) below. For defaults see the [Files excluded by default section below](#files-excluded-by-default). |
| `detect_unicode` | Alert for all non-ASCII unicode characters. Defaults to false.                                                                                                                                                                                                                                                      |
| `detect_bidi`    | Detect bidirectional control characters. These can cause the trojan source problem. Defaults to true.                                                                                                                                                                                                               |
| `parallelism`    | How many files to check in parallel. Defaults to 10.                                                                                                                                                                                                                                                                |

For an example you can take a look at the [.trojansourcedetector.json](.trojansourcedetector.json) in this repository.

If you want to use a different file name, you can change your GitHub Actions config:

```yaml
jobs:
  trojansource:
    name: Trojan Source Detection
    runs-on: ubuntu-latest
    steps:
      # Checkout your project with git
      - name: Checkout
        uses: actions/checkout@v2
      # Run trojansourcedetector
      - name: Trojan Source Detector
        uses: haveyoudebuggedit/trojansourcedetector@v1
        with:
          config: path/to/config/file
```

Or, if you are using the command line version, you can simply pass the `-config` option with the appropriate config file.

## Globbing

When including and excluding files the following patterns are supported:

- `?` matches any single character, except for the path separator.
- `*` matches any character sequence, except for the path separator.
- `**` matches zero or more path segments.
- `[a-z]` matches a single character that falls in this character class.
- `[^a]` matches a single character that is not `a`.
- `[a-z]*` matches a sequence of characters within this character class.

**Note:** In order to match files in subdirectories, patterns must be prefixed with `**/`.

**Note:** File patterns should always be written with the *nix notation (`/`) as a path separator.

## Files excluded by default

By default, the following patterns are excluded:

```
// Version control systems
.git/**

// Image files. Do not include SVG here.
**/*.gif
**/*.jpg
**/*.png
**/*.psd
**/*.ico

// Audio files.
**/*.wav
**/*.mp3
**/*.png
**/*.aiff
**/*.au
**/*.ogg

// Executable or object files
**/*.exe
**/*.so
**/*.dll
**/*.dylib

// Compressed files
**/*.tar
**/*.zip
**/*.gz
**/*.bz
**/*.z
**/*.rar
**/*.ear

// Python
**/*.whl
**/*.pyc
**/__pycache__/**

// Java
**/*.class
**/*.jar

// PHP
**/*.phar

// Translations
**/*.mo
**/*.pot
```

## Building

This tool can be built using Go 1.17 or higher:

```
go build cmd/trojansourcedetector/main.go
```

## Running tests

In order to run tests, you will need to run the following command:

```
go test -v ./...
```