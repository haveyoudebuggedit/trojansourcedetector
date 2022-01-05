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

| Field | Description |
|-------|-------------|
| `directory` | Directory to run the check on. Defaults to the current directory. |
| `include` | A list of files to include in the scan. Paths should always be written in Linux syntax with forward slashes and begin with the project directory. Basic pattern matching is supported via [Go filepath](https://pkg.go.dev/path/filepath#Match). Defaults to empty (all files). |
| `exclude` | A list of files to exclude from the scan. Paths should always be written in Linux syntax with forward slashes and begin with the project directory. Basic pattern matching is supported via [Go filepath](https://pkg.go.dev/path/filepath#Match). Defaults to `.git` and all its subdirectories. |
| `detect_unicode` | Alert for all non-ASCII unicode characters. Defaults to false. |
| `detect_bidi` | Detect bidirectional control characters. These can cause the trojan source problem. Defaults to true. |
| `parallelism` | How many files to check in parallel. Defaults to 10. |

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